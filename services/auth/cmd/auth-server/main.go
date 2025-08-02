package main

import (
	"log"
	"net"

	"github.com/empaid/estateedge/services/auth/internal/repository"
	"github.com/empaid/estateedge/services/common/genproto/auth"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

func main() {
	// log.Print("New server started")

	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Printf("failed to start the GRPC server: %s", ":3000")
	}

	grpcServer := grpc.NewServer()
	db := repository.NewConnection()

	authHandler := authService{
		store: repository.NewStorage(db.Conn),
	}
	auth.RegisterAuthServiceServer(grpcServer, authHandler)
	log.Print("Server started listening ")
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatal("Error while creating server")
	}

}
