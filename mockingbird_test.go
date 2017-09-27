package mockingbird

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/goware/mockingbird/dump"
	"github.com/goware/mockingbird/dump/gob"
	"github.com/goware/mockingbird/dump/json"
	"github.com/goware/mockingbird/dump/wire"
	"github.com/goware/mockingbird/store"
	"github.com/goware/mockingbird/store/bolt"
	"github.com/goware/mockingbird/store/file"
	"github.com/goware/mockingbird/store/mem"
	"github.com/stretchr/testify/assert"
)

func TestClientCachedResponse(t *testing.T) {
	client := http.Client{Transport: NewTransport()}

	for i := 0; i < 5; i++ {
		_, err := client.Get("https://golang.org/pkg/net/http/")
		assert.NoError(t, err)
	}
}

func TestClientMockCustomBody(t *testing.T) {
	tr := NewTransport()

	client := http.Client{Transport: tr}

	helloWorld := "Hello world!"

	req, err := http.NewRequest("GET", "https://golang.org/pkg/net/http/", nil)
	assert.NoError(t, err)

	err = tr.Record(req, &http.Response{
		Header: http.Header{
			"Content-Type": {"text/plain"},
		},
		Body: ioutil.NopCloser(strings.NewReader(helloWorld)),
	})
	assert.NoError(t, err)

	for i := 0; i < 5; i++ {
		res, err := client.Get("https://golang.org/pkg/net/http/")
		assert.NoError(t, err)

		buf, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, helloWorld, string(buf))

		err = res.Body.Close()
		assert.NoError(t, err)
	}
}

func TestStores(t *testing.T) {
	stores := []store.Store{
		file.NewStore("testdata", ".wire"),
		bolt.NewStore("testdata/bolt.db"),
		&mem.Mem{},
	}

	for _, s := range stores {
		tr := NewTransport()
		tr.SetStore(s)

		client := http.Client{Transport: tr}

		for i := 0; i < 5; i++ {
			_, err := client.Get("https://golang.org/pkg/net/http/")
			assert.NoError(t, err)
		}
	}
}

func TestEncoderDecoders(t *testing.T) {
	encodersDecoders := []dump.EncoderDecoder{
		&gob.Gob{},
		&json.JSON{},
		&wire.Wire{},
	}

	for _, ed := range encodersDecoders {
		tr := NewTransport()
		tr.SetEncoderDecoder(ed)

		client := http.Client{Transport: tr}

		for i := 0; i < 5; i++ {
			_, err := client.Get("https://golang.org/pkg/net/http/")
			assert.NoError(t, err)
		}
	}
}
