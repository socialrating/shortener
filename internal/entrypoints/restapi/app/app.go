package app

import (
	"fmt"
	"net/http"

	"github.com/socialrating/shortener/internal/entrypoints/restapi/handler"
	"github.com/socialrating/shortener/internal/service"
)

type RestApp struct {
	port string
	*http.Server
}

func NewRestApp(port string, s *service.Service) *RestApp {
	h := handler.NewRestHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /", h.CreateShortURL)
	mux.HandleFunc("GET /{short_url}", h.GetOriginalURL)

	return &RestApp{
		port: port,
		Server: &http.Server{
			Addr:    fmt.Sprintf(":%s", port),
			Handler: mux,
		},
	}
}
