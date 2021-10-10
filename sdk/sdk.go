package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"gogcli/logging"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Sdk struct {
	session string
	al      string
	maxRetries int64
	currentRetries int64
	logger  *logging.Logger
}

func NewSdk(cookie GogCookie, logSource *logging.Source) *Sdk {
	logger := logSource.CreateLogger(os.Stdout, "[sdk] ", log.Lmsgprefix)
	sdk := Sdk{session: cookie.Session, al: cookie.Al, maxRetries: 3, currentRetries: 0, logger: logger}
	return &sdk
}

func (s *Sdk) incRetries() {
	(*s).currentRetries += 1
}

func (s *Sdk) resetRetries() {
	(*s).currentRetries = 0
}

func (s *Sdk) maxRetriesReached() bool {
	return (*s).currentRetries == (*s).maxRetries
}

func (s *Sdk) getClient(followRedirects bool) http.Client {
	cs := []*http.Cookie{
		&http.Cookie{Name: "sessions_gog_com", Value: (*s).session},
		&http.Cookie{Name: "gog-al", Value: (*s).al},
	}
	j := Jar{cookies: cs}
	if followRedirects{
		return http.Client{Jar: &j}
	} else {
		return http.Client{
			Jar: &j,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}
	}
}

func (s *Sdk) getUrl(url string, fnCall string, jsonBody bool) ([]byte, int, error) {
	c := (*s).getClient(true)

	(*s).logger.Debug(fmt.Sprintf("%s -> GET %s", fnCall, url))
	
	r, err := c.Get(url)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fnCall, err.Error())
		return nil, -1, errors.New(msg)
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode > 299 {
		if r.StatusCode >= 500 && (!(*s).maxRetriesReached()) {
			(*s).logger.Warning(fmt.Sprintf("%s -> GET %s failed with code %d. Will retry.", fnCall, url, r.StatusCode))
			(*s).incRetries()
			return (*s).getUrl(url, fnCall, jsonBody)
		}
		(*s).resetRetries()
		msg := fmt.Sprintf("%s -> retrieval request error: did not expect status code of %d", fnCall, r.StatusCode)
		return nil, r.StatusCode, errors.New(msg)
	}

	(*s).resetRetries()

	b, bErr := ioutil.ReadAll(r.Body)
	if bErr != nil {
		msg := fmt.Sprintf("%s -> retrieval body error: %s", fnCall, bErr.Error())
		return nil, r.StatusCode, errors.New(msg)
	}

	if jsonBody {
		var out bytes.Buffer
		jErr := json.Indent(&out, b, "", "  ")
		if jErr != nil {
			msg := fmt.Sprintf("%s -> json parsing error: %s", fnCall, jErr.Error())
			return nil, r.StatusCode, errors.New(msg)
		}
		b = out.Bytes()
	}
	(*s).logger.Debug(fmt.Sprintf("%s -> response body: %s", fnCall, string(b)))

	return b, r.StatusCode, nil
}
