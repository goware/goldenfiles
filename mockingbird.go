package mockingbird

import (
	"errors"
	"net/http"

	"github.com/goware/mockingbird/store"
)

func New(s store.Store) *http.Client {
	if s == nil {
		s = &store.Mem{}
	}
	return &http.Client{
		Transport: &roundTripper{
			store: s,
		},
	}
}

func Remove(c *http.Client, req *http.Request) error {
	tr, ok := c.Transport.(*roundTripper)
	if !ok {
		return errors.New("missing roundtripper")
	}

	return tr.remove(req)
}

func Record(c *http.Client, req *http.Request, res *http.Response) error {
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

	tr, ok := c.Transport.(*roundTripper)
	if !ok {
		return errors.New("missing roundtripper")
	}

	return tr.set(req, res)
}
