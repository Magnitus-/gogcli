package sdk

import (
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
	"net/url"
)

func (s *Sdk) GetDownloadFilename(downloadPath string) (string, error) {
	fn := fmt.Sprintf("GetDownloadFilename(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	c := (*s).getClient(false)
	if (*s).debug {
		(*s).logger.Println(fmt.Sprintf("%s -> GET %s", fn, u))
	}

	r, err := c.Get(u)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fn, err.Error())
		return "", errors.New(msg)
	}

	if r.StatusCode != 302 {
		msg := fmt.Sprintf("%s -> Expected response status code of 302, but got %d", fn, r.StatusCode)
		return "", errors.New(msg)
	}

	locHeader, ok := r.Header["Location"]
	var location string
	if !ok {
		msg := fmt.Sprintf("%s -> Expected location header in response, but it was missing", fn)
		return "", errors.New(msg)
	} else {
		location = locHeader[0]
		if (*s).debug {
			(*s).logger.Println(fmt.Sprintf("%s -> Location Header: %s", fn, location))
		}
	}

	locUrl, err := url.Parse(location)
	if err != nil {
		msg := fmt.Sprintf("%s -> Error parsing location header url: %s", fn, err.Error())
		return "", errors.New(msg)
	}

	queryParams := locUrl.Query()
	pathParam, ok := queryParams["path"]
	if !ok {
		msg := fmt.Sprintf("%s -> Error location header url does not have the expected path query parameter", fn)
		return "", errors.New(msg)	
	}

	return path.Base(pathParam[0]), nil
}

func (s *Sdk) GetDownloadHandle(downloadPath string) (io.ReadCloser, int, string, error) {
	fn := fmt.Sprintf("GetDownloadHandle(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	var body io.ReadCloser
	bodyLength := 0
	filename := ""

	c := (*s).getClient(true)
	if (*s).debug {
		(*s).logger.Println(fmt.Sprintf("%s -> GET %s", fn, u))
	}

	r, err := c.Get(u)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fn, err.Error())
		return nil, 0, "", errors.New(msg)
	}

	body = r.Body

	clHeader, ok := r.Header["Content-Length"]
	if ok {
		l, lErr := strconv.Atoi(clHeader[0])
		if lErr != nil {
			(*s).logger.Println(fmt.Sprintf("%s -> Cannot return exact download size as Content-Length header is not parsable. Will set it to 0.", fn))
		} else {
			bodyLength = l
		}
	} else {
		(*s).logger.Println(fmt.Sprintf("%s -> Cannot return exact download size as Content-Length header is not found. Will set it to 0.", fn))
	}

	finalURL := r.Request.URL.String()
	if (*s).debug {
		(*s).logger.Println(fmt.Sprintf("    Final Url: %s", finalURL))
	}

	p := (*r.Request.URL).Path
	filename = path.Base(p)

	return body, bodyLength, filename, nil
}
