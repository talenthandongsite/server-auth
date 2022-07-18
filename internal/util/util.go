package util

import (
	"crypto"
	"encoding/hex"
)

func HashSHA256(plaintext string) string {
	hash := crypto.SHA256.New()
	hash.Write([]byte(plaintext))
	digest := hash.Sum(nil)
	return hex.EncodeToString(digest)
}
