package handler

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	pburlshortener "github.com/socialrating/shortener/api/proto/urlshortener"

	"github.com/socialrating/shortener/internal/storage"
)

type GRPCHandler struct {
	storage storage.Storage
	pburlshortener.UnimplementedUrlShortenerServiceServer
}

func NewGRPCHandler(storage storage.Storage) *GRPCHandler {
	return &GRPCHandler{storage: storage}
}

func (h *GRPCHandler) ShortenUrl(ctx context.Context, req *pburlshortener.ShortenUrlRequest) (*pburlshortener.ShortenUrlResponse, error) {
	if _, err := url.ParseRequestURI(req.Url); err != nil {
		return nil, fmt.Errorf("неверный URL: %w", err)
	}

	shortKey, err := h.storage.Save(ctx, req.Url)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			return &pburlshortener.ShortenUrlResponse{
				ShortUrl:    shortKey,
				OriginalUrl: req.Url,
			}, nil
		}
		return nil, fmt.Errorf("ошибка при создании короткой ссылки: %w", err)
	}

	return &pburlshortener.ShortenUrlResponse{
		ShortUrl:    shortKey,
		OriginalUrl: req.Url,
	}, nil
}

func (h *GRPCHandler) GetOriginalUrl(ctx context.Context, req *pburlshortener.GetOriginalUrlRequest) (*pburlshortener.GetOriginalUrlResponse, error) {
	originalURL, err := h.storage.Get(ctx, req.ShortKey)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить оригинальный URL: %w", err)
	}

	return &pburlshortener.GetOriginalUrlResponse{
		OriginalUrl: originalURL,
	}, nil
}
