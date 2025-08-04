package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/empaid/estateedge/services/worker/internal/repository"
	"github.com/segmentio/kafka-go"
)

type notifyPayload struct {
	Bucket  string `json:"bucket"`
	Key     string `json:"key"`
	Time    string `json:"eventTime"`
	Summary string `json:"summary"`
}

func NewReader(brokers []string, topic []string, groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupTopics: topic,
		GroupID:     groupId,
	})
}

func MessageHandler(ctx context.Context, m kafka.Message, s *repository.Storage) error {
	fileId := string(m.Key)
	log.Print("New Message from Kafka: ", fileId)
	var body notifyPayload
	if err := json.Unmarshal(m.Value, &body); err != nil {
		log.Fatal("Error while unmarshalling notification body, ", err)
	}
	if body.Summary == "" {
		file := repository.File{
			ID:     fileId,
			Status: "UPLOADED",
		}
		s.FileStore.ChangeFileStatus(ctx, &file)

		print("Message Value", string(m.Value))
		SendAnalyseNotification(ctx, &body)
	} else {
		file := repository.File{
			ID:      fileId,
			Status:  "COMPLETED",
			Summary: body.Summary,
		}
		s.FileStore.ChangeFileStatus(ctx, &file)
		s.FileStore.ChangeFileSummary(ctx, &file)
	}
	return nil
}

func SendAnalyseNotification(ctx context.Context, payload *notifyPayload) {
	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		log.Fatal("KAFKA_BROKERS env var is required")
	}
	brokers := strings.Split(brokersEnv, ",")

	topic := os.Getenv("KAFKA_TOPIC_FILE_ANALYZE")
	if topic == "" {
		log.Fatal("KAFKA_TOPIC_FILE_ANALYZE env var is required")
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
		log.Fatal("error while marshalling client")
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
		log.Printf("Failed to write Kafka message for %s/%s: %v", payload.Bucket, payload.Key, err)
		return
	}

	log.Printf("Published to Kafka: %s/%s", payload.Bucket, payload.Key)
}

func Listen(ctx context.Context, reader *kafka.Reader) error {
	db, err := repository.NewDBConection()
	if err != nil {
		log.Fatal("Unable to connet to DB ", err)
		return err
	}

	storage := repository.NewStorage(db)
	for {
		msg, err := reader.FetchMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				return nil
			}
			log.Print("Failed to load Message, ", err)
			continue
		}

		if err := MessageHandler(ctx, msg, storage); err != nil {
			log.Print("Failed to handle Message, ", err)
		}
		if err := reader.CommitMessages(ctx, msg); err != nil {
			log.Print("Failed to commit message, ", err)
		}
	}
}
