package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

func SendAnalysisCompleteNotification(ctx context.Context, payload *notifyPayload) {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		log.Fatal("KAFKA_BROKERS env var is required")
	}
	brokers := strings.Split(brokersEnv, ",")

	topic := os.Getenv("KAFKA_TOPIC_FILE_ANALYZE_COMPLETE")
	if topic == "" {
		log.Fatal("KAFKA_TOPIC_FILE_ANALYZE_COMPLETE env var is required")
	}

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Dialer: &kafka.Dialer{
			Timeout:   10 * time.Second, // was 3s
			KeepAlive: 30 * time.Second,
		},

		WriteTimeout: 15 * time.Second, // was 10s
	})
	payload.Time = time.Now().Local().Format(time.RFC3339)
	body, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("ANALYSING: error while marshalling client")
		return
	}
	err = kafkaWriter.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(payload.Key),
			Value: body,
			Time:  time.Now(),
		},
	)
	if err != nil {
		log.Printf("ANALYSING: Failed to write Kafka message for %s/%s: %v", payload.Bucket, payload.Key, err)
		return
	}

	log.Printf("ANALYSING: Published to Kafka: %s/%s", payload.Bucket, payload.Key)
}
