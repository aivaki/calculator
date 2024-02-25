package types

import (
	"time"
)

type StatusData struct {
	FreeWorkers int      `json:"free_workers"`
	MaxWorkers  int      `json:"max_workers"`
	Expressions []string `json:"expressions_in_process"`
}

type Expression struct {
	Status,
	Result,
	Expression string
	CreatedAt,
	CompletedAt time.Time
}

type Job struct {
	ID         string
	Expression string
	Result     any
}