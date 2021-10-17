package sdk

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

type GogCookie struct {
	Session string
	Al      string
}

func readDefaultCookie(lines []string) (string, string) {
	session := ""
	al := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "sessions_gog_com=") {
			session = strings.TrimPrefix(line, "sessions_gog_com=")
		} else if strings.HasPrefix(line, "gog-al=") {
			al = strings.TrimPrefix(line, "gog-al=")
		}
	}
	return session, al
}

func readNetscapeCookie(lines []string) (string, string) {
	session := ""
	al := ""
	for _, line := range lines {
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		lineFields := strings.Split(line, "\t")
		if len(lineFields) < 7 {
			continue
		}

		if lineFields[5] == "sessions_gog_com" {
			session = lineFields[6]
		} else if lineFields[5] == "gog-al" {
			al = lineFields[6]
		}
	}
	return session, al
}

func ReadCookie(path string, kind string) (GogCookie, error) {
	if kind != "default" && kind != "netscape" {
		msg := fmt.Sprintf("Cookie type of %s is not supported: %s", kind)
		return GogCookie{"", ""}, errors.New(msg)
	}

	bs, err := ioutil.ReadFile(path)
	if err != nil {
		msg := fmt.Sprintf("Error retrieving session: %s", err.Error())
		return GogCookie{"", ""}, errors.New(msg)
	}

	lines := strings.Split(string(bs), "\n")
	if kind == "netscape" {
		session, al := readNetscapeCookie(lines)
		return GogCookie{session, al}, nil
	}

	session, al := readDefaultCookie(lines)
	return GogCookie{session, al}, nil
}
