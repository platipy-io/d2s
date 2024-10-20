package telemetry

import (
	"context"

	"github.com/mdobak/go-xerrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var ErrInitTrace = xerrors.Message("failed initializing tracer")

func InitTrace(name string, opts ...otlptracehttp.Option) error {
	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx, opts...)
	if err != nil {
		return xerrors.WithWrapper(ErrInitTrace, err)
	}
	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName(name)))
	if err != nil {
		return xerrors.WithWrapper(ErrInitTrace, err)
	}
	otel.SetTracerProvider(trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	))
	return nil
}
