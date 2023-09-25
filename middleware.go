package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type CustomMetricsMiddleware struct {
	next     http.Handler
	registry *prometheus.Registry
	requests prometheus.Counter
	errors   prometheus.Counter
	duration prometheus.Summary
}

func NewCustomMetricsMiddleware() *CustomMetricsMiddleware {
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
		next:     nil,
		registry: registry,
		requests: requests,
		errors:   errors,
		duration: duration,
	}
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

func (m *CustomMetricsMiddleware) SetNext(next http.Handler) {
	m.next = next
}
