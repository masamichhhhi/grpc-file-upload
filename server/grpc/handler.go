package handler

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/storage"
	upload "github.com/masamichhhhi/grpc-upload/proto"
	"google.golang.org/api/option"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const credentialFilePath = "./key.json"

type server struct{}

func NewUploadServer(gserver *grpc.Server) {
	uploadServer := &server{}
	upload.RegisterUploadServiceServer(gserver, uploadServer)
	reflection.Register(gserver)
}

func (s *server) Upload(stream upload.UploadService_UploadServer) error {
	tempfile, err := ioutil.TempFile(os.TempDir(), "sample")
	defer tempfile.Close()
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}
		tempfile.Write(resp.MediaData)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
	if err != nil {
		log.Fatal(err)
	}

	bucketName := "grpc-test-masamichi"
	objectPath := "sample.mp4"
	obj := client.Bucket(bucketName).Object(objectPath)
	writer := obj.NewWriter(ctx)

	if _, err = io.Copy(writer, tempfile); err != nil {
		log.Println(err)
	}

	if err = writer.Close(); err != nil {
		log.Println(err)
	}

	log.Println("done")

	err = stream.SendAndClose(&upload.UploadReply{UploadStatus: "ok"})
	if err != nil {
		return err
	}

	return nil
}
