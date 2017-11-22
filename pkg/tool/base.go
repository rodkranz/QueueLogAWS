package tool

import (
	"crypto/md5"
	"fmt"
)

// MD5 create md5 from bytes.
func MD5(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
