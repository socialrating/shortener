package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/socialrating/shortener/internal/storage"
)

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(url string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к PostgreSQL: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("не удалось проверить подключение к PostgreSQL: %w", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
			id SERIAL PRIMARY KEY,
			short_key VARCHAR(10) NOT NULL UNIQUE,
			original_url TEXT NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_original_url ON urls(original_url);
	`)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать таблицу: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) Save(ctx context.Context, originalURL string) (string, error) {
	const maxAttempts = 5
	attempt := 0
	var key string
	var err error

	for attempt < maxAttempts {
		key = generateKey(originalURL, attempt)

		_, err = s.db.ExecContext(ctx, `
			INSERT INTO urls (short_key, original_url)
			VALUES ($1, $2)
			ON CONFLICT (original_url) DO NOTHING
		`, key, originalURL)

		if err == nil {
			// Проверяем, была ли вставка
			var inserted bool
			err = s.db.QueryRowContext(ctx, `
				SELECT EXISTS(SELECT 1 FROM urls WHERE original_url = $1)
			`, originalURL).Scan(&inserted)

			if err != nil {
				return "", fmt.Errorf("ошибка проверки вставки: %w", err)
			}

			if inserted {
				return key, nil
			}
			return "", storage.ErrAlreadyExists
		}

		// Если ошибка не связана с уникальностью ключа, прерываем
		if !isKeyConflict(err) {
			return "", fmt.Errorf("ошибка при сохранении URL: %w", err)
		}

		attempt++
	}

	return "", fmt.Errorf("не удалось сгенерировать уникальный ключ после %d попыток", maxAttempts)
}

func isKeyConflict(err error) bool {
	if err == nil {
		return false
	}
	// Проверяем на нарушение уникального ограничения
	return err.Error() == "pq: duplicate key value violates unique constraint \"urls_short_key_key\""
}

func (s *PostgresStorage) Get(ctx context.Context, shortKey string) (string, error) {
	if len(shortKey) != 10 {
		return "", storage.ErrNotFound
	}

	var originalURL string
	err := s.db.QueryRowContext(ctx, `
		SELECT original_url FROM urls WHERE short_key = $1
	`, shortKey).Scan(&originalURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrNotFound
		}
		return "", fmt.Errorf("ошибка при получении URL: %w", err)
	}

	return originalURL, nil
}

func (s *PostgresStorage) Close() error {
	return s.db.Close()
}

func generateKey(url string, attempt int) string {
	hash := fnvHash(url + string(rune(attempt)))
	key := make([]byte, 10)
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
