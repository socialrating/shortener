package storage

import (
	"context"
	"errors"
)

var (
	ErrURLNotFound = errors.New("url not found")
)

type Storage interface {
	SaveURL(ctx context.Context, originalURL string, shortURL string) error
	GetURL(ctx context.Context, shortURL string) (string, error)
	GetShortURL(ctx context.Context, originalURL string) (string, error)
}
