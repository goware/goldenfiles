package goldenfiles

import (
	"net/http"

	"github.com/goware/goldenfiles/dump"
	gobd "github.com/goware/goldenfiles/dump/gob"
	mems "github.com/goware/goldenfiles/store/mem"
)

func NewTransport() *Transport {
	return &Transport{
		store:          &mems.Mem{},
		encoderDecoder: dump.New(&gobd.Gob{}),
	}
}

func NewClient() *http.Client {
	return &http.Client{
		Transport: NewTransport(),
	}
}
