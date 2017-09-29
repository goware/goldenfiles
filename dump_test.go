package goldenfiles

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/goware/goldenfiles/dump"
	gobd "github.com/goware/goldenfiles/dump/gob"
	jsond "github.com/goware/goldenfiles/dump/json"
	wired "github.com/goware/goldenfiles/dump/wire"
	"github.com/stretchr/testify/assert"
)

func TestSerializeUnserialize(t *testing.T) {
	dumpers := []*dump.Dump{
		dump.New(&jsond.JSON{}),
		dump.New(&gobd.Gob{}),
		dump.New(&wired.Wire{}),
	}

	res, err := http.Get("https://golang.org/pkg/net/http/")
	assert.NoError(t, err)

	for _, dumper := range dumpers {

		buf, err := dumper.Encode(res)
		assert.NoError(t, err)

		res2, err := dumper.Decode(buf)
		assert.NoError(t, err)

		assert.Equal(t, res.Status, res2.Status)
		assert.Equal(t, res.StatusCode, res2.StatusCode)
		assert.Equal(t, res.ProtoMajor, res2.ProtoMajor)
		assert.Equal(t, res.ProtoMinor, res2.ProtoMinor)
		if res.ContentLength > 0 {
			assert.Equal(t, res.ContentLength, res2.ContentLength)
		}
		assert.Equal(t, res.TransferEncoding, res2.TransferEncoding)
		assert.Equal(t, res.Uncompressed, res2.Uncompressed)

		body, err := ioutil.ReadAll(res.Body)
		assert.NoError(t, err)

		body2, err := ioutil.ReadAll(res2.Body)
		assert.NoError(t, err)

		assert.Equal(t, body, body2)
	}
}
