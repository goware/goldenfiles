package mockingbird

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/goware/mockingbird/store"
	"github.com/stretchr/testify/assert"
)

func TestSerializeUnserialize(t *testing.T) {
	res, err := http.Get("https://golang.org/pkg/net/http/")
	assert.NoError(t, err)

	buf, err := serializeResponse(res)
	assert.NoError(t, err)

	res2, err := unserializeResponse(buf)
	assert.NoError(t, err)

	assert.Equal(t, res.Status, res2.Status)
	assert.Equal(t, res.StatusCode, res2.StatusCode)
	assert.Equal(t, res.ProtoMajor, res2.ProtoMajor)
	assert.Equal(t, res.ProtoMinor, res2.ProtoMinor)
	assert.Equal(t, res.ContentLength, res2.ContentLength)
	assert.Equal(t, res.TransferEncoding, res2.TransferEncoding)
	assert.Equal(t, res.Uncompressed, res2.Uncompressed)

	body, err := ioutil.ReadAll(res.Body)
	assert.NoError(t, err)

	body2, err := ioutil.ReadAll(res2.Body)
	assert.NoError(t, err)

	assert.Equal(t, body, body2)
}

func TestClientCachedResponse(t *testing.T) {
	client := New(nil)

	for i := 0; i < 5; i++ {
		_, err := client.Get("https://golang.org/pkg/net/http/")
		assert.NoError(t, err)
	}
}

func TestClientGoldenFile(t *testing.T) {
	db, err := bolt.Open("testdata/bolt", 0600, nil)
	assert.NoError(t, err)

	client := New(&store.Bolt{db})

	for i := 0; i < 5; i++ {
		_, err := client.Get("https://golang.org/pkg/net/http/")
		assert.NoError(t, err)
	}
}

func TestClientMockCustomBody(t *testing.T) {
	client := New(nil)

	helloWorld := "Hello world!"

	req, err := http.NewRequest("GET", "https://golang.org/pkg/net/http/", nil)
	assert.NoError(t, err)

	err = Record(client, req, &http.Response{
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
		&store.File{},
	}

	for _, s := range stores {
		client := New(s)
		for i := 0; i < 5; i++ {
			_, err := client.Get("https://golang.org/pkg/net/http/")
			assert.NoError(t, err)
		}
	}
}
