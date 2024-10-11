package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

var Provider *trace.TracerProvider

func init() {
	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("localhost:5080"),
		otlptracehttp.WithURLPath("/api/default/v1/traces"),
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithHeaders(map[string]string{
			"Authorization": "Basic cm9vdEBleGFtcGxlLmNvbTpEMlp5Qk0wdGNaV3BoenZM",
			"stream-name":   "default",
		}),
	)
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

	// Create the OTLP log exporter that sends logs to configured destination
	logExporter, err := otlploghttp.New(ctx,
		otlploghttp.WithInsecure(),
		otlploghttp.WithEndpoint("localhost:5080"),
		otlploghttp.WithURLPath("/api/default/v1/logs"),
		otlploghttp.WithHeaders(map[string]string{
			"Authorization": "Basic cm9vdEBleGFtcGxlLmNvbTpEMlp5Qk0wdGNaV3BoenZM",
			"stream-name":   "default",
		}),
	)
	if err != nil {
		panic("failed to initialize exporter")
	}

	// Create the logger provider
	lp := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)
	global.SetLoggerProvider(lp)
}
