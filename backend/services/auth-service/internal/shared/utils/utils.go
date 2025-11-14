package utils

import (
	"strings"

	"github.com/google/uuid"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
