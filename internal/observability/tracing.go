package observability

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
)

var Propagator = propagation.TraceContext{}

func InitTracer(serviceName string) func() {
	ctx := context.Background()

	// Create OTLP gRPC exporter (send to Tempo)
	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),                 // No TLS for local dev
		otlptracegrpc.WithEndpoint("localhost:4317"), // Tempo gRPC port
		otlptracegrpc.WithDialOption(grpc.WithBlock()),
		// otlptracegrpc.WithDialOption(grpc.WithBlock()),
	)
	if err != nil {
		log.Fatalf("failed to create exporter: %v", err)
	}

	// Resource describes your service
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	// Tracer provider setup
	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(Propagator)
	log.Println("Started tracer ....")

	// Shutdown func
	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("Error shutting down tracer provider: %v", err)
		}
	}
}
