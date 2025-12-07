package utils

import (
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func TraceErr(s trace.Span, err error, message string) {
	s.RecordError(err)
	s.SetStatus(codes.Error, message)
}
