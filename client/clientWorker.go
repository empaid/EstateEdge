package main

import (
	"context"
	"log"

	"github.com/empaid/estateedge/pkg/env"
	"github.com/empaid/estateedge/services/common/genproto/fileUpload"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunWorkerService() {

	conn, err := grpc.NewClient(env.GetString("WORKER_SERVICE_ADDR", ""), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Print("Error")
	}
	defer conn.Close()
	client := fileUpload.NewFileUploadServiceClient(conn)

	res, errr := client.UploadFile(context.Background(), &fileUpload.UploadRequest{
		UserId: 12,
	})
	if errr != nil {
		log.Fatal("error while calling grpc function", errr)
	}
	log.Print(res)

}
