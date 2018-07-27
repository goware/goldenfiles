package goldenfiles

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/goware/goldenfiles/dump"
	"github.com/goware/goldenfiles/store"
)

type Transport struct {
	store          store.Store
	encoderDecoder *dump.Dump
}

func (tr *Transport) SetEncoderDecoder(ed dump.EncoderDecoder) {
	tr.encoderDecoder = dump.New(ed)
}

func (tr *Transport) SetStore(s store.Store) {
	tr.store = s
}

func (tr *Transport) get(req *http.Request) (*http.Response, error) {
	id, err := tr.requestID(req)
	if err != nil {
		return nil, err
	}

	buf, err := tr.store.Get(id)
	if err != nil {
		return nil, err
	}

	return tr.encoderDecoder.Decode(buf)
}

func (tr *Transport) Remove(req *http.Request) error {
	id, err := tr.requestID(req)
	if err != nil {
		return err
	}
	return tr.store.Delete(id)
}

func (tr *Transport) set(req *http.Request, res *http.Response) error {
	id, err := tr.requestID(req)
	if err != nil {
		return err
	}

	var bodyReader = io.Reader(res.Body)

	switch res.Header.Get("Content-Encoding") {
	case "deflate":
		z, err := zlib.NewReader(bodyReader)
		if err != nil {
			return err
		}
		res.Header.Del("Content-Encoding")

		defer z.Close()
		bodyReader = io.Reader(z)
	case "gzip":
		gz, err := gzip.NewReader(bodyReader)
		if err != nil {
			return err
		}
		res.Header.Del("Content-Encoding")

		defer gz.Close()
		bodyReader = io.Reader(gz)
	}

	var body []byte
	body, err = ioutil.ReadAll(bodyReader)
	if err != nil {
		return err
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	buf, err := tr.encoderDecoder.Encode(res)
	if err != nil {
		return err
	}

	return tr.store.Set(id, buf)
}

func (tr *Transport) requestID(req *http.Request) (string, error) {
	buf := []byte{}

	if err := req.Context().Err(); err != nil {
		return "", err
	}

	if req.Body != nil {
		var err error
		buf, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return "", err
		}
		req.Body.Close()

		req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	}

	return fmt.Sprintf(
		"%s_%s",
		req.Method,
		strings.Trim(req.URL.String(), "/"),
	), nil
}

func (tr *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	t := &http.Transport{
		DisableCompression: true,
	}

	res, err := tr.get(req)
	if err == nil {
		res.Request = req
		return res, nil
	}

	// Get data
	res, err = t.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if err := tr.set(req, res); err != nil {
		return nil, err
	}

	return res, nil
}

func (tr *Transport) Record(req *http.Request, res *http.Response) error {
	if res.StatusCode == 0 {
		res.StatusCode = http.StatusOK
	}
	if res.Status == "" {
		res.Status = "200 OK"
	}
	if res.Proto == "" {
		res.Proto = "HTTP/1.0"
	}
	if res.ProtoMajor == 0 {
		res.ProtoMajor = 1
	}
	if res.Header == nil {
		res.Header = http.Header{
			"Content-Type": {"text/plain"},
		}
	}

	return tr.set(req, res)
}
