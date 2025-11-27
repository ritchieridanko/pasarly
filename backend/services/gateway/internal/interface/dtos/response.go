package dtos

import "time"

type Response[T any] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	RequestID string    `json:"request_id"`
	Page      *int      `json:"page,omitempty"`
	PageSize  *int      `json:"page_size,omitempty"`
	Total     *int      `json:"total,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
