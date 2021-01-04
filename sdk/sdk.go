package sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Sdk struct {
	session string
	al      string
}

func NewSdk(cookiePath string) Sdk {
	var sdk Sdk
	bs, err := ioutil.ReadFile(cookiePath)
	if err != nil {
		fmt.Println("Error retrieving session:", err)
		os.Exit(1)
	}

	fileLines := strings.Split(string(bs), "\n")
	for _, fileLine := range fileLines {
		if strings.HasPrefix(fileLine, "sessions_gog_com=") {
			sdk.session = strings.TrimPrefix(fileLine, "sessions_gog_com=")
		} else if strings.HasPrefix(fileLine, "gog-al=") {
			sdk.al = strings.TrimPrefix(fileLine, "gog-al=")
		}
	}
	return sdk
}

func (s Sdk) getClient() http.Client {
	cs := []*http.Cookie{
		&http.Cookie{Name: "sessions_gog_com", Value: s.session},
		&http.Cookie{Name: "gog-al", Value: s.al},
	}
	j := Jar{cookies: cs}
	return http.Client{Jar: &j}
}
