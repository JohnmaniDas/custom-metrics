package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
)

type CustomMetricsMiddleware struct {
	next     http.Handler
	requests prometheus.Counter
	errors   prometheus.Counter
	duration prometheus.Summary
}

func New(ctx context.Context, next http.Handler, config dynamic.Middleware, name string) (http.Handler, error) {
	registry := prometheus.NewRegistry()

	requests := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "requests_total",
		Help: "The total number of requests.",
	})

	errors := prometheus.NewCounter(prometheus.CounterOpts{
		Name: "errors_total",
		Help: "The total number of errors.",
	})

	duration := prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "duration_per_request",
		Help: "The duration of each request.",
	})

	registry.MustRegister(requests, errors, duration)

	return &CustomMetricsMiddleware{
		next:     next,
		requests: requests,
		errors:   errors,
		duration: duration,
	}, nil
}

func (m *CustomMetricsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	m.next.ServeHTTP(w, r)

	duration := time.Since(start)

	m.requests.Inc()

	if r.Response.StatusCode >= 400 {
		m.errors.Inc()
	}

	m.duration.Observe(float64(duration.Milliseconds()))
}
