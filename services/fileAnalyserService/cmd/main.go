package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/empaid/estateedge/services/fileAnalyserService/internal/kafka"
)

func main() {

	brokersEnv := os.Getenv("KAFKA_BROKERS")
	if brokersEnv == "" {
		log.Fatal("KAFKA_BROKERS env var is required")
	}
	brokers := strings.Split(brokersEnv, ",")
	topic := os.Getenv("KAFKA_TOPIC_FILE_ANALYZE")
	if topic == "" {
		log.Fatal("KAFKA_TOPIC_FILE_ANALYZE env var is required")
	}
	groupId := os.Getenv("KAFKA_GROUP_ID")
	if groupId == "" {
		log.Fatal("KAFKA_GROUP_ID env var is required")
	}

	reader := kafka.NewReader(brokers, topic, groupId)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	log.Print("Starting Kafka Listener")
	if err := kafka.Listen(ctx, reader); err != nil {
		log.Fatal("Failed kakfa listener, ", err)
	}

}
