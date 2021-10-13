package handler

import (
	"context"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	upload "github.com/masamichhhhi/grpc-upload/proto"
	"github.com/pkg/errors"

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
	req, err := stream.Recv()

	f, err := CreateTempFile(req)
	if err != nil {
		panic(err)
	}

	defer f.Close()
	defer os.Remove(f.Name())

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			return err
		}
		chunk := resp.GetChunk()
		if _, err := f.Write(chunk.GetMediaData()); err != nil {
			panic(err)
		}
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(credentialFilePath))
	if err != nil {
		log.Fatal(err)
	}

	bucketName := "grpc-test-masamichi"
	detail, err := f.Stat()
	if err != nil {
		panic(err)
	}
	objectPath := detail.Name()
	obj := client.Bucket(bucketName).Object(objectPath)
	writer := obj.NewWriter(ctx)

	if _, err = io.Copy(writer, f); err != nil {
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

func CreateTempFile(req *upload.FileRequest) (*os.File, error) {
	header := req.GetHeader()
	log.Println(header.GetName())
	// f, err := ioutil.TempFile(os.TempDir(), header.GetName())
	f, err := os.Create(header.GetName())
	if err != nil {
		return nil, errors.Wrap(err, "failed to create file")
	}
	for _, dict := range header.GetHeader() {
		log.Println("Key:", dict.GetKey())
		for _, val := range dict.GetValues() {
			log.Println("value: ", val)
		}
	}
	return f, nil
}
