package main

import (
	"context"
	"net/http"

	"github.com/JohnmaniDas/custom-metrics/middleware"
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
)

// CreateConfig creates and initializes the plugin configuration.
func CreateConfig() *Config {
	return &Config{}
}

// New creates a new plugin instance.
func New(ctx context.Context, next http.Handler, conf *dynamic.Middleware, name string) (http.Handler, error) {
	return middleware.New(ctx, next, conf, name)
}
