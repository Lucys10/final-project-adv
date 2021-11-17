package sha3hash

import (
	"fmt"
	"golang.org/x/crypto/sha3"
)

func GetHashSHA3(s string) string {
	h := sha3.Sum256([]byte(s))
	return fmt.Sprintf("%x", h)
}
