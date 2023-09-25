package main

import (
	"context"
	"net/http"

	"github.com/traefik/traefik/v2/pkg/config/dynamic"
)

// CustomMetrics holds a plugin instance.
type CustomMetrics struct {
	next http.Handler
	name string
	conf *Config
}

// New creates a new plugin instance.
func New(ctx context.Context, next http.Handler, conf *dynamic.Middleware, name string) (http.Handler, error) {
	return &CustomMetrics{
		next: next,
		name: name,
		conf: conf,
	}, nil
}

// Implement the ServeHTTP method for your CustomMetrics type here.
func (c *CustomMetrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement your middleware logic here.
}
