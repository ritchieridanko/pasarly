package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/channels"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/repositories"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"github.com/ritchieridanko/pasarly/backend/shared/events/v1"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/proto"
)

const authErrTracer string = "handler.auth"

type AuthHandler struct {
	timeout time.Duration
	er      repositories.EventRepository
	ec      channels.EmailChannel
}

func NewAuthHandler(er repositories.EventRepository, ec channels.EmailChannel, timeout time.Duration) *AuthHandler {
	return &AuthHandler{er: er, ec: ec, timeout: timeout}
}

func (h *AuthHandler) OnAuthCreated(ctx context.Context, m kafka.Message) error {
	ctx, span := otel.Tracer(authErrTracer).Start(ctx, "OnAuthCreated")
	defer span.End()

	var evt events.AuthCreated
	if err := proto.Unmarshal(m.Value, &evt); err != nil {
		e := fmt.Errorf("failed to handle message: %w", err)
		utils.TraceErr(span, e, ce.MsgInternalServer)
		return e
	}

	event, err := h.er.GetEventByID(ctx, evt.GetEventId())
	if err != nil {
		return err
	}

	// Idempotency check
	if event == nil {
		data := models.CreateEvent{
			ID:   evt.GetEventId(),
			Type: constants.EventTopicAuthCreated,
		}

		if err := h.er.CreateEvent(ctx, &data); err != nil {
			return err
		}
	}
	if event != nil {
		if event.CompletedAt != nil {
			return nil // already completed -> skip
		}

		if time.Since(event.ProcessedAt).Seconds() < h.timeout.Seconds() {
			e := fmt.Errorf("failed to handle message: %w", ce.ErrEventOnProcess)
			utils.TraceErr(span, e, ce.MsgInternalServer)
			return e
		}
	}

	if err := h.ec.SendWelcome(ctx, evt.GetEmail(), evt.GetToken()); err != nil {
		return err
	}

	return h.er.SetCompleted(ctx, evt.GetEventId())
}
