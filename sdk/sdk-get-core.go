package sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
)

type BodyReaderReply struct {
	BodyHandle io.ReadCloser
	BodyLength int64
	FinalUrl string
	StatusCode int
	RetriesLeft int64
}

func (s *Sdk) getUrlBodyReader(url string, fnCall string, retriesLeft int64) (BodyReaderReply, error) {
	c := s.getClient(true)

	(*s).logger.Debug(fmt.Sprintf("%s -> GET %s", fnCall, url))

	r, err := c.Get(url)
	if err != nil {
		if retriesLeft > 0 {
			(*s).logger.Warning(fmt.Sprintf("%s -> failed with retrieval request error %s. Will retry.", fnCall, err.Error()))
			s.pauseAfterError()
			return s.getUrlBodyReader(url, fnCall, retriesLeft - 1)
		}

		msg := fmt.Sprintf("%s -> retrieval request error: %s", fnCall, err.Error())
		return BodyReaderReply{
			BodyHandle: nil,
			BodyLength: int64(-1),
			FinalUrl: "",
			StatusCode: -1,
			RetriesLeft: retriesLeft,
		}, errors.New(msg)
	}

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		if r.StatusCode >= 500 && retriesLeft > 0 {
			r.Body.Close()
			(*s).logger.Warning(fmt.Sprintf("%s -> failed with code %d. Will retry.", fnCall, r.StatusCode))
			s.pauseAfterError()
			return s.getUrlBodyReader(url, fnCall, retriesLeft - 1)
		}
		msg := fmt.Sprintf("%s -> body download handle retrieval error: did not expect status code of %d", fnCall, r.StatusCode)
		return BodyReaderReply{
			BodyHandle: r.Body,
			BodyLength: int64(-1),
			FinalUrl: r.Request.URL.String(),
			StatusCode: r.StatusCode,
			RetriesLeft: retriesLeft,
		}, errors.New(msg)
	}

	bodyLength := int64(-1)
	var lErr error

	clHeader, ok := r.Header["Content-Length"]
	if ok {
		bodyLength, lErr = strconv.ParseInt(clHeader[0], 10, 64)
		if lErr != nil {
			msg := fmt.Sprintf("%s -> Cannot return exact download size as Content-Length header is not parsable.", fnCall)
			return BodyReaderReply{
				BodyHandle: r.Body,
				BodyLength: int64(-1),
				FinalUrl: r.Request.URL.String(),
				StatusCode: r.StatusCode,
				RetriesLeft: retriesLeft,
			}, errors.New(msg)
		}
	}

	return BodyReaderReply{
		BodyHandle: r.Body,
		BodyLength: bodyLength,
		FinalUrl: r.Request.URL.String(),
		StatusCode: r.StatusCode,
		RetriesLeft: retriesLeft,
	}, nil
}

type BodyChecksumReply struct {
	BodyChecksum string
	BodyLength int64
	FinalUrl string
	StatusCode int
	RetriesLeft int64
}

func (s *Sdk) getUrlBodyChecksum(url string, fnCall string, retriesLeft int64) (BodyChecksumReply, error) {
	reply, err := s.getUrlBodyReader(url, fnCall, retriesLeft)
	if reply.BodyHandle != nil {
		defer reply.BodyHandle.Close()
	}
	if err != nil {
		return BodyChecksumReply{
			BodyChecksum: "",
			BodyLength: reply.BodyLength,
			FinalUrl: reply.FinalUrl,
			StatusCode: reply.StatusCode,
			RetriesLeft: reply.RetriesLeft,
		}, err
	}

	h := md5.New()
	copiedAmount, copyErr := io.Copy(h, reply.BodyHandle)
	if copiedAmount != reply.BodyLength || copyErr != nil {
		if copyErr == nil {
			copyErr = errors.New(fmt.Sprintf("%s -> checksum computation processed %d bytes and expected %d", fnCall, copiedAmount, reply.BodyLength))
		}
		if reply.RetriesLeft > 0 {
			(*s).logger.Warning(fmt.Sprintf("%s -> checksum computation failed with error: %s. Will retry.", fnCall, copyErr.Error()))
			s.pauseAfterError()
			return s.getUrlBodyChecksum(url, fnCall, retriesLeft - 1)
		}
		msg := fmt.Sprintf("%s -> checksum computation failed with error: %s", fnCall, copyErr.Error())
		return BodyChecksumReply{
			BodyChecksum: "",
			BodyLength: reply.BodyLength,
			FinalUrl: reply.FinalUrl,
			StatusCode: reply.StatusCode,
			RetriesLeft: reply.RetriesLeft,
		}, errors.New(msg)
	}

	return BodyChecksumReply{
		BodyChecksum: hex.EncodeToString(h.Sum(nil)),
		BodyLength: reply.BodyLength,
		FinalUrl: reply.FinalUrl,
		StatusCode: reply.StatusCode,
		RetriesLeft: reply.RetriesLeft,
	}, nil
}

