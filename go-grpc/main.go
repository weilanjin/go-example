package main

import (
	v1 "github.com/weilanjin/go-example/go-grpc/idl/ecommerce/v1"
	"log"
	"net"

	"github.com/weilanjin/go-example/go-grpc/internal/domain/ecommerce/service"
	"google.golang.org/grpc"
)

const addrs = ":50051"

func main() {
	lis, err := net.Listen("tcp", addrs)
	if err != nil {
		log.Panic(err)
	}
	s := grpc.NewServer()

	v1.RegisterProductInfoServer(s, &service.ProductService{})

	log.Printf("Starting gRPC listener on port " + addrs)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	s.GracefulStop()
}