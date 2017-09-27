package gob

import (
	"bytes"
	"encoding/gob"

	"github.com/goware/mockingbird/dump"
)

type Gob struct {
}

func (g *Gob) Encode(res *dump.Response) ([]byte, error) {
	var out bytes.Buffer
	enc := gob.NewEncoder(&out)
	if err := enc.Encode(res); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func (g *Gob) Decode(buf []byte) (*dump.Response, error) {
	var res dump.Response

	in := bytes.NewBuffer(buf)
	dec := gob.NewDecoder(in)

	if err := dec.Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

var _ = dump.EncoderDecoder(&Gob{})
