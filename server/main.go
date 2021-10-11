package main

import (
	"net"

	handler "github.com/masamichhhhi/grpc-upload/server/grpc"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()

	handler.NewUploadServer(server)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
