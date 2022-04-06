package es2

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("github.com/contextgg/pkg/es2")
