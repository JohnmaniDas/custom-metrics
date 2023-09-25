package main

import (
	"context"
	"net/http"

	"github.com/JohnmaniDas/custom-metrics/middleware"
)

// Config holds the plugin configuration.
type Config struct{}

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// CustomMetrics holds a plugin instance.
type CustomMetrics struct {
	next http.Handler
	name string
	conf *Config
}

// New creates a new plugin instance.
func New(ctx context.Context, next http.Handler, conf *Config, name string) (http.Handler, error) {
	return middleware.New(ctx, next, conf, name)
}

// Implement the ServeHTTP method for your CustomMetrics type here.
func (c *CustomMetrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement your middleware logic here.
}
