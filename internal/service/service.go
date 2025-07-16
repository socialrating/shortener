package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"strings"

	"github.com/socialrating/shortener/internal/storage"
)

type Service struct {
	storage storage.Storage
}

func New(s storage.Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	shortURL, err := s.storage.GetShortURL(ctx, originalURL)
	if err == nil {
		return shortURL, nil
	}

	shortURL = s.generateShortURL(originalURL)

	err = s.storage.SaveURL(ctx, originalURL, shortURL)
	if err != nil {
		return "", err
	}

	return shortURL, nil
}

func (s *Service) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	return s.storage.GetURL(ctx, shortURL)
}

func (s *Service) generateShortURL(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	hash := hasher.Sum(nil)

	encoded := base64.URLEncoding.EncodeToString(hash)
	encoded = strings.ReplaceAll(encoded, "-", "_")
	encoded = strings.ReplaceAll(encoded, "=", "")

	return encoded[:10]
}
