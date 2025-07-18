package handler

import (
	"context"
	"errors"
	"net/http"

	api "github.com/socialrating/shortener/api/openapi/urlshortner"
	"github.com/socialrating/shortener/internal/service"
)

type RestHandler struct {
	service *service.Service
}

var _ api.StrictServerInterface = (*RestHandler)(nil)

func NewRestHandler(s *service.Service) *RestHandler {
	return &RestHandler{service: s}
}

func (h *RestHandler) Register(server *http.Server) {
	// Создаем обработчик сгенерированного OAPI
	handler := api.NewStrictHandler(h,nil)

	// Создаем стандартный роутер
	mux := http.NewServeMux()

	// Регистрируем маршруты на mux
	api.HandlerFromMux(handler, mux)

	// Устанавливаем mux как основной обработчик сервера
	server.Handler = mux
}

func (h *RestHandler) PostUrl(ctx context.Context, request api.PostUrlRequestObject) (api.PostUrlResponseObject, error) {
	if request.Body.Url == "" {
		errMsg := "URL is required"
		return api.PostUrl400JSONResponse{
			Error: &errMsg,
		}, nil
	}

	shortURL, err := h.service.CreateShortURL(ctx, request.Body.Url)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidURL):
			errMsg := "Invalid URL format"
			return api.PostUrl400JSONResponse{
				Error: &errMsg,
			}, nil
		default:
			errMsg := "Internal server error"
			return api.PostUrl500JSONResponse{
				Error: &errMsg,
			}, nil
		}
	}

	return api.PostUrl200JSONResponse{
		ShortUrl: shortURL,
	}, nil
}

func (h *RestHandler) GetShortUrl(ctx context.Context, request api.GetShortUrlRequestObject) (api.GetShortUrlResponseObject, error) {
	if request.ShortUrl == "" {
		errMsg := "short_url parameter is required"
		return api.GetShortUrl404JSONResponse{Error: &errMsg}, nil
	}

	originalURL, err := h.service.GetOriginalURL(ctx, request.ShortUrl)
	if err != nil {
		errMsg := "Internal server error"
		return api.GetShortUrl500JSONResponse{Error: &errMsg}, nil
	}
	return api.GetShortUrl307Response{
		Headers: api.GetShortUrl307ResponseHeaders{
			Location: originalURL,
		},
	}, nil
}
