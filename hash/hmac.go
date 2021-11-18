package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
)

// HMAC wraps a hash.Hash function so we can use it without knowledge of what kind of hash it contains.
// In this context it is implied it will be a HMAC.
type HMAC struct {
	hash hash.Hash
}

// NewHMAC creates a new object embedded with a HMAC hash.
// It wraps the stdlib hmac.New() but removes the users obligation to provide a
// func() that returns a hash.Hash as this handled by the constructor.
func NewHMAC(key string) HMAC {
	h := hmac.New(sha256.New, []byte(key))
	return HMAC{
		hash: h,
	}
}

// Hash is a helper function that will write the input to the hash and return the result as a string.
func (h HMAC) Hash(input string) string {
	h.hash.Reset()
	h.hash.Write([]byte(input))
	b := h.hash.Sum(nil)
	return base64.URLEncoding.EncodeToString(b)
}
