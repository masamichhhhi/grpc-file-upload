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
	"gopkg.in/h2non/filetype.v1"
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
		err = Upload(stream, file, fileHeader)
		if err != nil {
			fmt.Println(err)
		}

		log.Printf("complete")

	})
	http.ListenAndServe(":8081", nil)
}

func Upload(stream upload.UploadService_UploadClient, file multipart.File, info *multipart.FileHeader) error {
	h := MakeHeader(info)
	h.Add("content-type", mime(file, info))
	stream.Send(&upload.FileRequest{
		File: h.Cast(),
	})

	buf := make([]byte, 1024)
	for {
		_, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		err = stream.Send(&upload.FileRequest{
			File: &upload.FileRequest_Chunk{
				Chunk: &upload.ChunkType{
					MediaData: buf,
				},
			}})
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

type FileHeader upload.FileRequest_Header

func (f *FileHeader) Add(key, value string) {
	h := f.Header.Header
	for i := 0; i < len(h); i++ {
		if h[i].Key == key {
			h[i].Values = append(h[i].Values, value)
			return
		}
	}

	h = append(h, &upload.FileHeader_MIMEHeaderType{
		Key:    key,
		Values: []string{value},
	})
}

func (f *FileHeader) Cast() *upload.FileRequest_Header {
	if f == nil {
		return nil
	}

	h := upload.FileRequest_Header(*f)
	return &h
}

func MakeHeader(info *multipart.FileHeader) *FileHeader {
	log.Println(info.Filename)
	h := FileHeader(
		upload.FileRequest_Header{
			Header: &upload.FileHeader{
				Name: info.Filename,
			},
		},
	)
	return &h
}

func mime(f multipart.File, info *multipart.FileHeader) string {
	head := make([]byte, 261)
	f.Read(head)
	f.Seek(0, 0)
	kind, err := filetype.Match(head)
	if err != nil {
		return "unknown"
	}
	return kind.MIME.Value
}
