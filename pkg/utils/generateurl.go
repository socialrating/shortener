package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

const (
	ShortURLLength = 10
	AllowedChars   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

func GenerateShortURL(originalURL string, offset int) string {
	hash := sha256.Sum256([]byte(originalURL + string(rune(offset))))
	encoded := base64.URLEncoding.EncodeToString(hash[:])
	cleaned := strings.Map(func(r rune) rune {
		if strings.ContainsRune(AllowedChars, r) {
			return r
		}
		return -1
	}, encoded)
	if len(cleaned) >= ShortURLLength {
		return cleaned[:ShortURLLength]
	}
	return cleaned
}
