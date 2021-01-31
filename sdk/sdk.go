package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Sdk struct {
	session string
	al      string
	logger  *log.Logger
}

func NewSdk(cookiePath string, logger *log.Logger) (*Sdk, error) {
	sdk := Sdk{session: "", al: "", logger: logger}
	bs, err := ioutil.ReadFile(cookiePath)
	if err != nil {
		msg := fmt.Sprintf(" Error retrieving session: %s", err.Error())
		return &sdk, errors.New(msg)
	}

	fileLines := strings.Split(string(bs), "\n")
	for _, fileLine := range fileLines {
		if strings.HasPrefix(fileLine, "sessions_gog_com=") {
			sdk.session = strings.TrimPrefix(fileLine, "sessions_gog_com=")
		} else if strings.HasPrefix(fileLine, "gog-al=") {
			sdk.al = strings.TrimPrefix(fileLine, "gog-al=")
		}
	}
	return &sdk, nil
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

func (s *Sdk) getUrl(url string, fnCall string, debug bool, jsonBody bool) ([]byte, error) {
	c := (*s).getClient(true)

	if debug {
		(*s).logger.Println(fmt.Sprintf("%s -> GET %s", fnCall, url))
	}
	r, err := c.Get(url)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fnCall, err.Error())
		return nil, errors.New(msg)
	}

	b, bErr := ioutil.ReadAll(r.Body)
	if bErr != nil {
		msg := fmt.Sprintf("%s -> retrieval body error: %s", fnCall, bErr.Error())
		return nil, errors.New(msg)
	}
	if debug {
		if jsonBody {
			var out bytes.Buffer
			jErr := json.Indent(&out, b, "", "  ")
			if jErr != nil {
				msg := fmt.Sprintf("%s -> json parsing error: %s", fnCall, jErr.Error())
				return nil, errors.New(msg)
			}
			b = out.Bytes()
		}
		(*s).logger.Println(
			fmt.Sprintf("%s -> response body:", fnCall),
			string(b),
		)
	}

	return b, nil
}
