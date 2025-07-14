package app

import (
	"log"
	"net"

	"github.com/socialrating/shortener/internal/storage"
	"google.golang.org/grpc"
)

func StartGRPCServer(stor storage.Storage, port string) {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
		// pb.RegisterURLShortenerServer(s, NewServer(stor))
	log.Printf("gRPC server listening on %s", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
