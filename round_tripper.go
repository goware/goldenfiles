package mockingbird

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/goware/mockingbird/store"
)

type roundTripper struct {
	store store.Store
}

func (rt *roundTripper) get(req *http.Request) (*http.Response, error) {
	id, err := rt.requestID(req)
	if err != nil {
		return nil, err
	}

	buf, err := rt.store.Get(id)
	if err != nil {
		return nil, err
	}

	return unserializeResponse(buf)
}

func (rt *roundTripper) remove(req *http.Request) error {
	id, err := rt.requestID(req)
	if err != nil {
		return err
	}
	return rt.store.Delete(id)
}

func (rt *roundTripper) set(req *http.Request, res *http.Response) error {
	id, err := rt.requestID(req)
	if err != nil {
		return err
	}

	buf, err := serializeResponse(res)
	if err != nil {
		return err
	}

	return rt.store.Set(id, buf)
}

func (rt *roundTripper) requestID(req *http.Request) (string, error) {
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
		"%s:%s:%s",
		req.Method,
		req.URL.String(),
		fmt.Sprintf("%x", sha256.Sum256(buf)),
	), nil
}

func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	t := &http.Transport{}

	res, err := rt.get(req)
	if err == nil {
		res.Request = req
		return res, nil
	}

	res, err = t.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	if err := rt.set(req, res); err != nil {
		return nil, err
	}

	return res, nil
}
