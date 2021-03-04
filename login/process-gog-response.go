package login

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func processHeaders(r *http.Response) map[string]string {
	headers := make(map[string]string)
	for name, value := range r.Header() {
		if name == "Location" {
			value = strings.Replace(value, "https://login.gog.com", "http://localhost:8080/login", -1)
			value = strings.Replace(value, "https://embed.gog.com", "http://localhost:8080/embed", -1)
			value = strings.Replace(value, "https://auth.gog.com", "http://localhost:8080/auth", -1)
			value = strings.Replace(value, "https://static.gog.com", "http://localhost:8080/static", -1)
			value = strings.Replace(value, "https://www.gog.com", "http://localhost:8080/www", -1)
		}
		headers[name] = value
	}

	return headers
}

func processBody(r *http.Response, headers map[string]string) []byte, error {
	defer r.Body.Close()
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(headers["Content-Type"], "text/html") ||  strings.HasPrefix(headers["Content-Type"], "application/x-javascript") {
		body := string(b)
		body = strings.replace(body, "location.href.replace(/^http:/, 'https:')", "location.href")
		body = strings.Replace(body, "https://login.gog.com", "http://localhost:8080/login", -1)
		body = strings.Replace(body, "https://embed.gog.com", "http://localhost:8080/embed", -1)
		body = strings.Replace(body, "https://auth.gog.com", "http://localhost:8080/auth", -1)
		body = strings.Replace(body, "https://static.gog.com", "http://localhost:8080/static", -1)
		body = strings.Replace(body, "https://www.gog.com", "http://localhost:8080/www", -1)
		b = []byte(s)
	}

	return b, nil
}

//https://stackoverflow.com/a/35747382
//https://github.com/gin-gonic/gin#serving-data-from-reader
func processGogResponse(r *http.Response, c *gin.Context) {
	headers := processHeaders(r)
	body, err := processBody(r, headers)

	contentLength := len(body)
	contentType := headers["Content-Type"]
	delete(headers, "Content-Type")
	delete(headers, "Content-Length")
	c.DataFromReader(r.StatusCode, contentLength, contentType, bytes.NewReader(body), headers)
}