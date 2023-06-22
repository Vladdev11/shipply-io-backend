package models

import (
	"bytes"
	"crypto/sha256"
)

func HashPassword(password string, salt []byte) []byte {
	h := sha256.New()
	h.Write([]byte(password))
	h.Write(salt)
	return h.Sum(nil)
}

func ComparePassword(password string, salt []byte, hashedPassword []byte) bool {
	return bytes.Equal(hashedPassword, HashPassword(password, salt))
}
