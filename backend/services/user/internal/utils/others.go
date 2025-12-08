package utils

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func TraceErr(s trace.Span, err error, message string) {
	s.RecordError(err)
	s.SetStatus(codes.Error, message)
}

func UnwrapString(value *wrappers.StringValue) *string {
	if value != nil {
		return &value.Value
	}
	return nil
}

func UnwrapTimestamp(value *timestamp.Timestamp) *time.Time {
	if value != nil {
		t := value.AsTime()
		return &t
	}
	return nil
}

func WrapString(value *string) *wrappers.StringValue {
	if value != nil {
		return wrapperspb.String(*value)
	}
	return nil
}

func WrapTime(value *time.Time) *timestamppb.Timestamp {
	if value != nil {
		return timestamppb.New(*value)
	}
	return nil
}
