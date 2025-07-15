package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/socialrating/shortener/internal/storage"
)

// URLService handles URL shortening operations
type URLService struct {
	storage storage.Storage
}

// NewURLService creates a new URL service instance
func NewURLService(storage storage.Storage) *URLService {
	return &URLService{storage: storage}
}

const (
	ShortURLLength = 10
	AllowedChars   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

// GenerateShortKey generates a short key for a given URL with an offset
func GenerateShortKey(originalURL string, offset int) string {
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


// GetOriginalURL retrieves the original URL for a given short key
func (s *URLService) GetOriginalURL(ctx context.Context, shortKey string) (string, error) {
	originalURL, err := s.storage.Get(ctx, shortKey)
	if err != nil {
		return "", fmt.Errorf("не удалось получить оригинальный URL: %w", err)
	}
	return originalURL, nil
}
