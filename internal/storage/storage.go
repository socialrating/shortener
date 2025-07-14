package storage

import (
	"context"
	"errors"
)

var (
	ErrNotFound      = errors.New("ссылка не найдена")
	ErrAlreadyExists = errors.New("ссылка уже существует")
)

type Storage interface {
	Save(ctx context.Context, originalURL string) (string, error)
	Get(ctx context.Context, shortKey string) (string, error)
	Close() error
}
