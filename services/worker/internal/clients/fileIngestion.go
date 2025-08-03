package clients

import (
	"log"

	"github.com/empaid/estateedge/pkg/env"
	"github.com/empaid/estateedge/services/common/genproto/fileIngestion"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewFileIngestionServiceClient() (*grpc.ClientConn, fileIngestion.FileIngestionServiceClient) {
	conn, err := grpc.NewClient(env.GetString("FILE_INGESTION_SERVICE_ADDR", ""), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Print("Error")
	}
	client := fileIngestion.NewFileIngestionServiceClient(conn)
	return conn, client
}
