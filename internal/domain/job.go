package domain

import "github.com/google/uuid"

type JobType int

const (
	JobTypeFetch JobType = iota
	JobTypeUpdate
)

type Job struct {
	ID     uuid.UUID
	Type   JobType
	Email  string
	Period string
}
