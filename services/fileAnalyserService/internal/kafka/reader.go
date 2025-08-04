package kafka

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/empaid/estateedge/pkg/env"
	"github.com/segmentio/kafka-go"
)

type notifyPayload struct {
	Bucket  string `json:"bucket"`
	Key     string `json:"key"`
	Time    string `json:"eventTime"`
	Summary string `json:"summary"`
}

func NewReader(brokers []string, topic string, groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupId,
	})
}

func MessageHandler(ctx context.Context, m kafka.Message, c *s3.Client) error {
	fileId := string(m.Key)
	log.Print("New Message from Kafka: ", fileId)
	var body notifyPayload
	if err := json.Unmarshal(m.Value, &body); err != nil {
		log.Fatal("Error while unmarshalling notification body, ", err)
	}

	getInput := &s3.GetObjectInput{
		Bucket: &body.Bucket,
		Key:    &body.Key,
	}
	out, err := c.GetObject(ctx, getInput)
	if err != nil {
		log.Fatal("Object not present in s3")
		return err
	}
	defer out.Body.Close()

	imgBuf := &bytes.Buffer{}
	if _, err := io.Copy(imgBuf, out.Body); err != nil {
		log.Fatal("Unable to convert image to b64")
		return err
	}
	b64 := base64.StdEncoding.EncodeToString(imgBuf.Bytes())

	reqBody := map[string]interface{}{
		"model": "gpt-4o", // or another vision-capable model
		"messages": []interface{}{
			map[string]interface{}{
				"role": "system",
				"content": []interface{}{
					map[string]string{"type": "text", "text": "You are an assistant that analyzes images."},
				},
			},
			map[string]interface{}{
				"role": "user",
				"content": []interface{}{
					map[string]string{"type": "text", "text": fmt.Sprintf("Please describe the image stored at s3://%s/%s", body.Bucket, body.Key)},
					map[string]interface{}{
						"type": "image_url",
						"image_url": map[string]string{
							"url": "data:image/png;base64," + b64,
							// "mime_type": http.DetectContentType(imgBuf.Bytes()),
						},
					},
				},
			},
		},
	}
	bodyBytes, err := json.Marshal(reqBody)

	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	// 5) Send the HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+env.GetString("OPEN_AI_API_KEY", ""))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("OpenAI request failed: %w", err)
	}
	defer resp.Body.Close()

	// log.Printf("Raw OpenAI response:\n%s\n", string(resbodyBytes))

	var respJSON struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&respJSON); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	if len(respJSON.Choices) == 0 {
		return fmt.Errorf("no choices in response")
	}

	log.Printf("AI analysis for %s/%s:\n%s",
		body.Bucket, body.Key,
		respJSON.Choices[0].Message.Content,
	)

	body.Summary = respJSON.Choices[0].Message.Content

	SendAnalysisCompleteNotification(ctx, &body)

	return nil
}

func Listen(ctx context.Context, reader *kafka.Reader) error {

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(env.GetString("AWS_DEFAULT_REGION", "")),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		),
	)
	if err != nil {
		log.Printf("Error: %v", err)
		return err
	}
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(env.GetString("AWS_BASE_ENDPOINT", ""))
		o.UsePathStyle = true
	})

	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return nil
			}
			log.Print("Failed to load Message, ", err)
			continue
		}

		if err := MessageHandler(ctx, msg, client); err != nil {
			log.Print("Failed to handle Message, ", err)
		}
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Print("Failed to commit message, ", err)
		}
	}
}
