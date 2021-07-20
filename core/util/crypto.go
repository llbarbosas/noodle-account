package util

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/google/uuid"
)

func NewUUIDStr() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func RandomString(n uint) (string, error) {
	b := make([]byte, n)

	_, err := rand.Read(b)

	if err != nil {
		return "", err
	}

	bBase64 := base64.RawURLEncoding.EncodeToString(b)

	return bBase64, nil
}

func Sha256Base64(value string) string {
	hash := sha256.New().Sum([]byte(value))
	base64 := base64.URLEncoding.EncodeToString(hash)

	return base64
}
