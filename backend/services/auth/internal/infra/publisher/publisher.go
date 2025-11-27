package publisher

import (
	"context"

	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/constants"
	"github.com/ritchieridanko/pasarly/backend/services/auth/internal/infra/logger"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/proto"
)

type Publisher struct {
	writer *kafka.Writer
	logger *logger.Logger
}

func NewPublisher(w *kafka.Writer, l *logger.Logger) *Publisher {
	return &Publisher{writer: w, logger: l}
}

func (p *Publisher) Publish(ctx context.Context, key string, m proto.Message) error {
	value, err := proto.Marshal(m)
	if err != nil {
		return err
	}

	traceID := trace.SpanFromContext(ctx).SpanContext().TraceID().String()

	var requestID string
	if v := ctx.Value(constants.CtxKeyRequestID); v != nil {
		requestID, _ = v.(string)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: value,
		Headers: []kafka.Header{
			{Key: "trace_id", Value: []byte(traceID)},
			{Key: "correlation_id", Value: []byte(requestID)},
			{Key: "content_type", Value: []byte("application/x-protobuf")},
		},
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		p.logger.Sugar().Warnf("failed to publish message (topic=%s, key=%s): %s", p.writer.Topic, key, err.Error())
	}

	return err
}
