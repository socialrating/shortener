package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	urlshortner "github.com/socialrating/shortener/api/proto/urlshortner"
	"github.com/socialrating/shortener/internal/storage"
)

// RESTHandler handles HTTP requests for URL shortening and redirection.
type RESTHandler struct {
	storage storage.Storage
}

func NewRESTHandler(storage storage.Storage) *RESTHandler {
	return &RESTHandler{storage: storage}
}

func (h *RESTHandler) ShortenURL(w http.ResponseWriter, r *http.Request) {
	var req urlshortner.ShortenUrlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос", http.StatusBadRequest)
		return
	}

	if req.Url == "" {
		http.Error(w, "URL обязателен", http.StatusBadRequest)
		return
	}

	if _, err := url.ParseRequestURI(req.Url); err != nil {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}

	shortKey, err := h.storage.Save(r.Context(), req.Url)
	if err != nil {
		if errors.Is(err, storage.ErrAlreadyExists) {
			// Возвращаем существующую ссылку
			h.sendResponse(w, shortKey, req.Url)
			return
		}
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	h.sendResponse(w, shortKey, req.Url)
}

func (h *RESTHandler) sendResponse(w http.ResponseWriter, shortKey, originalURL string) {
	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortKey) // Предполагаем, что сервер работает на порту 8080
	// Формируем ответ					

	response := urlshortner.ShortenUrlResponse{
		ShortUrl:    shortURL,
		OriginalUrl: originalURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *RESTHandler) RedirectToOriginal(w http.ResponseWriter, r *http.Request, shortKey string) {
	if len(shortKey) != 10 {
		http.Error(w, "Неверный формат короткой ссылки", http.StatusBadRequest)
		return
	}

	originalURL, err := h.storage.Get(r.Context(), shortKey)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			http.Error(w, "Ссылка не найдена", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusFound)
}
