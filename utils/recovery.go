package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateRecoveryKey() (string, error) {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}
