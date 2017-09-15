package mockingbird

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"net/http"
)

type responseDump struct {
	Header           http.Header
	Body             []byte
	Status           string
	StatusCode       int
	ProtoMajor       int
	ProtoMinor       int
	ContentLength    int64
	TransferEncoding []string
	Uncompressed     bool
	Trailer          http.Header
}

func serializeResponse(res *http.Response) ([]byte, error) {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Silently replacing body
	res.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	dump := responseDump{
		Header:           res.Header,
		Body:             buf,
		Status:           res.Status,
		StatusCode:       res.StatusCode,
		ProtoMajor:       res.ProtoMajor,
		ProtoMinor:       res.ProtoMinor,
		ContentLength:    res.ContentLength,
		TransferEncoding: res.TransferEncoding,
		Uncompressed:     res.Uncompressed,
		Trailer:          res.Trailer,
	}

	var out bytes.Buffer
	enc := gob.NewEncoder(&out)
	if err := enc.Encode(dump); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func unserializeResponse(buf []byte) (*http.Response, error) {
	dump := responseDump{}

	in := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(in)

	if err := dec.Decode(&dump); err != nil {
		return nil, err
	}

	res := http.Response{
		Header:           dump.Header,
		Body:             ioutil.NopCloser(bytes.NewBuffer(dump.Body)),
		Status:           dump.Status,
		StatusCode:       dump.StatusCode,
		ProtoMajor:       dump.ProtoMajor,
		ProtoMinor:       dump.ProtoMinor,
		ContentLength:    dump.ContentLength,
		TransferEncoding: dump.TransferEncoding,
		Uncompressed:     dump.Uncompressed,
		Trailer:          dump.Trailer,
	}

	res.Body = ioutil.NopCloser(bytes.NewBuffer(dump.Body))

	return &res, nil
}
