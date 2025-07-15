package app

import (
	"fmt"
	"net"
	pburlshortener "github.com/socialrating/shortener/api/proto/urlshortner"

	"github.com/socialrating/shortener/config"
	"github.com/socialrating/shortener/internal/entrypoints/grpcapi/handler"
	"github.com/socialrating/shortener/internal/storage"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	server  *grpc.Server
	config  *config.Config
	storage storage.Storage
}

func NewGRPCApp(cfg *config.Config, storage storage.Storage) *GRPCApp {
	grpcServer := grpc.NewServer()

	// Регистрация обработчиков
	urlHandler := handler.NewGRPCHandler(storage)
	pburlshortener.RegisterUrlShortenerServiceServer(grpcServer, urlHandler)
	return &GRPCApp{
		server:  grpcServer,
		config:  cfg,
		storage: storage,
	}
}

func (a *GRPCApp) Start() error {
	lis, err := net.Listen("tcp", ":"+a.config.Server.GRPCPort)
	if err != nil {
		return fmt.Errorf("не удалось прослушать порт: %w", err)
	}

	fmt.Printf("gRPC сервер запущен на порту %s\n", a.config.Server.GRPCPort)
	if err := a.server.Serve(lis); err != nil {
		return fmt.Errorf("ошибка gRPC сервера: %w", err)
	}

	return nil
}

func (a *GRPCApp) Stop() {
	a.server.GracefulStop()
	fmt.Println("gRPC сервер остановлен")

	if err := a.storage.Close(); err != nil {
		fmt.Printf("Ошибка при закрытии хранилища: %v\n", err)
	}
}
