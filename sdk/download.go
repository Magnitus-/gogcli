package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"
	"net/url"
)

func getFilenameFromUrl(location string, fn string) (string, error) {
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

func convertDownloadUrlToMetadataUrl(downloadUrl string) (string, error) {
	parsedUrl, err := url.Parse(downloadUrl)
	if err != nil {
		msg := fmt.Sprintf("convertDownloadToMetadata(downloadUrl=%s) -> Could not parse url", downloadUrl)
		return "", errors.New(msg)
	}
	(*parsedUrl).Path = (*parsedUrl).Path + ".xml"
	return parsedUrl.String(), nil
}

type XmlFile struct {
	XMLName  xml.Name `xml:"file"`
	Name     string `xml:"name,attr"`
	Checksum string `xml:"md5,attr"`
	Size     int `xml:"total_size,attr"`
	Chunks   []XmlFileChunk
}

type XmlFileChunk struct {
	Id   int `xml:"id,attr"`
	From int `xml:"from,attr"`
	To   int `xml:"to,attr"`
	Checksum string `xml:",chardata"`
}

func retrieveDownloadMetadata(c http.Client, metadataUrl string, fn string) (bool, string, string, int, error) {
	fileInfo := XmlFile{Chunks: make([]XmlFileChunk, 0)}

	r, err := c.Get(metadataUrl)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fn, err.Error())
		return true, "", "", -1, errors.New(msg)
	}
	if r.StatusCode != 200 {
		if r.StatusCode == 404 {
			return false, "", "", -1, nil
		} else {
			msg := fmt.Sprintf("%s -> Expected response status code of 200, but got %d", fn, r.StatusCode)
			return true, "", "", -1, errors.New(msg)
		}
	}

	b, bErr := ioutil.ReadAll(r.Body)
	if bErr != nil {
		msg := fmt.Sprintf("%s -> retrieval body error: %s", fn, bErr.Error())
		return true, "", "", -1, errors.New(msg)
	}

	err = xml.Unmarshal(b, &fileInfo)
	if err != nil {
		msg := fmt.Sprintf("%s -> Could not parse file metadata: %s", fn, err.Error())
		return true, "", "", -1, errors.New(msg)
	}

	return true, fileInfo.Name, fileInfo.Checksum, fileInfo.Size, nil
}

func retrieveUrlRedirectLocation(c http.Client, redirectingUrl string, fn string) (string, error) {
	var location string
	r, err := c.Get(redirectingUrl)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fn, err.Error())
		return "", errors.New(msg)
	}

	if r.StatusCode != 302 {
		msg := fmt.Sprintf("%s -> Expected response status code of 302, but got %d", fn, r.StatusCode)
		return "", errors.New(msg)
	}

	locHeader, ok := r.Header["Location"]
	if !ok {
		msg := fmt.Sprintf("%s -> Expected location header in response, but it was missing", fn)
		return "", errors.New(msg)
	} else {
		location = locHeader[0]
	}
	return location, nil
}

//Gets the filename and checksum of the url path, requires 3 requests
func (s *Sdk) GetDownloadInfo(downloadPath string) (string, string, int, error) {
	var filenameLoc string
	var downloadLoc string
	var err error
	fn := fmt.Sprintf("GetDownloadInfo(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	c := (*s).getClient(false)
	if (*s).debug {
		(*s).logger.Println(fmt.Sprintf("%s -> GET %s", fn, u))
	}

	filenameLoc, err = retrieveUrlRedirectLocation(c, u, fn)
	if err != nil {
		return "", "", -1, err
	}
	downloadLoc, err = retrieveUrlRedirectLocation(c, filenameLoc, fn)
	if err != nil {
		return "", "", -1, err
	}

	//Finally, retrieve the metadata
	metadataUrl, metadataUrlErr := convertDownloadUrlToMetadataUrl(downloadLoc)
	if metadataUrlErr != nil {
		return "", "", -1, metadataUrlErr
	}

	found, filename, checksum, size, retrieveMetaErr := retrieveDownloadMetadata(c, metadataUrl, fn)
	if retrieveMetaErr != nil {
		return "", "", -1, retrieveMetaErr
	}
	if !found {
		filename, err = getFilenameFromUrl(filenameLoc, fn)
		if err != nil {
			return "", "", -1, err
		}
		return filename, "", -1, nil
	} else {
		return filename, checksum, size, nil
	}
}

//Gets just the filename of the url path, requires 2 requests
func (s *Sdk) GetDownloadFilename(downloadPath string) (string, error) {
	fn := fmt.Sprintf("GetDownloadFilename(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	c := (*s).getClient(false)
	if (*s).debug {
		(*s).logger.Println(fmt.Sprintf("%s -> GET %s", fn, u))
	}

	redirectLocation, err := retrieveUrlRedirectLocation(c, u, fn)
	if err != nil {
		return "", err
	}

	return getFilenameFromUrl(redirectLocation, fn)
}

func (s *Sdk) GetDownloadHandle(downloadPath string) (io.ReadCloser, int64, string, error) {
	fn := fmt.Sprintf("GetDownloadHandle(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	var body io.ReadCloser
	bodyLength := int64(0)
	filename := ""

	c := (*s).getClient(true)
	if (*s).debug {
		(*s).logger.Println(fmt.Sprintf("%s -> GET %s", fn, u))
	}

	r, err := c.Get(u)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fn, err.Error())
		return nil, int64(0), "", errors.New(msg)
	}

	body = r.Body

	clHeader, ok := r.Header["Content-Length"]
	if ok {
		l, lErr := strconv.ParseInt(clHeader[0], 10, 64)
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
