package tracer

import (
	"context"
	"fmt"

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

func Init(appName, endpoint string, l *zap.Logger) (*Tracer, error) {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(endpoint),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tracer: %w", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(appName),
			),
		),
	)

	otel.SetTracerProvider(tp)

	l.Sugar().Infof("✅ [TRACER] initialized (app_name=%s, endpoint=%s)", appName, endpoint)
	return &Tracer{Cleanup: func() { _ = tp.Shutdown(ctx) }}, nil
}
