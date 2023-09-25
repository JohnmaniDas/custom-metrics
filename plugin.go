package main

import (
	"context"
	"net/http"

	"github.com/JohnmaniDas/custom-metrics/middleware"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
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
func New(ctx context.Context, next http.Handler, conf *dynamic.Middleware, name string) (http.Handler, error) {
	return middleware.New(ctx, next, conf, name)
}