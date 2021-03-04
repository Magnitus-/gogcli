package login

import (
	"net/http"
)

func getClient(cs []*http.Cookie) http.Client {
	j := Jar{cookies: cs}
	return http.Client{
		Jar: &j,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}