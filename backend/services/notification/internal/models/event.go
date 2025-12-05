package models

import "time"

type Event struct {
	ID          string
	Type        string
	ProcessedAt time.Time
	CompletedAt *time.Time
}

type CreateEvent struct {
	ID   string
	Type string
}
