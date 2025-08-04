package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/empaid/estateedge/pkg/env"
	"github.com/empaid/estateedge/services/worker/internal/kafka"
)

func main() {

	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		log.Fatal("KAFKA_BROKERS env var is required")
	}
	brokers := strings.Split(brokersEnv, ",")
	topic_file_upload := os.Getenv("KAFKA_TOPIC_FILE_UPLOAD")
	if topic_file_upload == "" {
		log.Fatal("KAFKA_TOPIC_FILE_UPLOAD env var is required")
	}

	topic_analyze_complete := os.Getenv("KAFKA_TOPIC_FILE_ANALYZE_COMPLETE")
	if topic_analyze_complete == "" {
		log.Fatal("KAFKA_TOPIC_FILE_ANALYZE_COMPLETE env var is required")
	}

	groupId := os.Getenv("KAFKA_GROUP_ID")
	if groupId == "" {
		log.Fatal("KAFKA_GROUP_ID env var is required")
	}
	topics := []string{topic_file_upload, topic_analyze_complete}
	reader := kafka.NewReader(brokers, topics, groupId)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// storage := repository.NewStorage(db)
	go func() {
		log.Print("Starting Kafka Listener")
		if err := kafka.Listen(ctx, reader); err != nil {
			log.Fatal("Failed kakfa listener, ", err)
		}
	}()

	NewGrpcServer(env.GetString("WORKER_SERVICE_ADDR", ""))

}
