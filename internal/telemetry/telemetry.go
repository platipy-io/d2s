package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var Provider *trace.TracerProvider

func init() {
	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx, otlptracehttp.WithInsecure())
	if err != nil {
		panic(err)
	}

	res, err := resource.New(ctx, resource.WithAttributes(semconv.ServiceName("d2s")))
	if err != nil {
		panic(err)
	}
	Provider = trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)
}
