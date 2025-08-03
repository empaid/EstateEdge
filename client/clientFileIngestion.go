package main

import (
	"context"
	"log"

	"github.com/empaid/estateedge/services/common/genproto/fileIngestion"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunFileIngestionClient() {

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	conn, err := grpc.NewClient("localhost:4000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Print("Error")
	}
	defer conn.Close()
	client := fileIngestion.NewFileIngestionServiceClient(conn)
	// res, errr := client.Register(context.Background(), &auth.RegisterRequest{
	// 	Username: "dishapurohit25",
	// 	Password: "testtesttest",
	// 	Email:    "disha25@gmail.com",
	// })
	// if errr != nil {
	// 	log.Fatal("error while calling grpc function", errr)
	// }
	// log.Print(res)

	loginRes, errr := client.ReturnPreSignedUploadURL(context.Background(), &fileIngestion.UploadRequest{
		File: &fileIngestion.File{
			Name:      "Test",
			Extension: "png",
			Id:        "45",
			Location:  "aws",
			Bucket:    "files",
		},
	})
	if errr != nil {
		log.Fatal("error while calling grpc function", errr)
	}
	log.Print(loginRes)

}
