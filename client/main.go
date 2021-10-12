package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	upload "github.com/masamichhhhi/grpc-upload/proto"
	"google.golang.org/grpc"
)

func main() {
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		file, fileHeader, err := r.FormFile("file")

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer file.Close()

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
		err = Upload(stream, file, fileHeader.Filename)
		if err != nil {
			fmt.Println(err)
		}

		log.Printf("complete")

	})
	http.ListenAndServe(":8081", nil)
}

func Upload(stream upload.UploadService_UploadClient, file multipart.File, fileName string) error {
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

	resp, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}
	fmt.Println(resp.UploadStatus)

	return nil

}
