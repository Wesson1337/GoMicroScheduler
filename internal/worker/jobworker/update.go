package jobworker

import (
	"gomicrosched/internal/domain"

	"github.com/google/uuid"
)

type updateWorker struct {
	jobId uuid.UUID
}

func NewUpdateWorker(job domain.Job) *updateWorker {
	return &updateWorker{
		jobId: job.ID,
	}
}

func (w *updateWorker) Run() (result string, err error) {
	// some update logic here
	// for example SQL query UPDATE mytable SET status = 'done' WHERE id = w.JobId
	return "Data updated successfully", nil
}
