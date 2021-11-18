// Package rand is a convenience to access the math/rand or crypto/rand packages,
// wrapping some of their functionality in easy to call, commonly used functions.
package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// RememberTokenBytes is the default number of bytes to use when generating a RememberToken.
const RememberTokenBytes = 32

// RandomBytes returns a slice of random bytes of length n.
// It uses crypto/rand so is pretty random.
func RandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// RandomString returns a base64 URL-encoded, version of a random slice of bytes of length n.
func RandomString(n int) (string, error) {
	b, err := RandomBytes(n)
	if err != nil {
		return "", fmt.Errorf("Could not generate random bytes: %w", err)
	}
	s := base64.URLEncoding.EncodeToString(b)
	return s, nil
}

// GenerateToken will generate a new base64 url-encoded string to be used for remember tokens.
func GenerateToken() (string, error) {
	return RandomString(RememberTokenBytes)
}

// NumBytes returns the length of a strings underlying byte-slice.
func NumBytes(s string) int {
	b := []byte(s)
	return len(b)
}
