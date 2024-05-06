package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	v1 "lovec.wlj/go-grpc/idl/ecommerce/v1"
	"lovec.wlj/go-grpc/internal/domain/ecommerce/service"
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
