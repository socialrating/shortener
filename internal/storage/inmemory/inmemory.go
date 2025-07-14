package inmemory

import (
	"context"
	"sync"

	"github.com/socialrating/shortener/internal/storage"
)

type InMemoryStorage struct {
	mu        sync.RWMutex
	urls      map[string]string // shortKey -> originalURL
	urlToKey  map[string]string // originalURL -> shortKey
	keyGen    func(string) string
	keyLength int
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		urls:      make(map[string]string),
		urlToKey:  make(map[string]string),
		keyLength: 10,
	}
}

func (s *InMemoryStorage) generateKey(url string) string {
	// Простая реализация для примера
	hash := fnvHash(url)
	key := make([]byte, s.keyLength)
	for i := range key {
		idx := int(hash) % len(charset)
		key[i] = charset[idx]
		hash = hash >> 6
	}
	return string(key)
}

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_"

func fnvHash(s string) uint32 {
	hash := uint32(2166136261)
	for _, c := range s {
		hash *= 16777619
		hash ^= uint32(c)
	}
	return hash
}

func (s *InMemoryStorage) Save(ctx context.Context, originalURL string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Проверяем, существует ли уже URL
	if key, exists := s.urlToKey[originalURL]; exists {
		return key, storage.ErrAlreadyExists
	}

	// Генерируем уникальный ключ
	attempt := 0
	var key string
	for {
		key = s.generateKey(originalURL + string(rune(attempt)))
		if _, exists := s.urls[key]; !exists {
			break
		}
		attempt++
	}

	// Сохраняем
	s.urls[key] = originalURL
	s.urlToKey[originalURL] = key

	return key, nil
}

func (s *InMemoryStorage) Get(ctx context.Context, shortKey string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if len(shortKey) != s.keyLength {
		return "", storage.ErrNotFound
	}

	url, exists := s.urls[shortKey]
	if !exists {
		return "", storage.ErrNotFound
	}

	return url, nil
}

func (s *InMemoryStorage) Close() error {
	return nil
}
