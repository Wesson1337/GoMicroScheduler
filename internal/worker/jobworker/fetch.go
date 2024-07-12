package jobworker

import (
	"gomicrosched/internal/domain"

	"github.com/google/uuid"
)

type fetchWorker struct {
	jobId uuid.UUID
}

func NewFetchWorker(job domain.Job) *fetchWorker {
	return &fetchWorker{
		jobId: job.ID,
	}
}

func (w *fetchWorker) Run() (result string, err error) {
	// some fetch logic here
	// for example sql query SELECT * FROM mytable WHERE id = w.JobId
	return "Data to process - 1, 2, 3", nil
}
