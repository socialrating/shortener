package postgres

import (
	"context"
	"database/sql"
	"github.com/socialrating/shortener/internal/storage"
)

type PostgresStorage struct {
	db *sql.DB
}

func New(databaseURL string) (*PostgresStorage, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS urls (
            id SERIAL PRIMARY KEY,
            original_url TEXT NOT NULL UNIQUE,
            short_url VARCHAR(10) NOT NULL UNIQUE
        );
    `)
	if err != nil {
		return nil, err
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) SaveURL(ctx context.Context, originalURL string, shortURL string) error {
	_, err := s.db.ExecContext(ctx, "INSERT INTO urls(original_url, short_url) VALUES($1, $2)", originalURL, shortURL)
	return err
}

func (s *PostgresStorage) GetURL(ctx context.Context, shortURL string) (string, error) {
	var originalURL string
	err := s.db.QueryRowContext(ctx, "SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err == sql.ErrNoRows {
		return "", storage.ErrURLNotFound
	}
	return originalURL, err
}

func (s *PostgresStorage) GetShortURL(ctx context.Context, originalURL string) (string, error) {
	var shortURL string
	err := s.db.QueryRowContext(ctx, "SELECT short_url FROM urls WHERE original_url = $1", originalURL).Scan(&shortURL)
	if err == sql.ErrNoRows {
		return "", storage.ErrURLNotFound
	}
	return shortURL, err
}
