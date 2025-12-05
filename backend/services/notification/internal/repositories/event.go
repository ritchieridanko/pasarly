package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/infra/database"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/models"
	"github.com/ritchieridanko/pasarly/backend/services/notification/internal/utils"
	"github.com/ritchieridanko/pasarly/backend/shared/ce"
	"go.opentelemetry.io/otel"
)

const eventErrTracer string = "repository.event"

type EventRepository interface {
	CreateEvent(ctx context.Context, data *models.CreateEvent) (err error)
	GetEventByID(ctx context.Context, eventID string) (event *models.Event, err error)
	SetCompleted(ctx context.Context, eventID string) (err error)
}

type eventRepository struct {
	database *database.Database
}

func NewEventRepository(db *database.Database) EventRepository {
	return &eventRepository{database: db}
}

func (r *eventRepository) CreateEvent(ctx context.Context, data *models.CreateEvent) error {
	ctx, span := otel.Tracer(eventErrTracer).Start(ctx, "CreateEvent")
	defer span.End()

	query := "INSERT INTO events (event_id, event_type) VALUES ($1, $2)"

	if err := r.database.Execute(ctx, query, data.ID, data.Type); err != nil {
		e := fmt.Errorf("failed to create event: %w", err)
		utils.TraceErr(span, e, ce.MsgInternalServer)
		return e
	}

	return nil
}

func (r *eventRepository) GetEventByID(ctx context.Context, eventID string) (*models.Event, error) {
	ctx, span := otel.Tracer(eventErrTracer).Start(ctx, "GetEventByID")
	defer span.End()

	query := `
		SELECT event_id, event_type, processed_at, completed_at
		FROM events
		WHERE event_id = $1
	`
	if r.database.InTx(ctx) {
		query += " FOR UPDATE"
	}

	row := r.database.QueryRow(ctx, query, eventID)

	var e models.Event
	if err := row.Scan(&e.ID, &e.Type, &e.ProcessedAt, &e.CompletedAt); err != nil {
		if errors.Is(err, ce.ErrDBReturnNoRows) {
			return nil, nil
		}

		e := fmt.Errorf("failed to fetch event by id: %w", err)
		utils.TraceErr(span, e, ce.MsgInternalServer)
		return nil, e
	}

	return &e, nil
}

func (r *eventRepository) SetCompleted(ctx context.Context, eventID string) error {
	ctx, span := otel.Tracer(eventErrTracer).Start(ctx, "SetCompleted")
	defer span.End()

	query := "UPDATE events SET completed_at = NOW() WHERE event_id = $1"

	if err := r.database.Execute(ctx, query, eventID); err != nil {
		e := fmt.Errorf("failed to set event completed: %w", err)
		utils.TraceErr(span, e, ce.MsgInternalServer)
		return e
	}

	return nil
}
