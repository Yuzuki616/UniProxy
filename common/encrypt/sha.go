package encrypt

import (
	"crypto/sha1"
	"encoding/hex"
)

func Sha(b []byte) string {
	e := sha1.Sum(b)
	return hex.EncodeToString(e[:])
}
