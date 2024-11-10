package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelmetric "go.opentelemetry.io/otel/metric"
)

var (
	meter = otel.Meter("github.com/michaelpeterswa/talvi/backend/internal/middleware")

	requestsCounter otelmetric.Int64Counter
)

func init() {
	var err error

	requestsCounter, err = meter.Int64Counter(
		"requests",
		otelmetric.WithDescription("number of requests to the movies service"),
	)

	if err != nil {
		fmt.Println("Failed to create requests counter")
		os.Exit(1)
	}
}

func RequestsCounterMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestsCounter.Add(r.Context(), 1, otelmetric.WithAttributes(
				attribute.String("path", r.URL.Path),
				attribute.String("method", r.Method),
			))

			next.ServeHTTP(w, r)
		})
	}
}
