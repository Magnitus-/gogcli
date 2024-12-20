package sdk

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
)

type XmlFile struct {
	XMLName  xml.Name `xml:"file"`
	Name     string   `xml:"name,attr"`
	Checksum string   `xml:"md5,attr"`
	Size     int64    `xml:"total_size,attr"`
	Chunks   []XmlFileChunk
}

type XmlFileChunk struct {
	Id       int    `xml:"id,attr"`
	From     int64  `xml:"from,attr"`
	To       int64  `xml:"to,attr"`
	Checksum string `xml:",chardata"`
}

var DOWNLOAD_METADATA_REGEX *regexp.Regexp

func getDownloadMetadataRegex() *regexp.Regexp {
	//ex: <file name="planescape_torment_pl_2.0.0.14.pkg" available="1" notavailablemsg="" md5="4fd4855bc907665c964aebe457dd39eb" chunks="144" timestamp="2016-10-06 11:30:44" total_size="1506711017">
	return regexp.MustCompile(`^<file name="(?P<name>.+)" available="(?:1|0)" notavailablemsg="(?:.*)" md5="(?P<checksum>[0-9a-z]+)" chunks="(?:\d+)" timestamp="(?:.+)" total_size="(?P<size>\d+)">$`)
}

type DownloadMetadata struct {
	Filename    string
	Checksum    string
	Size        int64
	Found       bool
	BadMetadata bool
}

func (s *Sdk) retrieveDownloadMetadata(metadataUrl string, fn string, retriesLeft int64) (DownloadMetadata, error) {
	fileInfo := XmlFile{Chunks: make([]XmlFileChunk, 0)}

	reply, err := s.getUrlBody(
		metadataUrl,
		fn,
		false,
		retriesLeft,
	)
	if err != nil {
		return DownloadMetadata{
			Filename: "",
			Checksum: "",
			Size: int64(-1),
			Found: reply.StatusCode != 403 && reply.StatusCode != 404,
			BadMetadata: reply.StatusCode != 403 && reply.StatusCode != 404,
		}, err
	}

	err = xml.Unmarshal(reply.Body, &fileInfo)
	if err != nil {
		//Fallback to applying a regex on the first line as a last resort
		first_line := strings.Split(string(reply.Body), "\n")[0]
		if !DOWNLOAD_METADATA_REGEX.MatchString(first_line) {
			msg := fmt.Sprintf("%s -> Could not parse file xml metadata and the first line was not the expected format: %s", fn, err.Error())
			return DownloadMetadata{
				Filename: "",
				Checksum: "",
				Size: int64(-1),
				Found: true,
				BadMetadata: true,
			}, errors.New(msg)
		}

		match := DOWNLOAD_METADATA_REGEX.FindStringSubmatch(first_line)
		size, _ := strconv.ParseInt(match[3], 10, 64)
		return DownloadMetadata{
			Filename: match[1],
			Checksum: match[2],
			Size: size,
			Found: true,
			BadMetadata: false,
		}, nil
	}

	return DownloadMetadata{
		Filename: fileInfo.Name,
		Checksum: fileInfo.Checksum,
		Size: fileInfo.Size,
		Found: true,
		BadMetadata: false,
	}, nil
}

