package notifier

import (
	"context"
	"gomicrosched/pkg/env"
	"gomicrosched/pkg/kafka"
	"gomicrosched/pkg/logger"
	"os/signal"
	"strings"
	"syscall"
)

func init() {
	logger.ConfigureLogger()
}

func Run() {
	var kConfig = kafka.KafkaConfig{
		Topic:    env.MustGetEnv("KAFKA_TOPIC_NOTIFIER"),
		Brokers:  strings.Split(env.MustGetEnv("KAFKA_BROKER"), ","),
		Username: env.MustGetEnv("KAFKA_USERNAME"),
		Password: env.MustGetEnv("KAFKA_PASSWORD"),
	}

	kafkaReader := kafka.NewReader(kConfig)
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	kHandler := kafka.NewHandler(ctx, kafkaReader, nil, 100)
	kHandler.RunFunc(processMessage)
	defer kHandler.Close()
}
