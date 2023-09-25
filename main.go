package main

import (
	// ... other imports
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
)

func init() {
	// Register your plugin type
	dynamic.Add("custommetrics", func() interface{} {
		return &CustomMetrics{}
	})
}
