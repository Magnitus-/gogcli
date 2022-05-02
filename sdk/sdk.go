package sdk

import (
	"gogcli/logging"
	"log"
	"net/http"
	"os"
	"time"
)

type Sdk struct {
	session        string
	al             string
	maxRetries     int64
	retryPause     time.Duration
	logger         *logging.Logger
}

func NewSdk(cookie GogCookie, logSource *logging.Source) *Sdk {
	logger := logSource.CreateLogger(os.Stdout, "[sdk] ", log.Lmsgprefix)
	pause, _ := time.ParseDuration("100ms")
	sdk := Sdk{
		session: cookie.Session, 
		al: cookie.Al, 
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
	cs := []*http.Cookie{
		&http.Cookie{Name: "sessions_gog_com", Value: (*s).session},
		&http.Cookie{Name: "gog-al", Value: (*s).al},
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