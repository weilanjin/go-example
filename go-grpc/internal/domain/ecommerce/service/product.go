package service

import (
	"context"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "lovec.wlj/go-grpc/idl/ecommerce/v1"
	"lovec.wlj/pkg/uid"
)

var storage = sync.Map{}

var _ v1.ProductInfoServer = (*ProductService)(nil)

type ProductService struct {
	v1.UnimplementedProductInfoServer
}

// AddProduct implements v1.ProductInfoServer.
func (e *ProductService) AddProduct(ctx context.Context, in *v1.Product) (*v1.ProductID, error) {
	in.Id = uid.UUID()
	storage.Store(in.Id, in)
	return &v1.ProductID{Value: in.Id}, status.New(codes.OK, "").Err()
}

// GetProduct implements v1.ProductInfoServer.
func (e *ProductService) GetProduct(ctx context.Context, in *v1.ProductID) (*v1.Product, error) {
	if in.Value == "" {
		return nil, status.New(codes.InvalidArgument, "").Err()
	}
	if v, ok := storage.Load(in.Value); ok {
		return v.(*v1.Product), status.New(codes.OK, "").Err()
	}
	return nil, status.New(codes.NotFound, "").Err()
}
