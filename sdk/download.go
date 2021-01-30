package sdk

import (
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
)

func (s *Sdk) GetDownloadHandle(downloadPath string, debug bool) (io.ReadCloser, int, string, error) {
	fn := fmt.Sprintf("GetDownloadHandle(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://gog.com%s", downloadPath)

	var body io.ReadCloser
	bodyLength := 0
	filename := ""

	c := (*s).getClient()
	if debug {
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
	if debug {
		(*s).logger.Println(fmt.Sprintf("    Final Url: %s", finalURL))
	}

	p := (*r.Request.URL).Path
	filename = path.Base(p)

	return body, bodyLength, filename, nil
}
