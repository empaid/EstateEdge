package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/empaid/estateedge/services/worker/internal/repository"
	"github.com/segmentio/kafka-go"
)

type notifyPayload struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	Time   string `json:"eventTime"`
}

func NewReader(brokers []string, topic string, groupId string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupId,
	})
}

func MessageHandler(ctx context.Context, m kafka.Message, s *repository.Storage) error {
	fileId := string(m.Key)
	log.Print("New Message from Kafka: ", fileId)
	var body notifyPayload
	if err := json.Unmarshal(m.Value, &body); err != nil {
		log.Fatal("Error while unmarshalling notification body, ", err)
	}
	file := repository.File{
		ID:     fileId,
		Status: "UPLOADED",
	}
	s.FileStore.ChangeFileStatus(ctx, &file)

	print("Message Value", string(m.Value))
	return nil
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
