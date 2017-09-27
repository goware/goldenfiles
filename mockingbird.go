package mockingbird

import (
	"net/http"

	"github.com/goware/mockingbird/dump"
	gobd "github.com/goware/mockingbird/dump/gob"
	mems "github.com/goware/mockingbird/store/mem"
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
