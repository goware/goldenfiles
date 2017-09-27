package dump

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Response struct {
	Header           http.Header
	Body             string
	Status           string
	StatusCode       int
	ProtoMajor       int
	ProtoMinor       int
	ContentLength    int64
	TransferEncoding []string
	Uncompressed     bool
	Trailer          http.Header
}

type EncoderDecoder interface {
	Encode(*Response) ([]byte, error)
	Decode([]byte) (*Response, error)
}

type Dump struct {
	ed EncoderDecoder
}

func (d *Dump) Encode(res *http.Response) ([]byte, error) {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	res.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	dump := Response{
		Header:           res.Header,
		Body:             string(buf),
		Status:           res.Status,
		StatusCode:       res.StatusCode,
		ProtoMajor:       res.ProtoMajor,
		ProtoMinor:       res.ProtoMinor,
		ContentLength:    res.ContentLength,
		TransferEncoding: res.TransferEncoding,
		Uncompressed:     res.Uncompressed,
		Trailer:          res.Trailer,
	}

	return d.ed.Encode(&dump)
}

func (d *Dump) Decode(buf []byte) (*http.Response, error) {
	dump, err := d.ed.Decode(buf)
	if err != nil {
		return nil, err
	}

	res := http.Response{
		Header:           dump.Header,
		Body:             ioutil.NopCloser(bytes.NewBufferString(dump.Body)),
		Status:           dump.Status,
		StatusCode:       dump.StatusCode,
		Proto:            fmt.Sprintf("HTTP/%d.%d", dump.ProtoMajor, dump.ProtoMinor),
		ProtoMajor:       dump.ProtoMajor,
		ProtoMinor:       dump.ProtoMinor,
		ContentLength:    dump.ContentLength,
		TransferEncoding: dump.TransferEncoding,
		Uncompressed:     dump.Uncompressed,
		Trailer:          dump.Trailer,
	}

	return &res, nil
}

func New(ed EncoderDecoder) *Dump {
	return &Dump{ed: ed}
}