type BodyLengthReply struct {
	BodyLength int64
	FinalUrl string
	StatusCode int
	RetriesLeft int64
}

func (s *Sdk) getUrlBodyLength(url string, fnCall string, retriesLeft int64) (BodyLengthReply, error) {
	c := s.getClient(true)
	r, err := c.Head(url)
	if err != nil {
		if retriesLeft > 0 {
			(*s).logger.Warning(fmt.Sprintf("%s -> content length retrieval error: %s. Will retry.", fnCall, err.Error()))
			s.pauseAfterError()
			return s.getUrlBodyLength(url, fnCall, retriesLeft - 1)
		}
		msg := fmt.Sprintf("%s -> content length retrieval request error: %s", fnCall, err.Error())
		return BodyLengthReply{
			BodyLength: int64(-1),
			FinalUrl: "",
			StatusCode: -1,
			RetriesLeft: retriesLeft,
		}, errors.New(msg)
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		if r.StatusCode >= 500 && retriesLeft > 0 {
			(*s).logger.Warning(fmt.Sprintf("%s -> content length retrieval failed with code %d. Will retry.", fnCall, r.StatusCode))
			s.pauseAfterError()
			return s.getUrlBodyLength(url, fnCall, retriesLeft - 1)
		}
		msg := fmt.Sprintf("%s -> content length retrieval error: did not expect status code of %d", fnCall, r.StatusCode)
		return BodyLengthReply{
			BodyLength: int64(-1),
			FinalUrl: r.Request.URL.String(),
			StatusCode: r.StatusCode,
			RetriesLeft: retriesLeft,
		}, errors.New(msg)
	}

	clHeader, ok := r.Header["Content-Length"]
	if !ok {
		msg := fmt.Sprintf("%s -> content length retrieval error: cannot return body length as Content-Length header is not found.", fnCall)
		return BodyLengthReply{
			BodyLength: int64(-1),
			FinalUrl: r.Request.URL.String(),
			StatusCode: r.StatusCode,
			RetriesLeft: retriesLeft,
		}, errors.New(msg)
	}

	length, lErr := strconv.ParseInt(clHeader[0], 10, 64)
	if lErr != nil {
		msg := fmt.Sprintf("%s -> content length retrieval error: cannot return body length as Content-Length header is not parsable.", fnCall)
		return BodyLengthReply{
			BodyLength: int64(-1),
			FinalUrl: r.Request.URL.String(),
			StatusCode: r.StatusCode,
			RetriesLeft: retriesLeft,
		}, errors.New(msg)
	}

	return BodyLengthReply{
		BodyLength: length,
		FinalUrl: r.Request.URL.String(),
		StatusCode: r.StatusCode,
		RetriesLeft: retriesLeft,
	}, nil
}

type BodyReply struct {
	Body []byte
	FinalUrl string
	StatusCode int
	RetriesLeft int64
}

