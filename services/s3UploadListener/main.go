package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/segmentio/kafka-go"
)

var (
	kafkaWriter *kafka.Writer
)

func init() {

	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		log.Fatal("KAFKA_BROKERS env var is required")
	}
	brokers := strings.Split(brokersEnv, ",")
	fmt.Print(brokers)
	topic := os.Getenv("KAFKA_TOPIC_FILE_UPLOAD")
	fmt.Print(topic)
	if topic == "" {
		log.Fatal("KAFKA_TOPIC_FILE_UPLOAD env var is required")
	}

	kafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second, // was 3s
			KeepAlive: 30 * time.Second,
		},

		WriteTimeout: 15 * time.Second, // was 10s
	})
}

type notifyPayload struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Time   string `json:"eventTime"`
}

func handler(ctx context.Context, evt events.S3Event) error {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		log.Fatal("KAFKA_BROKERS env var is required")
	}
	brokers := strings.Split(brokersEnv, ",")
	fmt.Print("BROKERS: ", brokers)
	topic := os.Getenv("KAFKA_TOPIC_FILE_UPLOAD")
	fmt.Print(topic)
	apiURL := os.Getenv("API_ENDPOINT")
	fmt.Print(apiURL)
	if apiURL == "" {
		return fmt.Errorf("API_ENDOINT don't exists")
	}
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for _, record := range evt.Records {
		bucket := record.S3.Bucket.Name
		key := record.S3.Object.Key

		payload := notifyPayload{
			Bucket: bucket,
			Key:    key,
			Time:   record.EventTime.Format(time.RFC3339),
		}

		body, err := json.Marshal(payload)

		if err != nil {
			return fmt.Errorf("failed to marshal notify object: %w", err)
		}
		writeCtx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		err = kafkaWriter.WriteMessages(writeCtx,
			kafka.Message{
				Key:   []byte(payload.Key),
				Value: body,
				Time:  time.Now(),
			},
		)
		if err != nil {
			log.Printf("Failed to write Kafka message for %s/%s: %v", payload.Bucket, payload.Key, err)
			return err
		}

		log.Printf("Published to Kafka: %s/%s", payload.Bucket, payload.Key)

		req, err := http.NewRequestWithContext(ctx, "POST", apiURL, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("failed to create HTTP request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("failed to call API: %w", err)
		}
		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return fmt.Errorf("API returned non-2xx status: %d", resp.StatusCode)
		}
		fmt.Printf("✅ Notified API for %s/%s\n", bucket, key)

	}
	return nil
}

func main() {
	defer kafkaWriter.Close()
	lambda.Start(handler)
}
