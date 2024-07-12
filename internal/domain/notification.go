package domain

import "github.com/google/uuid"

type Notification struct {
	JobID  uuid.UUID
	Result string
	Error  string
	Email  string
}
