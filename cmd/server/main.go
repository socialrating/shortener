package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/socialrating/shortener/config"
	grpcapp "github.com/socialrating/shortener/internal/entrypoints/grpcapi/app"
	restapp "github.com/socialrating/shortener/internal/entrypoints/restapi/app"
	"github.com/socialrating/shortener/internal/service"
	"github.com/socialrating/shortener/internal/storage"
	"github.com/socialrating/shortener/internal/storage/inmemory"
	"github.com/socialrating/shortener/internal/storage/postgres"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	var store storage.Storage
	switch cfg.Storage.Type {
	case "in-memory":
		store = inmemory.New()
		log.Println("Using in-memory storage")
	case "postgres":
		store, err = postgres.New(cfg.Storage.Postgres.URL)
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL: %v", err)
		}
		log.Println("Using PostgreSQL storage")
	default:
		log.Fatalf("Unknown storage type: %s", cfg.Storage.Type)
	}

	urlService := service.New(store)
	httpServer := restapp.NewRestApp(cfg.Server.HTTPPort, urlService)
	gRPCServer := grpcapp.NewGRPCApp(cfg.Server.GRPCPort, urlService)

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		shutdownSignal := make(chan os.Signal, 1)
		signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM)

		select {
		case sig := <-shutdownSignal:
			log.Printf("Received shutdown signal: %v", sig)
			return errors.New("graceful shutdown")
		case <-ctx.Done():
			return ctx.Err()
		}
	})

	g.Go(func() error {
		log.Printf("Starting gRPC server on port :%s", cfg.Server.GRPCPort)
		return gRPCServer.Run()
	})

	g.Go(func() error {
		log.Printf("Starting HTTP server on port :%s", httpServer.Addr)
		return httpServer.ListenAndServe()
	})

	if err := g.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) && !errors.Is(err, context.Canceled) {
		log.Printf("Server stopped with error: %v", err)
	}

	log.Println("Starting graceful shutdown of servers...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gRPCServer.Stop()
	log.Println("gRPC server stopped.")

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error during HTTP server shutdown: %v", err)
	} else {
		log.Println("HTTP server stopped.")
	}
}
