package main

import (
	"crypto"
	"fmt"

	"github.com/rauwekost/silo/Godeps/_workspace/src/github.com/rauwekost/go-checksum"
)

func main() {
	//md5 of string
	c, err := checksum.String("password", crypto.MD5)
	if err != nil {
		panic(err)
	}
	fmt.Printf("MD5 = %s\n", c)

	//sha1 of string
	c, err = checksum.String("password", crypto.SHA1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("SHA1 = %s\n", c)

	//sha256 of file
	c, err = checksum.File("testfile.txt", crypto.SHA256)
	if err != nil {
		panic(err)
	}
	fmt.Printf("SHA256 = %s\n", c)
}
