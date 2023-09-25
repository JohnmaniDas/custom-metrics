package custom_metrics

import (
	"context"
	"log"
	"net/http"
)

type CustomMetrics struct {
	next http.Handler
	name string
	conf *Config
}

type Config struct{}

func CreateConfig() *Config {
	return &Config{}
}

// Note: We've changed the signature of the New function here
func New(ctx context.Context, next http.Handler, name string) (http.Handler, error) {
	return &CustomMetrics{
		next: next,
		name: name,
		conf: CreateConfig(),
	}, nil
}

func (c *CustomMetrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Log information about the incoming request.
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)

	// Call the next handler in the chain.
	c.next.ServeHTTP(w, r)
}
