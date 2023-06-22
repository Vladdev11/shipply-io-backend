package util

import (
	"crypto/rand"
	"fmt"
)

func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, fmt.Errorf("error generating salt: %v", err)
	}
	return salt, nil
}
