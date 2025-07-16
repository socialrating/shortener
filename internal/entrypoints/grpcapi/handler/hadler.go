package handler

import (
	"context"

	pb "github.com/socialrating/shortener/api/proto/urlshortner"
	"github.com/socialrating/shortener/internal/service"
)

type GRPCHandler struct {
	pb.UnimplementedURLShortenerServer
	service *service.Service
}

func NewGRPCHandler(s *service.Service) *GRPCHandler {
	return &GRPCHandler{service: s}
}

func (h *GRPCHandler) CreateShortURL(ctx context.Context, req *pb.CreateShortURLRequest) (*pb.CreateShortURLResponse, error) {
	shortURL, err := h.service.CreateShortURL(ctx, req.GetOriginalUrl())
	if err != nil {
		return nil, err
	}
	return &pb.CreateShortURLResponse{ShortUrl: shortURL}, nil
}

func (h *GRPCHandler) GetOriginalURL(ctx context.Context, req *pb.GetOriginalURLRequest) (*pb.GetOriginalURLResponse, error) {
	originalURL, err := h.service.GetOriginalURL(ctx, req.GetShortUrl())
	if err != nil {
		return nil, err
	}
	return &pb.GetOriginalURLResponse{OriginalUrl: originalURL}, nil
}
