package sdk

import (
	"gogcli/logging"
	"log"
	"net/http"
	"os"
	"time"
)

type Sdk struct {
	cookie_values map[string]string
	maxRetries    int64
	retryPause    time.Duration
	logger        *logging.Logger
}

func NewSdk(cookie GogCookie, logSource *logging.Source) *Sdk {
	logger := logSource.CreateLogger(os.Stdout, "[sdk] ", log.Lmsgprefix)
	pause, _ := time.ParseDuration("100ms")
	
	cookie_values := make(map[string]string)
	cookie_values["gog-al"] = cookie.Al
	cookie_values["sessions_gog_com"] = cookie.Session
	
	sdk := Sdk{
		cookie_values: cookie_values, 
		maxRetries: 5, 
		retryPause: pause, 
		logger: logger,
	}
	return &sdk
}

func (s *Sdk) pauseAfterError() {
	time.Sleep((*s).retryPause)
}

func (s *Sdk) getClient(followRedirects bool) http.Client {
	cs := []*http.Cookie{}
	for key, val := range (*s).cookie_values {
		cs = append(cs, &http.Cookie{Name: key, Value: val})
	}

	j := Jar{cookies: cs}
	if followRedirects {
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