package utils

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/google/uuid"
	"github.com/ritchieridanko/pasarly/backend/services/gateway/internal/constants"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func CtxAuthID(ctx context.Context) (int64, error) {
	authID, ok := ctx.Value(constants.CtxKeyAuthID).(int64)
	if !ok {
		return 0, errors.New("auth id not provided")
	}

	return authID, nil
}

func CtxRole(ctx context.Context) (string, error) {
	role, ok := ctx.Value(constants.CtxKeyRole).(string)
	if !ok {
		return "", errors.New("role not provided")
	}

	return role, nil
}

func NewUUID() uuid.UUID {
	return uuid.New()
}

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
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
