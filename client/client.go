package main

import (
	"context"
	"log"

	"github.com/empaid/estateedge/services/common/genproto/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	conn, err := grpc.NewClient("localhost:3000", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Print("Error")
	}
	defer conn.Close()
	client := auth.NewAuthServiceClient(conn)
	// res, errr := client.Register(context.Background(), &auth.RegisterRequest{
	// 	Username: "dishapurohit25",
	// 	Password: "testtesttest",
	// 	Email:    "disha25@gmail.com",
	// })
	// if errr != nil {
	// 	log.Fatal("error while calling grpc function", errr)
	// }
	// log.Print(res)

	loginRes, errr := client.Login(context.Background(), &auth.LoginRequest{
		Username: "dishapurohit25",
		Password: "testtesttest",
	})
	if errr != nil {
		log.Fatal("error while calling grpc function", errr)
	}
	log.Print(loginRes)

	validateResponse, errr := client.Validate(context.Background(), &auth.ValidateRequest{
		AuthToken: loginRes.AuthToken,
	})
	if errr != nil {
		log.Fatal("error while calling grpc function", errr)
	}
	log.Print(validateResponse)

}