func (s *Sdk) getUrlBody(url string, fnCall string, jsonBody bool, retriesLeft int64) (BodyReply, error) {
	(*s).logger.Debug(fmt.Sprintf("%s -> GET %s", fnCall, url))

	reply, err := s.getUrlBodyReader(url, fnCall, retriesLeft)
	if reply.BodyHandle != nil {
		defer reply.BodyHandle.Close()
	}
	if err != nil {
		return BodyReply{
			Body: nil,
			FinalUrl: reply.FinalUrl,
			StatusCode: reply.StatusCode,
			RetriesLeft: reply.RetriesLeft,
		}, err
	}

	b, bErr := ioutil.ReadAll(reply.BodyHandle)
	if bErr != nil {
		if reply.RetriesLeft > 0 {
			(*s).logger.Warning(fmt.Sprintf("%s -> body retrieval error: %s. Will retry.", fnCall, bErr.Error()))
			s.pauseAfterError()
			return s.getUrlBody(url, fnCall, jsonBody, reply.RetriesLeft - 1)
		}
		msg := fmt.Sprintf("%s -> body retrieval error: %s", fnCall, bErr.Error())
		return BodyReply{
			Body: nil,
			FinalUrl: reply.FinalUrl,
			StatusCode: reply.StatusCode,
			RetriesLeft: reply.RetriesLeft,
		}, errors.New(msg)
	}

	if jsonBody {
		var out bytes.Buffer
		jErr := json.Indent(&out, b, "", "  ")
		if jErr != nil {
			msg := fmt.Sprintf("%s -> json parsing error: %s", fnCall, jErr.Error())
			return BodyReply{
				Body: nil,
				FinalUrl: reply.FinalUrl,
				StatusCode: reply.StatusCode,
				RetriesLeft: reply.RetriesLeft,
			}, errors.New(msg)
		}
		b = out.Bytes()
	}
	(*s).logger.Debug(fmt.Sprintf("%s -> response body: %s", fnCall, string(b)))

	return BodyReply{
		Body: b,
		FinalUrl: reply.FinalUrl,
		StatusCode: reply.StatusCode,
		RetriesLeft: reply.RetriesLeft,
	}, nil
}

type RedirectReply struct {
	RedirectUrl string
	StatusCode int
	RetriesLeft int64
}

func (s *Sdk) getUrlRedirect(url string, fnCall string, retriesLeft int64) (RedirectReply, error) {
	c := s.getClient(false)
	
	var location string
	r, err := c.Get(url)
	if err != nil {
		if retriesLeft > 0 {
			(*s).logger.Warning(fmt.Sprintf("%s -> redirect retrieval error: %s. Will retry.", fnCall, err.Error()))
			s.pauseAfterError()
			return s.getUrlRedirect(url, fnCall, retriesLeft - 1)
		}
		msg := fmt.Sprintf("%s -> redirect retrieval error: %s", fnCall, err.Error())
		return RedirectReply{
			RedirectUrl: "",
			StatusCode: -1,
			RetriesLeft:  retriesLeft,
		}, errors.New(msg)
	}
	defer r.Body.Close()

	if r.StatusCode < 300 || r.StatusCode >= 400 {
		if r.StatusCode >= 500 && retriesLeft > 0 {
			(*s).logger.Warning(fmt.Sprintf("%s -> redirect retrieval error: expected response status code of 3xx, but got %d. Will retry.", fnCall, r.StatusCode))
			s.pauseAfterError()
			return s.getUrlRedirect(url, fnCall, retriesLeft - 1)
		}

		msg := fmt.Sprintf("%s -> redirect retrieval error: expected response status code of 3xx, but got %d", fnCall, r.StatusCode)
		return RedirectReply{
			RedirectUrl: "",
			StatusCode: r.StatusCode,
			RetriesLeft:  retriesLeft,
		}, errors.New(msg)
	}

	locHeader, ok := r.Header["Location"]
	if !ok {
		msg := fmt.Sprintf("%s -> redirect retrieval error: expected location header in response, but it was missing", fnCall)
		return RedirectReply{
			RedirectUrl: "",
			StatusCode: r.StatusCode,
			RetriesLeft:  retriesLeft,
		}, errors.New(msg)
	} else {
		location = locHeader[0]
	}

	return RedirectReply{
		RedirectUrl: location,
		StatusCode: r.StatusCode,
		RetriesLeft:  retriesLeft,
	}, nil
}