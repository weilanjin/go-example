### go grpc 依赖

macos brew install protobuf

go install github.com/golang/protobuf/protoc-gen-go

go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

go get -u google.golang.org/grpc

lanjin@lovec go-grpc % protoc --go_out=./idl/ecommerce --go-grpc_out=./idl/ecommerce ./idl/ecommerce/v1/*.proto
