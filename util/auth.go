package util

import (
	"crypto/rand"
	"fmt"
)

var (
	//ErrGenerateSaltFailed is returned when generating a salt fails
	ErrGenerateSaltFailed = fmt.Errorf("failed to generate salt")
)

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, ErrGenerateSaltFailed
	}
	return salt, nil
}
