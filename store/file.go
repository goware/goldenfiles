package store

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	reUnsafeFileChars = regexp.MustCompile(`[^a-zA-Z0-9\-]`)
	reConsecutiveDash = regexp.MustCompile(`-+`)
)

var TestdataDir = "testdata"

type File struct {
	mu sync.Mutex
}

func safeFileName(s string) string {
	s = reUnsafeFileChars.ReplaceAllString(s, "-")
	s = reConsecutiveDash.ReplaceAllString(strings.Trim(s, "-"), "-")
	_ = os.MkdirAll(TestdataDir, 0666)
	return TestdataDir + "/" + s + ".bin"
}

func (f *File) Get(key string) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileName := safeFileName(key)

	fp, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer fp.Close()
	return ioutil.ReadAll(fp)
}

func (f *File) Set(key string, val []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileName := safeFileName(key)

	return ioutil.WriteFile(fileName, val, 0666)
}

func (f *File) Delete(key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return os.Remove(key)
}

var _ = Store(&File{})
