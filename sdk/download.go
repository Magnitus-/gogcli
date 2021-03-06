package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strconv"
	"strings"
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
	Size     int64 `xml:"total_size,attr"`
	Chunks   []XmlFileChunk
}

type XmlFileChunk struct {
	Id   int `xml:"id,attr"`
	From int64 `xml:"from,attr"`
	To   int64 `xml:"to,attr"`
	Checksum string `xml:",chardata"`
}

var DOWNLOAD_METADATA_REGEX *regexp.Regexp
func getDownloadMetadataRegex() *regexp.Regexp {
	//ex: <file name="planescape_torment_pl_2.0.0.14.pkg" available="1" notavailablemsg="" md5="4fd4855bc907665c964aebe457dd39eb" chunks="144" timestamp="2016-10-06 11:30:44" total_size="1506711017">
	return regexp.MustCompile(`^<file name="(?P<name>.+)" available="(?:1|0)" notavailablemsg="(?:.*)" md5="(?P<checksum>[0-9a-z]+)" chunks="(?:\d+)" timestamp="(?:.+)" total_size="(?P<size>\d+)">$`)
}

func retrieveDownloadMetadata(c http.Client, metadataUrl string, fn string) (bool, string, string, int64, error, bool) {
	fileInfo := XmlFile{Chunks: make([]XmlFileChunk, 0)}

	r, err := c.Get(metadataUrl)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fn, err.Error())
		return true, "", "", int64(-1), errors.New(msg), false
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		if r.StatusCode == 404 || r.StatusCode == 403 {
			return false, "", "", int64(-1), nil, false
		} else {
			msg := fmt.Sprintf("%s -> Expected response status code of 200, but got %d", fn, r.StatusCode)
			return true, "", "", int64(-1), errors.New(msg), r.StatusCode >= 500
		}
	}

	b, bErr := ioutil.ReadAll(r.Body)
	if bErr != nil {
		msg := fmt.Sprintf("%s -> retrieval body error: %s", fn, bErr.Error())
		return true, "", "", int64(-1), errors.New(msg), false
	}

	err = xml.Unmarshal(b, &fileInfo)
	if err != nil {
		//Fallback to applying a regex on the first line as a last resort
		first_line := strings.Split(string(b), "\n")[0]
		if !DOWNLOAD_METADATA_REGEX.MatchString(first_line) {
			msg := fmt.Sprintf("%s -> Could not parse file xml metadata and the first line was not the expected format: %s", fn, err.Error())
			return true, "", "", int64(-1), errors.New(msg), false
		}

		match := DOWNLOAD_METADATA_REGEX.FindStringSubmatch(first_line)
		size, _ := strconv.ParseInt(match[3], 10, 64)
		return true, match[1], match[2], size, nil, false
	}

	return true, fileInfo.Name, fileInfo.Checksum, fileInfo.Size, nil, false
}

func retrieveUrlRedirectLocation(c http.Client, redirectingUrl string, fn string) (string, error, bool) {
	var location string
	r, err := c.Get(redirectingUrl)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fn, err.Error())
		return "", errors.New(msg), false
	}
	defer r.Body.Close()

	if r.StatusCode != 302 {
		msg := fmt.Sprintf("%s -> Expected response status code of 302, but got %d", fn, r.StatusCode)
		return "", errors.New(msg), r.StatusCode >= 500
	}

	locHeader, ok := r.Header["Location"]
	if !ok {
		msg := fmt.Sprintf("%s -> Expected location header in response, but it was missing", fn)
		return "", errors.New(msg), false
	} else {
		location = locHeader[0]
	}
	return location, nil, false
}

func retrieveUrlContentLength(c http.Client, downloadUrl string, fn string) (int64, error, bool, bool) {
	r, err := c.Head(downloadUrl)
	if err != nil {
		msg := fmt.Sprintf("%s -> retrieval request error: %s", fn, err.Error())
		return int64(0), errors.New(msg), false, false
	}
	defer r.Body.Close()

	if r.StatusCode != 200 {
		msg := fmt.Sprintf("%s -> Expected response status code of 200, but got %d", fn, r.StatusCode)
		return int64(0), errors.New(msg), (r.StatusCode == 404 || r.StatusCode == 403 || r.StatusCode == 301), r.StatusCode >= 500
	}

	clHeader, ok := r.Header["Content-Length"]
	if !ok {
		msg := fmt.Sprintf("%s -> Cannot return exact download size as Content-Length header is not found.", fn)
		return int64(0), errors.New(msg), false, false
	}
	
	length, lErr := strconv.ParseInt(clHeader[0], 10, 64)
	if lErr != nil {
		msg := fmt.Sprintf("%s -> Cannot return exact download size as Content-Length header is not parsable.", fn)
		return int64(0), errors.New(msg), false,false
	}

	return length, nil, false, false
}

