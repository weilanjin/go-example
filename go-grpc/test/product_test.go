package test

import (
	"context"
	"log"
	"testing"
	"time"

	v1 "github.com/weilanjin/go-example/go-grpc/idl/ecommerce/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const addrs = "localhost:50051"

func TestProduct(t *testing.T) {
	conn, err := grpc.Dial(addrs, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := v1.NewProductInfoClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// 添加商品
	pid, err := c.AddProduct(ctx, &v1.Product{
		Name:        "Apple iPhone 15 Plus",
		Description: "Apple iPhone 15 Plus 128GB",
		Price:       1099.0,
	})
	if err != nil {
		t.Fatal(err)
	}
	// 获取商品
	p, err := c.GetProduct(ctx, &v1.ProductID{Value: pid.Value})
	if err != nil {
		t.Fatal(err)
	}
	log.Println(p)
}
