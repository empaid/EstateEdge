package main

import (
	"context"
	"log"
	"net"

	"github.com/empaid/estateedge/services/common/genproto/fileIngestion"
	"github.com/empaid/estateedge/services/fileIngestion/internal/storage"
	"github.com/empaid/estateedge/services/fileIngestion/internal/types"
	"google.golang.org/grpc"
)

type FileIngestionService struct {
	storage types.Storage
	fileIngestion.UnimplementedFileIngestionServiceServer
}

func (f *FileIngestionService) ReturnPreSignedUploadURL(ctx context.Context, req *fileIngestion.UploadRequest) (*fileIngestion.UploadResponse, error) {

	url, err := f.storage.ReturnPreSignedUploadURL(ctx, req.File)
	if err != nil {
		log.Fatal("Error while returning the string", err)
		return nil, err
	}

	return &fileIngestion.UploadResponse{
		URL:  url,
		File: req.File,
	}, nil

}

func NewGrpcServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Unable to start auth server:", err)
	}

	ctx := context.Background()
	awsStorage, err := storage.NewAwsStorage(ctx)
	if err != nil {
		log.Fatal("Unable to initialize aws storage", err)
	}
	fileIngestionHandler := &FileIngestionService{
		storage: awsStorage,
	}

	grpc := grpc.NewServer()

	fileIngestion.RegisterFileIngestionServiceServer(grpc, fileIngestionHandler)
	log.Print("Server started listening")
	err = grpc.Serve(lis)
	if err != nil {
		log.Fatal("Unable to start grpc server:", err)
	}

}
