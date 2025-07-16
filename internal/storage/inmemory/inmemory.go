package inmemory

import (
	"context"
	"sync"

	"github.com/socialrating/shortener/internal/storage"
)

type InMemoryStorage struct {
	mu          sync.RWMutex
	shortToOrig map[string]string
	origToShort map[string]string
}

func New() *InMemoryStorage {
	return &InMemoryStorage{
		shortToOrig: make(map[string]string),
		origToShort: make(map[string]string),
	}
}

func (s *InMemoryStorage) SaveURL(ctx context.Context, originalURL string, shortURL string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.shortToOrig[shortURL] = originalURL
	s.origToShort[originalURL] = shortURL
	return nil
}

func (s *InMemoryStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	originalURL, ok := s.shortToOrig[shortURL]
	if !ok {
		return "", storage.ErrURLNotFound
	}
	return originalURL, nil
}

func (s *InMemoryStorage) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	shortURL, ok := s.origToShort[originalURL]
	if !ok {
		return "", storage.ErrURLNotFound
	}
	return shortURL, nil
}
