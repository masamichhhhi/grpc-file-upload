package main

import (
	"context"
	"fmt"
	"io"
	"os"

	upload "github.com/masamichhhhi/grpc-upload/proto"
	"google.golang.org/grpc"
)

func main() {
	connect, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer connect.Close()
	uploadService := upload.NewUploadServiceClient(connect)
	stream, err := uploadService.Upload(context.Background())
	if err != nil {
		panic(err)
	}
	err = Upload(stream)
	if err != nil {
		fmt.Println(err)
	}
}

func Upload(stream upload.UploadService_UploadClient) error {
	file, err := os.Open("./sample.mp4")
	if err != nil {
		return err
	}
	defer file.Close()
	buf := make([]byte, 1024)
	for {
		_, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = stream.Send(&upload.UploadRequest{MediaData: buf})
		if err != nil {
			fmt.Println(err)
		}
	}
	if err != nil {
		return err
	}

	return nil
}
