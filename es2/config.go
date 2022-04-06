package es2

import (
	"context"

	"github.com/contextgg/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.uber.org/zap"
)

// newExporter returns a console exporter.
func newExporter(url string) (trace.SpanExporter, error) {
	// Create the Jaeger exporter
	return jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
}

// newResource returns a resource describing this application.
func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("fib"),
			semconv.ServiceVersionKey.String("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

type App interface {
	CommandBus
	Close()
}

type app struct {
	CommandBus

	log            logger.Logger
	tracerProvider *trace.TracerProvider
}

func (a *app) Close() {
	if err := a.tracerProvider.Shutdown(context.Background()); err != nil {
		a.log.Fatal("Shutdown", "err", err)
	}
}

func Build(commandRegistry CommandRegistry) (App, error) {
	z := zap.NewNop()
	l := logger.NewLogger(z)

	exp, err := newExporter("http://localhost:14268/api/traces")
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(newResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	commandBus := NewCommandBus(commandRegistry)

	return &app{
		CommandBus:     commandBus,
		log:            l,
		tracerProvider: tp,
	}, nil
}
