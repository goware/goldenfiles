package mockingbird

import (
	"bytes"
	//"encoding/gob"
	//"encoding/json"
	//"errors"

	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	reHTTPStatusLine = regexp.MustCompile(`^HTTP/(\d+)\.(\d+)\s+(\d+)\s+(.+)$`)
	reHTTPHeader     = regexp.MustCompile(`^([a-zA-Z0-9\-]+):\s+(.+)$`)
)

type responseDump struct {
	Header           http.Header
	Body             string
	Status           string
	StatusCode       int
	ProtoMajor       int
	ProtoMinor       int
	ContentLength    int64
	TransferEncoding []string
	Uncompressed     bool
	Trailer          http.Header
}

func encode(dump *responseDump) ([]byte, error) {
	w := bytes.NewBuffer(nil)

	fmt.Fprintf(w,
		"HTTP/%d.%d %s\r\n",
		dump.ProtoMajor,
		dump.ProtoMinor,
		dump.Status,
	)
	for k, vv := range dump.Header {
		for _, v := range vv {
			fmt.Fprintf(w, "%s: %s\r\n", k, v)
		}
	}
	fmt.Fprintf(w, "\r\n")
	fmt.Fprintf(w, "%s", dump.Body)

	return w.Bytes(), nil
	/*
		var out bytes.Buffer
		enc := gob.NewEncoder(&out)
		if err := enc.Encode(dump); err != nil {
			return nil, err
		}

		return out.Bytes(), nil
	*/

	/*
		return json.MarshalIndent(dump, "", "\t")
	*/
}

func decode(buf []byte) (*responseDump, error) {
	dump := responseDump{
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
			case dump.StatusCode < 1:
				m := reHTTPStatusLine.FindAllStringSubmatch(s, -1)
				if len(m) > 0 {
					dump.ProtoMajor, _ = strconv.Atoi(m[0][1])
					dump.ProtoMinor, _ = strconv.Atoi(m[0][2])
					dump.StatusCode, _ = strconv.Atoi(m[0][3])
					dump.Status = fmt.Sprintf("%d %s", dump.StatusCode, m[0][4])
				}
			case readingHeaders:
				if s == "" {
					readingHeaders = false
					break
				}
				m := reHTTPHeader.FindAllStringSubmatch(s, -1)
				if len(m) > 0 {
					dump.Header.Add(m[0][0], m[0][1])
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

	if cl := dump.Header.Get("Content-Length"); cl != "" {
		dump.ContentLength, _ = strconv.ParseInt(cl, 10, 64)
	}

	dump.Body = string(body)

	return &dump, nil

	/*
		var dump responseDump
		if err := json.Unmarshal(buf, &dump); err != nil {
			return nil, err
		}
		return &dump, nil
	*/

	/*
		var dump responseDump

		in := bytes.NewBuffer(buf)
		dec := gob.NewDecoder(in)

		if err := dec.Decode(&dump); err != nil {
			return nil, err
		}

		return &dump, nil
	*/
}

func serializeResponse(res *http.Response) ([]byte, error) {
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Silently replacing body
	res.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	dump := responseDump{
		Header:           res.Header,
		Body:             string(buf),
		Status:           res.Status,
		StatusCode:       res.StatusCode,
		ProtoMajor:       res.ProtoMajor,
		ProtoMinor:       res.ProtoMinor,
		ContentLength:    res.ContentLength,
		TransferEncoding: res.TransferEncoding,
		Uncompressed:     res.Uncompressed,
		Trailer:          res.Trailer,
	}

	return encode(&dump)
}

func unserializeResponse(buf []byte) (*http.Response, error) {
	dump, err := decode(buf)
	if err != nil {
		return nil, err
	}

	res := http.Response{
		Header:           dump.Header,
		Body:             ioutil.NopCloser(bytes.NewBufferString(dump.Body)),
		Status:           dump.Status,
		StatusCode:       dump.StatusCode,
		ProtoMajor:       dump.ProtoMajor,
		ProtoMinor:       dump.ProtoMinor,
		ContentLength:    dump.ContentLength,
		TransferEncoding: dump.TransferEncoding,
		Uncompressed:     dump.Uncompressed,
		Trailer:          dump.Trailer,
	}

	return &res, nil
}
