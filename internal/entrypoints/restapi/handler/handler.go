package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/socialrating/shortener/internal/service"
)

type RestHandler struct {
	service *service.Service
}

func NewRestHandler(s *service.Service) *RestHandler {
	return &RestHandler{service: s}
}

func (h *RestHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Cannot read body", http.StatusBadRequest)
		return
	}

	var request struct {
		URL string `json:"url"`
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	shortURL, err := h.service.CreateShortURL(r.Context(), request.URL)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := struct {
		ShortURL string `json:"short_url"`
	}{ShortURL: shortURL}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *RestHandler) GetOriginalURL(w http.ResponseWriter, r *http.Request) {
	shortURL := r.URL.Path[1:]
	originalURL, err := h.service.GetOriginalURL(r.Context(), shortURL)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
