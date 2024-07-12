package notifier

import (
	"context"
	"gomicrosched/pkg/env"
	"gomicrosched/pkg/kafka"
	"gomicrosched/pkg/logger"
	"gomicrosched/pkg/smtp"
	"os/signal"
	"strings"
	"syscall"
)

func init() {
	logger.ConfigureLogger()
}

func Run() {
	kConfig := kafka.KafkaConfig{
		Topic:    env.MustGetEnv("KAFKA_TOPIC_NOTIFIER"),
		Brokers:  strings.Split(env.MustGetEnv("KAFKA_BROKER"), ","),
		Username: env.MustGetEnv("KAFKA_USERNAME"),
		Password: env.MustGetEnv("KAFKA_PASSWORD"),
		GroupID:  "notifier",
	}

	smtpConfig := smtp.SMTPConfig{
		Addr:     env.MustGetEnv("SMTP_ADDR"),
		Username: env.MustGetEnv("SMTP_USERNAME"),
		Password: env.MustGetEnv("SMTP_PASSWORD"),
	}

	kReader := kafka.NewReader(kConfig)
	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	smtpClient := smtp.NewSMTPClient(smtpConfig)

	kHandler := kafka.NewHandler(ctx, kReader, nil, 10000)
	defer kHandler.Close()
	kHandler.RunFunc(processMessage(smtpClient))
}
