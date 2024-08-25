package sdk

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type FirefoxCookie struct {
	RequestCookies map[string]string `json:"Request Cookies"`
}

func readDefaultCookie(content string) ([]*http.Cookie, error) {
	cs := []*http.Cookie{}
	lines := strings.Split(strings.Replace(content, "\r\n", "\n", -1), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "sessions_gog_com=") {
			cs = append(cs, &http.Cookie{Name: "sessions_gog_com", Value: strings.TrimPrefix(line, "sessions_gog_com=")})
		} else if strings.HasPrefix(line, "gog-al=") {
			cs = append(cs, &http.Cookie{Name: "gog-al", Value: strings.TrimPrefix(line, "gog-al=")})
		}
	}

	if len(cs) == 0 {
		return cs, errors.New("Could not parse cookie with default format. No cookie values were read.")
	}

    return cs, nil
}

func readStringCookie(content string) ([]*http.Cookie, error) {
    header := http.Header{}
    header.Add("Cookie", content)
    request := http.Request{Header: header}

	cs := request.Cookies()

	if len(cs) == 0 {
		return cs, errors.New("Could not parse cookie string. No cookie values were read.")
	}

    return cs, nil
}

func readNetscapeCookie(content string) ([]*http.Cookie, error) {
	cs := []*http.Cookie{}
	lines := strings.Split(strings.Replace(content, "\r\n", "\n", -1), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		lineFields := strings.Split(line, "\t")
		if len(lineFields) < 7 {
			continue
		}

		cs = append(cs, &http.Cookie{Name: lineFields[5], Value: lineFields[6]})
	}

	if len(cs) == 0 {
		return cs, errors.New("Could not parse Netscape cookie. No cookie values were read.")
	}

	return cs, nil
}

func readFirefoxCookie(content string) ([]*http.Cookie, error) {
	cs := []*http.Cookie{}

	fCookie := FirefoxCookie{}

	err := json.Unmarshal([]byte(content), &fCookie)
	if err != nil {
		return cs, errors.New(fmt.Sprintf("Following error occured while parsing Firefox cookie: %s", err.Error()))
	}

	for key, val := range fCookie.RequestCookies {
		cs = append(cs, &http.Cookie{Name: key, Value: val})
	}

	if len(cs) == 0 {
		return cs, errors.New("Could not parse Firefox cookie. No cookie values were read.")
	}

	return cs, nil
}

func ReadCookie(path string, kind string) ([]*http.Cookie, error) {
	if kind != "default" && kind != "string" && kind != "netscape" && kind != "firefox" {
		msg := fmt.Sprintf("Cookie type of %s is not supported", kind)
		return []*http.Cookie{}, errors.New(msg)
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		msg := fmt.Sprintf("Error reading file to retrieve cookie: %s", err.Error())
		return []*http.Cookie{}, errors.New(msg)
	}

	switch kind {
	case "default":
		return readDefaultCookie(string(bs))
	case "string":
		return readStringCookie(string(bs))
	case "firefox":
		return readFirefoxCookie(string(bs))
	case "netscape":
		return readNetscapeCookie(string(bs))
	}

	return []*http.Cookie{}, nil
}
