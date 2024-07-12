package worker

import (
	"context"
	"gomicrosched/pkg/env"
	"gomicrosched/pkg/kafka"
	"os/signal"
	"strings"
	"syscall"
)

func Run() {
	var kReaderConfig = kafka.KafkaConfig{
		Topic:    env.MustGetEnv("KAFKA_TOPIC_WORKER"),
		Brokers:  strings.Split(env.MustGetEnv("KAFKA_BROKER"), ","),
		Username: env.MustGetEnv("KAFKA_USERNAME"),
		Password: env.MustGetEnv("KAFKA_PASSWORD"),
		GroupID:  "worker",
	}
	var kWriterConfig = kafka.KafkaConfig{
		Topic:    env.MustGetEnv("KAFKA_TOPIC_NOTIFIER"),
		Brokers:  strings.Split(env.MustGetEnv("KAFKA_BROKER"), ","),
		Username: env.MustGetEnv("KAFKA_USERNAME"),
		Password: env.MustGetEnv("KAFKA_PASSWORD"),
	}

	kReader := kafka.NewReader(kReaderConfig)
	kWriter := kafka.NewWriter(kWriterConfig)
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	kHandler := kafka.NewHandler(ctx, kReader, kWriter, 10000)
	defer kHandler.Close()
	kHandler.RunFunc(processMessage)

}