func getFilenameFromUrlPath(location string, fn string) (string, error) {
	locUrl, err := url.Parse(location)
	if err != nil {
		msg := fmt.Sprintf("%s -> Error parsing location header url: %s", fn, err.Error())
		return "", errors.New(msg)
	}

	return path.Base(locUrl.Path), nil
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

//Gets the filename and checksum of the url path, requires 3 requests
func (s *Sdk) getDownloadFileInfo(downloadPath string) (string, string, int64, error, bool, bool) {
	fn := fmt.Sprintf("getDownloadFileInfo(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	(*s).logger.Debug(fmt.Sprintf("%s -> GET %s", fn, u))

	//Redirection
	reply, err := s.getUrlRedirect(u, fn, (*s).maxRetries)
	if err != nil {
		return "", "", int64(-1), err, reply.StatusCode == 403 || reply.StatusCode == 404, reply.StatusCode != 403 && reply.StatusCode != 404
	}
	filenameLoc := reply.RedirectUrl

	//Convert final download url to metadata url
	metadataUrl, metadataUrlErr := convertDownloadUrlToMetadataUrl(filenameLoc)
	if metadataUrlErr != nil {
		return "", "", int64(-1), metadataUrlErr, false, true
	}

	//Finally, retrieve the metadata
	metadata, metadataErr := s.retrieveDownloadMetadata(metadataUrl, fn, (*s).maxRetries)
	if metadataErr != nil && metadata.Found {
		return "", "", int64(-1), metadataErr, !metadata.Found, metadata.BadMetadata
	}
	
	if !metadata.Found {
		var filename string
		filename, err = getFilenameFromUrlPath(filenameLoc, fn)
		if err != nil {
			return "", "", int64(-1), err, false, true
		}

		lengthReply, lengthErr := s.getUrlBodyLength(filenameLoc, fn, (*s).maxRetries)
		if lengthErr != nil {
			return "", "", int64(-1), lengthErr, false, true
		}

		return filename, "", lengthReply.BodyLength, nil, false, false
	}

	return metadata.Filename, metadata.Checksum, metadata.Size, nil, false, false
}

func (s *Sdk) getDownloadFileInfoWorkaroundWay(downloadPath string) (string, string, int64, error) {
	fn := fmt.Sprintf(" getDownloadFileInfoWorkaroundWay(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	(*s).logger.Debug(fmt.Sprintf("%s -> GET %s", fn, u))

	reply, err := s.getUrlBodyChecksum(u, fn, (*s).maxRetries)
	if err != nil {
		return  "", "", int64(0), err
	}

	finalUrl, _ := url.Parse(reply.FinalUrl)
	filename := path.Base(finalUrl.Path)

	return filename, reply.BodyChecksum, reply.BodyLength, nil
}

type GetFileInfoReturn struct {
	Url         string
	Name        string
	Checksum    string
	Size        int64
	Error       error
	Dangling    bool
	BadMetadata bool
}

//Calls getDownloadFileInfo and fallbacks to getDownloadFileInfoWorkaroundWay if tolerateBadMetadata is true
func (s *Sdk) GetFileInfo(downloadPath string, tolerateBadMetadata bool) GetFileInfoReturn {
	name, checksum, size, err, dangling, badMetadata := s.getDownloadFileInfo(downloadPath)
	if badMetadata && tolerateBadMetadata {
		var workaroundErr error
		name, checksum, size, workaroundErr = s.getDownloadFileInfoWorkaroundWay(downloadPath)
		if workaroundErr != nil {
			return GetFileInfoReturn{
				Url: downloadPath,
				Name: name,
				Checksum: checksum,
				Size: size,
				Error: workaroundErr,
				Dangling: false,
				BadMetadata: false,
			}
		}

		return GetFileInfoReturn{
			Url: downloadPath,
			Name: name,
			Checksum: checksum,
			Size: size,
			Error: err,
			Dangling: false,
			BadMetadata: true,
		}
	}

	return GetFileInfoReturn{
		Url: downloadPath, 
		Name: name,
		Checksum: checksum,
		Size: size,
		Error: err,
		Dangling: dangling,
		BadMetadata: badMetadata,
	}
}

//Gets just the filename of the url path, requires 2 requests
func (s *Sdk) GetDownloadFilename(downloadPath string) (string, error) {
	fn := fmt.Sprintf("GetDownloadFilename(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	(*s).logger.Debug(fmt.Sprintf("%s -> GET %s", fn, u))
	
	reply, err := s.getUrlRedirect(u, fn, (*s).maxRetries)
	if err != nil {
		return "", err
	}

	name, err := getFilenameFromUrlPath(reply.RedirectUrl, fn)
	return name, err
}

func (s *Sdk) GetDownloadHandle(downloadPath string) (io.ReadCloser, int64, string, error) {
	fn := fmt.Sprintf("GetDownloadHandle(downloadPath=%s)", downloadPath)
	u := fmt.Sprintf("https://www.gog.com%s", downloadPath)

	reply, err := s.getUrlBodyReader(u, fn, (*s).maxRetries)
	if err != nil {
		return nil, int64(0), "", err
	}
	
	(*s).logger.Debug(fmt.Sprintf("Final Url: %s", reply.FinalUrl))

	downloadUrl, _ := url.Parse(reply.FinalUrl)
	filename := path.Base(downloadUrl.Path)

	return reply.BodyHandle, reply.BodyLength, filename, nil
}