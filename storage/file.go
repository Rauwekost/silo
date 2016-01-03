package storage

import (
	"bytes"
	"crypto"

	"github.com/rauwekost/silo/Godeps/_workspace/src/github.com/rauwekost/go-checksum"
)

type File interface {
	Size() int64
	Read(b []byte) (int, error)
	ReadAll() ([]byte, error)
	Close() error
	Checksum() string
}

type ContentFile struct {
	*bytes.Buffer
}

func NewContentFile(content []byte) *ContentFile {
	return &ContentFile{bytes.NewBuffer(content)}
}

func (f *ContentFile) Close() error {
	return nil
}

func (f *ContentFile) Size() int64 {
	return int64(f.Len())
}

func (f *ContentFile) ReadAll() ([]byte, error) {
	return f.Bytes(), nil
}

func (f *ContentFile) Checksum() string {
	sha, err := checksum.Bytes(f.Bytes(), crypto.SHA256)
	if err != nil {
		panic(err)
	}
	return sha
}
