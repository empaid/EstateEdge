package storage

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/empaid/estateedge/services/common/genproto/fileIngestion"
)

type AwsStorage struct {
	presigner *s3.PresignClient
	client    *s3.Client
}

func (s *AwsStorage) ReturnPreSignedUploadURL(ctx context.Context, file *fileIngestion.File) (string, error) {

	params := &s3.PutObjectInput{
		Bucket: aws.String(file.Bucket),
		Key:    aws.String(file.Id),
	}

	if _, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(file.Bucket),
	}); err != nil {
		log.Print("Bucket not present..Creating new bucket ", file.Bucket)
		if _, err := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(file.Bucket),
		}); err != nil {
			log.Print("Unable to create bucket", err)
			return "", err
		}

	}
	presignedReq, err := s.presigner.PresignPutObject(ctx, params)
	if err != nil {
		log.Print("Error while creating new presigned URL")
		return "", err
	}

	return presignedReq.URL, nil

}

func NewAwsStorage(ctx context.Context) (*AwsStorage, error) {

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil, err
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://localhost:4566")
		o.UsePathStyle = true
	})
	presigner := s3.NewPresignClient(client)
	return &AwsStorage{
		presigner: presigner,
		client:    client,
	}, nil
}
