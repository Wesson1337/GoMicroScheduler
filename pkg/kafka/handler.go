package kafka

import (
	"context"
	"fmt"
	"log"
	"sync"

	kafkago "github.com/segmentio/kafka-go"
)

const logCount = 1000

type Handler interface {
	RunFunc(msgProcessor ProcessFunc)
	Close()
}

type handler struct {
	ctx            context.Context
	reader         MessageReader
	writer         MessageWriter
	countMutex     sync.Mutex
	processedCount int
	processSem     chan struct{}
}

type ProcessFunc func(ctx context.Context, w MessageWriter, msg kafkago.Message)

func NewHandler(ctx context.Context, r MessageReader, w MessageWriter, maxProcMessage int) Handler {
	return &handler{
		ctx:            ctx,
		reader:         r,
		writer:         w,
		countMutex:     sync.Mutex{},
		processedCount: 0,
		processSem:     make(chan struct{}, maxProcMessage),
	}
}

func (handler *handler) RunFunc(msgProcessor ProcessFunc) {
	log.Printf("Started reading and processing messages with reader: %#v, writer: %#v\n", handler.reader, handler.writer)
	for {
		select {
		case <-handler.ctx.Done():
			log.Println("Context is cancelled. Waiting for all msgProcessors..")
			handler.cancelProcessing()
			return
		default:
			handler.handleProcessing(msgProcessor)
		}

	}
}

func (handler *handler) cancelProcessing() {
	for i := 0; i < cap(handler.processSem); i++ {
		handler.processSem <- struct{}{}
	}
	close(handler.processSem)
}

func (handler *handler) handleProcessing(msgProcessor ProcessFunc) {
	msg, err := handler.reader.ReadMessage(context.WithoutCancel(handler.ctx))
	if err != nil {
		panic(fmt.Sprintf("error reading message: %v", err))
	}

	handler.runProcessFunc(msg, msgProcessor)
}

func (handler *handler) runProcessFunc(msg kafkago.Message, msgProcessor ProcessFunc) {
	handler.processSem <- struct{}{}
	go func(msg kafkago.Message) {
		defer func() { <-handler.processSem }()

		msgProcessor(handler.ctx, handler.writer, msg)

		handler.incAndLogProcCount()
	}(msg)
}

func (handler *handler) incAndLogProcCount() {
	handler.countMutex.Lock()
	defer handler.countMutex.Unlock()

	handler.processedCount++
	if handler.processedCount%logCount == 0 {
		log.Printf("Processed %d messages", handler.processedCount)
	}
}

func (handler *handler) Close() {
	if handler.reader != nil {
		handler.reader.Close()
	}
	if handler.writer != nil {
		handler.writer.Close()
	}
}
