package sdk

import (
	"gogcli/logging"
	"log"
	"net/http"
	"os"
	"time"
)

type Sdk struct {
	cookies    []*http.Cookie
	maxRetries int64
	retryPause time.Duration
	logger     *logging.Logger
}

func NewSdk(cookies []*http.Cookie, logSource *logging.Source) *Sdk {
	logger := logSource.CreateLogger(os.Stdout, "[sdk] ", log.Lmsgprefix)
	pause, _ := time.ParseDuration("100ms")
	
	sdk := Sdk{
		cookies: cookies, 
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
	j := Jar{cookies: (*s).cookies}
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