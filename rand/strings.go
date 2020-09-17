package rand

import (
	"crypto/rand"
	"encoding/base64"
)

const RememberTokenBytes = 32

// Bytes will help us genrate n random bytes, or will return an error if there was
// one. This uses the crypt/rand package so it is safe to use with things like remember token
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// String will gnerate a byte slice of size nByes and then return a string
// that is the base64 URL encoded version of that byte slice
func String(nBytes int) (string, error) {
	b, err := Bytes(nBytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken is a helper function designed to generate remember tokens of a predetermined byte size
func RememberToken() (string, error) {
	return String(RememberTokenBytes)
}

func NBytes(base64String string) (int, error) {
	b, err := base64.URLEncoding.DecodeString(base64String)
	if err != nil {
		return -1, err
	}
	return len(b), nil
}
