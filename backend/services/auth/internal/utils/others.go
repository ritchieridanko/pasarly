package utils

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/constants"
	"google.golang.org/grpc/metadata"
)

func CtxRequestMeta(ctx context.Context) (userAgent, ipAddress string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", ""
	}

	if vals := md.Get(constants.CtxKeyUserAgent); len(vals) > 0 {
		userAgent = vals[0]
	}
	if vals := md.Get(constants.CtxKeyIPAddress); len(vals) > 0 {
		ipAddress = vals[0]
	}

	return
}

func NewUUID() uuid.UUID {
	return uuid.New()
}

func NormalizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
