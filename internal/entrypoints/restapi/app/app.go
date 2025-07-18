package app

import (
	"fmt"
	"net/http"

	"github.com/socialrating/shortener/internal/entrypoints/restapi/handler"
	"github.com/socialrating/shortener/internal/service"
)

type RestApp struct {
	*http.Server
	port string
}

func NewRestApp(port string, s *service.Service) *RestApp {
	h := handler.NewRestHandler(s)

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
	}

	h.Register(server)

	return &RestApp{
		port:   port,
		Server: server,
	}
}
