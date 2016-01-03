package checksum

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"fmt"
	"hash"
	"io"
	"math"
	"os"
)

const fileChunk = 8192 //8K

func createHash(method crypto.Hash) (hash.Hash, error) {
	var h hash.Hash

	switch method {
	case crypto.MD5:
		h = md5.New()
	case crypto.SHA1:
		h = sha1.New()
	case crypto.SHA224:
		h = sha256.New224()
	case crypto.SHA256:
		h = sha256.New()
	case crypto.SHA384:
		h = sha512.New384()
	case crypto.SHA512:
		h = sha512.New()
	default:
		return h, errors.New("Unsupported hashing method")
	}
	return h, nil
}

//String returns a hash for the given string
func String(s string, method crypto.Hash) (string, error) {
	h, err := createHash(method)
	if err != nil {
		return "", err
	}

	io.WriteString(h, s)

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

//Bytes return a hash for the given bytes
func Bytes(b []byte, method crypto.Hash) (string, error) {
	h, err := createHash(method)
	if err != nil {
		return "", err
	}
	io.WriteString(h, string(b))
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

//File returns a hash for the given file
func File(path string, method crypto.Hash) (string, error) {
	file, err := os.Open(path)

	if err != nil {
		return "", err
	}

	defer file.Close()

	stat, _ := file.Stat()
	size := stat.Size()
	chunks := uint64(math.Ceil(float64(size) / float64(fileChunk)))

	h, err := createHash(method)
	if err != nil {
		return "", err
	}

	for i := uint64(0); i < chunks; i++ {
		csize := int(math.Min(fileChunk, float64(size-int64(i*fileChunk))))
		buf := make([]byte, csize)

		file.Read(buf)
		io.WriteString(h, string(buf))
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
