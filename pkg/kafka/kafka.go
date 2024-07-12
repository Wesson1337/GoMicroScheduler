package kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

type Message = kafka.Message

type KafkaConfig struct {
	Brokers  []string
	Topic    string
	Username string
	Password string
	GroupID  string
}

type MessageWriter interface {
	WriteMessages(ctx context.Context, msgs ...Message) error
	Close() error
}

type MessageReader interface {
	ReadMessage(ctx context.Context) (Message, error)

	FetchMessage(ctx context.Context) (Message, error)

	CommitMessages(ctx context.Context, msgs ...Message) error

	Close() error
}

func NewWriter(config KafkaConfig) MessageWriter {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	wConfig := kafka.WriterConfig{
		Brokers:      config.Brokers,
		Topic:        config.Topic,
		Dialer:       dialer,
		Balancer:     &kafka.LeastBytes{},
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		BatchTimeout: 100 * time.Millisecond,
		RequiredAcks: -1,
	}

	return kafka.NewWriter(wConfig)
}

func NewReader(config KafkaConfig) MessageReader {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
	}

	rConfig := kafka.ReaderConfig{
		Brokers:        config.Brokers,
		Topic:          config.Topic,
		GroupID:        config.GroupID,
		Dialer:         dialer,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
	}

	return kafka.NewReader(rConfig)
}
