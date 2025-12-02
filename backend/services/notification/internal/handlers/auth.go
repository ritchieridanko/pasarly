package handlers

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/channels"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"github.com/ritchieridanko/pasarly/backend/shared/events/v1"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/proto"
)

const authErrTracer string = "handler.auth"

type AuthHandler struct {
	ec channels.EmailChannel
}

func NewAuthHandler(ec channels.EmailChannel) *AuthHandler {
	return &AuthHandler{ec: ec}
}

func (h *AuthHandler) AuthCreated(ctx context.Context, m kafka.Message) error {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "AuthCreated")
	defer span.End()

	var evt events.AuthCreated
	if err := proto.Unmarshal(m.Value, &evt); err != nil {
		e := fmt.Errorf("failed to handle message: %w", err)
		utils.TraceErr(span, e, ce.MsgInternalServer)
		return e
	}

	return h.ec.SendWelcome(ctx, evt.GetEmail(), evt.GetToken())
}
