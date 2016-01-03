Checksum
--------

Create checksums with ease

##Installation
Download and install <code>go-checksum</code> with go get

```
go get github.com/rauwekost/go-checksum
```

##Usage
```go
package main

import (
  "crypto"
  "fmt"

  "github.com/rauwekost/go-checksum"
)

func main() {
  //md5 of a string
  c, err := checksum.String("password", crypto.MD5)
  if err != nil {
    panic(err)
  }
  fmt.Printf("MD5 = %s\n", c)

  //sha1 of a string
  c, err = checksum.String("password", crypto.SHA1)
  if err != nil {
    panic(err)
  }
  fmt.Printf("SHA1 = %s\n", c)

  //sha256 of a file
  c, err = checksum.File("testfile.txt", crypto.SHA256)
  if err != nil {
    panic(err)
  }
  fmt.Printf("SHA256 = %s\n", c)
}

```