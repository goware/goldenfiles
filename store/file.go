package store

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
)

var (
	reUnsafeFileChars = regexp.MustCompile(`[^a-zA-Z0-9\-_]`)
	reConsecutiveDash = regexp.MustCompile(`-+`)
)

var TestdataDir = "testdata"

type File struct {
	mu  sync.Mutex
	Ext string
}

const fileNameMaxLength = 128

func (f *File) safeFileName(s string) string {
	s = reUnsafeFileChars.ReplaceAllString(s, "-")
	s = reConsecutiveDash.ReplaceAllString(strings.Trim(s, "-"), "-")
	if len(s) > fileNameMaxLength {
		hash := fmt.Sprintf("%x", md5.Sum([]byte(s)))
		s = s[0:fileNameMaxLength-len(hash)-1] + "-" + hash
	}
	return TestdataDir + "/" + s + f.Ext
}

func (f *File) Get(key string) ([]byte, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	fileName := f.safeFileName(key)

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

	fileName := f.safeFileName(key)
	if err := os.MkdirAll(path.Dir(fileName), 0666); err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, val, 0666)
}

func (f *File) Delete(key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	return os.Remove(key)
}

var _ = Store(&File{})
