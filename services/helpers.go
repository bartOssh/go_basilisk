package services

import (
	"crypto/rand"
	"encoding/hex"
)

// RandToken generates a random hex value of n length
func RandToken(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
