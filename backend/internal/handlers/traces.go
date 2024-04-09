package handlers

import "go.opentelemetry.io/otel"

var (
	tracer = otel.Tracer("github.com/michaelpeterswa/talvi/backend/internal/handlers")
)
