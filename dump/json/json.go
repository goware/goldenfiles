package json

import (
	"encoding/json"

	"github.com/goware/goldenfiles/dump"
)

type JSON struct {
}

func (j *JSON) Encode(res *dump.Response) ([]byte, error) {
	return json.MarshalIndent(res, "", "\t")
}

func (j *JSON) Decode(buf []byte) (*dump.Response, error) {
	var res dump.Response
	if err := json.Unmarshal(buf, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

var _ = dump.EncoderDecoder(&JSON{})
