package sdk

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

func (s *Sdk) GetDownloadHandle(path string, debug bool) (io.ReadCloser, int, error) {
	fn := fmt.Sprintf("GetDownloadHandle(path=%s)", path)
	url := fmt.Sprintf("https://gog.com%s", path)
	redirect := true

	var body io.ReadCloser
	bodyLength := 0

	c := (*s).getClient()
	for redirect {
		if debug {
			(*s).logger.Println(fmt.Sprintf("%s -> GET %s", fn, url))
		}

		r, err := c.Get(url)
		if err != nil {
			msg := fmt.Sprintf("%s -> retrieval request error: %s", fn, err.Error())
			return nil, 0, errors.New(msg)
		}

		if r.StatusCode == 302 {
			r.Body.Close()
			if _, ok := r.Header["Location"]; ok {
				url = r.Header["Location"][0]
			} else if _, ok := r.Header["location"]; ok {
				url = r.Header["location"][0]
			} else {
				msg := fmt.Sprintf("%s -> retrieval request error: Location header not found", fn)
				return nil, 0, errors.New(msg)
			}
		} else {
			redirect = false
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
		}
	}

	return body, bodyLength, nil
}
