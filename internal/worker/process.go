package worker

import (
	"context"
	"encoding/json"
	"gomicrosched/internal/domain"
	"gomicrosched/internal/worker/jobworker"
	"gomicrosched/pkg/kafka"
	"log"
)

type JobWorker interface {
	Run() (result string, err error)
}

func processMessage(ctx context.Context, w kafka.MessageWriter, msg kafka.Message) {
	var job domain.Job
	err := json.Unmarshal(msg.Value, &job)
	if err != nil {
		log.Printf("Error unmarshalling message: %v, err - %v", msg, err)
		return
	}

	worker := getWorker(job)
	result, err := worker.Run()

	notification := getEncodedNotification(job, result, err)
	message := kafka.Message{
		Key:   []byte(job.ID.String()),
		Value: notification,
	}
	if err := w.WriteMessages(ctx, message); err != nil {
		log.Printf("Error writing message: %v, err - %v", message, err)
	}
}

func getWorker(job domain.Job) JobWorker {
	switch job.Type {
	case domain.JobTypeUpdate:
		return jobworker.NewUpdateWorker(job)
	case domain.JobTypeFetch:
		return jobworker.NewFetchWorker(job)
	default:
		return nil
	}
}

func getEncodedNotification(job domain.Job, result string, err error) []byte {
	notification := domain.Notification{
		JobID:  job.ID,
		Result: result,
		Error:  "",
	}
	if err != nil {
		notification.Error = err.Error()
	}
	notificationJson, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Error marshalling notification: %v, err - %v", notification, err)
		return nil
	}
	return notificationJson
}
