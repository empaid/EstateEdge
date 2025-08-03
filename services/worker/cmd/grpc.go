package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/empaid/estateedge/services/common/genproto/fileIngestion"
	"github.com/empaid/estateedge/services/common/genproto/fileUpload"
	"github.com/empaid/estateedge/services/worker/internal/clients"
	"github.com/empaid/estateedge/services/worker/internal/repository"
	"google.golang.org/grpc"
)

type FileUploadService struct {
	storage                    *repository.Storage
	fileIngestionServiceClient fileIngestion.FileIngestionServiceClient
	fileUpload.UnimplementedFileUploadServiceServer
}

func (s *FileUploadService) ReturnPreSignedUploadURL(ctx context.Context, req *fileUpload.UploadRequest) (*fileUpload.UploadResponse, error) {
	if userExists := s.storage.UserStore.CheckIfUserExists(ctx, req.UserId); !userExists {
		log.Fatal("User doesn't exists ")
		return nil, errors.New("User doesn't exist")
	}
	file := repository.File{
		UserId: int(req.UserId),
		Status: "UPLOAD_URL_GENERATE",
	}

	if err := s.storage.FileStore.CreateFile(ctx, &file); err != nil {
		log.Fatal("Error while creating new file ", err)
		return nil, err
	}

	res, err := s.fileIngestionServiceClient.ReturnPreSignedUploadURL(ctx, &fileIngestion.UploadRequest{
		File: &fileIngestion.File{
			Id:       file.ID,
			Name:     file.ID,
			Location: "aws",
		},
	})
	if err != nil {
		log.Fatal("error while grpc request to client, ", err)
		return nil, err
	}

	return &fileUpload.UploadResponse{
		UploadURL: res.URL,
	}, nil

}

func NewGrpcServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("Error while loading the service ", err)
	}
	grpc := grpc.NewServer()
	db, err := repository.NewDBConection()
	if err != nil {
		log.Fatal("Unable to connet to DB ", err)
		return
	}

	storage := repository.NewStorage(db)
	grpcConn, fileIngestionServiceClient := clients.NewFileIngestionServiceClient()
	defer grpcConn.Close()
	fileUploadHandler := &FileUploadService{
		storage:                    storage,
		fileIngestionServiceClient: fileIngestionServiceClient,
	}
	fileUpload.RegisterFileUploadServiceServer(grpc, fileUploadHandler)

	grpc.Serve(lis)

}
