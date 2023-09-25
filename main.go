package main

import (
	"github.com/traefik/traefik/v2/pkg/config/dynamic"
)

func init() {
	dynamic.Add("custommetrics", func() interface{} {
		return &CustomMetricsMiddleware{}
	})
}
