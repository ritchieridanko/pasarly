package models

import "time"

type CreateSession struct {
	ParentID  *int64
	Token     string
	UserAgent string
	IPAddress string
	ExpiresAt time.Time
}
