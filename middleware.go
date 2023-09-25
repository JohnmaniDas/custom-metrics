package custommetrics

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

// Config holds the configuration for your middleware.
type Config struct {
	// Add middleware-specific configuration here.
}

// CreateConfig initializes the default Config.
func CreateConfig() *Config {
	return &Config{}
}

// CustomMetricsMiddleware is your custom middleware struct.
type CustomMetricsMiddleware struct {
	next       http.Handler
	requestCnt uint64
	errorCnt   uint64
	totalDur   uint64
}

// New creates a new instance of your middleware.
// It's the function that Traefik will call in order to create an instance of your middleware.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// Here, you can initialize your middleware with the config.
	return &CustomMetricsMiddleware{next: next}, nil
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (m *CustomMetricsMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	rw := &responseWriterWrapper{ResponseWriter: w}
	m.next.ServeHTTP(rw, r)

	duration := time.Since(start).Seconds()

	atomic.AddUint64(&m.requestCnt, 1)
	atomic.AddUint64(&m.totalDur, uint64(duration*1000))

	if rw.statusCode >= 400 {
		atomic.AddUint64(&m.errorCnt, 1)
	}
}

func (m *CustomMetricsMiddleware) GetMetricsString() string {
	reqs := atomic.LoadUint64(&m.requestCnt)
	errs := atomic.LoadUint64(&m.errorCnt)
	dur := atomic.LoadUint64(&m.totalDur)

	return fmt.Sprintf("requests_total %d\nerrors_total %d\nduration_milliseconds_total %d", reqs, errs, dur)
}

// This MetricsHandler may not be needed unless you want a separate endpoint just for showing metrics.
// If you do, you'll have to set it up separately in your main function or application.
func MetricsHandler(m *CustomMetricsMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(m.GetMetricsString()))
	}
}

func main() {
	mux := http.NewServeMux()

	middleware := New(context.Background(), mux, CreateConfig(), "custom-metrics")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	http.HandleFunc("/metrics", MetricsHandler(middleware.(*CustomMetricsMiddleware)))
	http.ListenAndServe(":9090", middleware)
}
