package app

import (
	"fmt"
	"net"

	pb "github.com/socialrating/shortener/api/proto/urlshortner"
	"github.com/socialrating/shortener/internal/entrypoints/grpcapi/handler"
	"github.com/socialrating/shortener/internal/service"
	"google.golang.org/grpc"
)

type GRPCApp struct {
	port   string
	server *grpc.Server
}

func NewGRPCApp(port string, s *service.Service) *GRPCApp {
	grpcServer := grpc.NewServer()
	h := handler.NewGRPCHandler(s)
	pb.RegisterURLShortenerServer(grpcServer, h)

	return &GRPCApp{
		port:   port,
		server: grpcServer,
	}
}

func (a *GRPCApp) Run() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", a.port))
	if err != nil {
		return fmt.Errorf("failed to start listening to port %s: %w", a.port, err)
	}

	if err := a.server.Serve(lis); err != nil {
		return fmt.Errorf("gRPC server error: %w", err)
	}
	return nil
}

func (a *GRPCApp) Stop() {
	a.server.GracefulStop()
}