//Gets the filename and checksum of the url path, requires 3 requests
func (s *Sdk) GetDownloadFileInfo(downloadPath string) (string, string, int64, error, bool) {
	var filenameLoc string
	var downloadLoc string
	var serverUnavailable bool
	var err error
	var dangling bool
	fn := fmt.Sprintf("GetDownloadFileInfo(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	c := (*s).getClient(false)
	(*s).logger.Debug(fmt.Sprintf("%s -> GET %s", fn, u))

	filenameLoc, err, serverUnavailable = retrieveUrlRedirectLocation(c, u, fn)
	if err != nil {
		if serverUnavailable && (!(*s).maxRetriesReached()) {
			(*s).logger.Warning(fmt.Sprintf("%s -> GET %s failed due to server error. Will retry.", fn, u))
			(*s).incRetries()
			return (*s).GetDownloadFileInfo(downloadPath)
		}
		(*s).resetRetries()
		return "", "", int64(0), err, false
	}
	downloadLoc, err, serverUnavailable = retrieveUrlRedirectLocation(c, filenameLoc, fn)
	if err != nil {
		if serverUnavailable && (!(*s).maxRetriesReached()) {
			(*s).logger.Warning(fmt.Sprintf("%s -> GET %s failed due to server error. Will retry.", fn, u))
			(*s).incRetries()
			return (*s).GetDownloadFileInfo(downloadPath)
		}
		(*s).resetRetries()
		return "", "", int64(0), err, false
	}

	//Finally, retrieve the metadata
	metadataUrl, metadataUrlErr := convertDownloadUrlToMetadataUrl(downloadLoc)
	if metadataUrlErr != nil {
		(*s).resetRetries()
		return "", "", int64(0), metadataUrlErr, false
	}

	found, filename, checksum, size, retrieveMetaErr, serverUnavailable := retrieveDownloadMetadata(c, metadataUrl, fn)
	if retrieveMetaErr != nil {
		if serverUnavailable && (!(*s).maxRetriesReached()) {
			(*s).logger.Warning(fmt.Sprintf("%s -> GET %s failed due to server error. Will retry.", fn, u))
			(*s).incRetries()
			return (*s).GetDownloadFileInfo(downloadPath)
		}
		(*s).resetRetries()
		return "", "", int64(0), retrieveMetaErr, false
	}
	if !found {
		filename, err = getFilenameFromUrl(filenameLoc, fn)
		if err != nil {
			(*s).resetRetries()
			return "", "", int64(0), err, false
		}

		size, err, dangling, serverUnavailable = retrieveUrlContentLength(c, downloadLoc, fn)
		if err != nil {
			if serverUnavailable && (!(*s).maxRetriesReached()) {
				(*s).logger.Warning(fmt.Sprintf("%s -> GET %s failed due to server error. Will retry.", fn, u))
				(*s).incRetries()
				return (*s).GetDownloadFileInfo(downloadPath)
			}
			(*s).resetRetries()
			return "", "", int64(0), err, dangling
		}

		(*s).resetRetries()
		return filename, "", size, nil, false
	} else {
		(*s).resetRetries()
		return filename, checksum, size, nil, false
	}
}

//Gets just the filename of the url path, requires 2 requests
func (s *Sdk) GetDownloadFilename(downloadPath string) (string, error) {
	fn := fmt.Sprintf("GetDownloadFilename(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	c := (*s).getClient(false)
	(*s).logger.Debug(fmt.Sprintf("%s -> GET %s", fn, u))

	redirectLocation, err, serverUnavailable := retrieveUrlRedirectLocation(c, u, fn)
	if err != nil {
		if serverUnavailable && (!(*s).maxRetriesReached()) {
			(*s).logger.Warning(fmt.Sprintf("%s -> GET %s failed due to server error. Will retry.", fn, u))
			(*s).incRetries()
			return (*s).GetDownloadFilename(downloadPath)
		}
		(*s).resetRetries()
		return "", err
	}

	(*s).resetRetries()

	name, err := getFilenameFromUrl(redirectLocation, fn)
	return name, err
}

func (s *Sdk) GetDownloadHandle(downloadPath string) (io.ReadCloser, int64, string, error) {
	fn := fmt.Sprintf("GetDownloadHandle(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	var body io.ReadCloser
	bodyLength := int64(0)
	filename := ""

	c := (*s).getClient(true)
	(*s).logger.Debug(fmt.Sprintf("%s -> GET %s", fn, u))

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
			(*s).logger.Warning(fmt.Sprintf("%s -> Cannot return exact download size as Content-Length header is not parsable. Will set it to 0.", fn))
		} else {
			bodyLength = l
		}
	} else {
		(*s).logger.Warning(fmt.Sprintf("%s -> Cannot return exact download size as Content-Length header is not found. Will set it to 0.", fn))
	}

	finalURL := r.Request.URL.String()
	(*s).logger.Debug(fmt.Sprintf("    Final Url: %s", finalURL))

	p := (*r.Request.URL).Path
	filename = path.Base(p)

	return body, bodyLength, filename, nil
}
