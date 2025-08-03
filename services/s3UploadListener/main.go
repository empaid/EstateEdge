package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type notifyPayload struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Time   string `json:"eventTime"`
}

func handler(ctx context.Context, evt events.S3Event) error {
	apiURL := os.Getenv("API_ENDPOINT")
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
			return fmt.Errorf("Failed to marshal notify object: ", err)
		}

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
	lambda.Start(handler)
}
