package notifier

import (
	"context"
	"gomicrosched/pkg/kafka"

	kafkago "github.com/segmentio/kafka-go"
)

func processMessage(ctx context.Context, w kafka.MessageWriter, msg kafkago.Message) {

}