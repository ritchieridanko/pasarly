package utils

import (
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var titlecaser = cases.Title(language.English)

func NewUUID() uuid.UUID {
	return uuid.New()
}

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func NormalizeStringPtr(s *string) *string {
	if s == nil {
		return nil
	}
	res := NormalizeString(*s)
	return &res
}

func ToTitlecase(s string) string {
	values := strings.Fields(s)
	if len(values) == 0 {
		return ""
	}

	switch strings.ToLower(values[0]) {
	case "dki", "di":
		return strings.ToUpper(values[0]) + " " + titlecaser.String(strings.Join(values[1:], " "))
	default:
		return titlecaser.String(strings.Join(values, " "))
	}
}

func ToTitlecasePtr(s *string) *string {
	if s == nil {
		return nil
	}
	res := ToTitlecase(*s)
	return &res
}

func TraceErr(s trace.Span, err error, message string) {
	s.RecordError(err)
	s.SetStatus(codes.Error, message)
}

func UnwrapString(sv *wrappers.StringValue) *string {
	if sv != nil {
		return &sv.Value
	}
	return nil
}

func UnwrapTimestamp(ts *timestamp.Timestamp) *time.Time {
	if ts != nil {
		t := ts.AsTime()
		return &t
	}
	return nil
}

func WrapString(s *string) *wrappers.StringValue {
	if s != nil {
		return wrapperspb.String(*s)
	}
	return nil
}

func WrapTime(t *time.Time) *timestamppb.Timestamp {
	if t != nil {
		return timestamppb.New(*t)
	}
	return nil
}
