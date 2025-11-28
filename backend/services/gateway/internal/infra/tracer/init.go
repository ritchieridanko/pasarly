package tracer

import (
	"context"
	"fmt"

	"github.com/ritchieridanko/pasarly/backend/services/gateway/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.uber.org/zap"
)

type Tracer struct {
	Cleanup func()
}

func Init(cfg *configs.Config, l *zap.Logger) (*Tracer, error) {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(cfg.Tracer.Endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(cfg.App.Name),
			),
		),
	)

	otel.SetTracerProvider(tp)

	l.Sugar().Infof("✅ [TRACER] initialized (endpoint=%s)", cfg.Endpoint)
	return &Tracer{Cleanup: func() { _ = tp.Shutdown(ctx) }}, nil
}
