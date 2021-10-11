package handler

import (
	"io"
	"os"
	"path/filepath"

	upload "github.com/masamichhhhi/grpc-upload/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func NewUploadServer(gserver *grpc.Server) {
	uploadServer := &server{}
	upload.RegisterUploadServiceServer(gserver, uploadServer)
	reflection.Register(gserver)
}

func (s *server) Upload(stream upload.UploadService_UploadServer) error {
	err := os.MkdirAll("Sample", 0777)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath.Join("Sample", "tmp.mp4"))
	defer file.Close()
	if err != nil {
		return err
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		file.Write(resp.MediaData)
	}

	err = stream.SendAndClose(&upload.UploadReply{UploadStatus: "ok"})
	if err != nil {
		return err
	}

	return nil
}
