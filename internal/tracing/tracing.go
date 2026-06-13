package tracing

import (
	"context"
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

const defaultServiceName = "go-server-template"

// ServiceName returns OTEL_SERVICE_NAME or a hardcoded fallback. Used as the
// instrumentation scope name for the chi server middleware.
func ServiceName() string {
	if v := os.Getenv("OTEL_SERVICE_NAME"); v != "" {
		return v
	}
	return defaultServiceName
}

// Setup installs the global tracer provider and propagator. The returned
// shutdown function flushes and stops the provider; it is always non-nil and
// safe to call even when tracing is disabled.
func Setup(ctx context.Context) (func(context.Context) error, error) {
	// Always install the W3C trace-context + baggage propagator so an incoming
	// `traceparent` (e.g. from the k8s gateway) is parsed and forwarded on
	// outgoing calls, regardless of whether we export spans ourselves.
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_ENDPOINT")
	if endpoint == "" {
		endpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}
	if endpoint == "" {
		return func(context.Context) error { return nil }, nil
	}

	exp, err := newExporter(ctx)
	if err != nil {
		return nil, fmt.Errorf("create otlp exporter: %w", err)
	}

	res, err := newResource(ctx)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		// No explicit sampler: SDK default is parentbased_always_on, and
		// OTEL_TRACES_SAMPLER / OTEL_TRACES_SAMPLER_ARG are honored.
	)
	otel.SetTracerProvider(tp)

	return func(shutdownCtx context.Context) error {
		_ = tp.ForceFlush(shutdownCtx)
		return tp.Shutdown(shutdownCtx)
	}, nil
}

func newExporter(ctx context.Context) (sdktrace.SpanExporter, error) {
	protocol := os.Getenv("OTEL_EXPORTER_OTLP_TRACES_PROTOCOL")
	if protocol == "" {
		protocol = os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
	}
	switch protocol {
	case "http/protobuf", "http":
		return otlptracehttp.New(ctx)
	case "", "grpc":
		return otlptracegrpc.New(ctx)
	default:
		return nil, fmt.Errorf("unsupported OTEL_EXPORTER_OTLP_PROTOCOL: %q", protocol)
	}
}

func newResource(ctx context.Context) (*resource.Resource, error) {
	opts := []resource.Option{
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
		resource.WithFromEnv(),
	}
	// Provide a fallback service.name only when the env var isn't set, so
	// resource.WithFromEnv() always wins when present.
	if os.Getenv("OTEL_SERVICE_NAME") == "" {
		opts = append(opts, resource.WithAttributes(semconv.ServiceName(defaultServiceName)))
	}
	return resource.New(ctx, opts...)
}
