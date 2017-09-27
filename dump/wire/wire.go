package wire

import (
	"bytes"

	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/goware/mockingbird/dump"
)

var (
	reHTTPStatusLine = regexp.MustCompile(`^HTTP/(\d+)\.(\d+)\s+(\d+)\s+(.+)$`)
	reHTTPHeader     = regexp.MustCompile(`^([a-zA-Z0-9\-]+):\s+(.+)$`)
)

type Wire struct {
}

func (wire *Wire) Encode(res *dump.Response) ([]byte, error) {
	w := bytes.NewBuffer(nil)

	fmt.Fprintf(w,
		"HTTP/%d.%d %s\r\n",
		res.ProtoMajor,
		res.ProtoMinor,
		res.Status,
	)
	for k, vv := range res.Header {
		for _, v := range vv {
			fmt.Fprintf(w, "%s: %s\r\n", k, v)
		}
	}
	fmt.Fprintf(w, "\r\n")
	fmt.Fprintf(w, "%s", res.Body)

	return w.Bytes(), nil
}

func (wire *Wire) Decode(buf []byte) (*dump.Response, error) {
	res := dump.Response{
		Header:        http.Header{},
		ContentLength: -1,
		Uncompressed:  true,
	}

	var line, body []byte

	readingHeaders := true
	for i, c := range buf {
		line = append(line, c)
		if i < 1 {
			continue
		}

		if c == '\n' && buf[i-1] == '\r' {
			s := strings.TrimSpace(string(line))
			switch {
			case res.StatusCode < 1:
				m := reHTTPStatusLine.FindAllStringSubmatch(s, -1)
				if len(m) > 0 {
					res.ProtoMajor, _ = strconv.Atoi(m[0][1])
					res.ProtoMinor, _ = strconv.Atoi(m[0][2])
					res.StatusCode, _ = strconv.Atoi(m[0][3])
					res.Status = fmt.Sprintf("%d %s", res.StatusCode, m[0][4])
				}
			case readingHeaders:
				if s == "" {
					readingHeaders = false
					break
				}
				m := reHTTPHeader.FindAllStringSubmatch(s, -1)
				if len(m) > 0 {
					res.Header.Add(m[0][0], m[0][1])
				}
			default:
				body = append(body, line...)
			}
			line = []byte{}
		}
	}
	if len(line) > 0 {
		body = append(body, line...)
	}

	if cl := res.Header.Get("Content-Length"); cl != "" {
		res.ContentLength, _ = strconv.ParseInt(cl, 10, 64)
	}

	res.Body = string(body)

	return &res, nil
}

var _ = dump.EncoderDecoder(&Wire{})
