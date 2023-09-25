package custom_metrics

import (
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

type CustomMetricsMiddleware struct {
	next       http.Handler
	requestCnt uint64
	errorCnt   uint64
	totalDur   uint64
}

func New(next http.Handler) *CustomMetricsMiddleware {
	return &CustomMetricsMiddleware{
		next: next,
	}
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

func MetricsHandler(m *CustomMetricsMiddleware) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(m.GetMetricsString()))
	}
}

func main() {
	mux := http.NewServeMux()

	middleware := New(mux)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, world!"))
	})

	http.HandleFunc("/metrics", MetricsHandler(middleware))
	http.ListenAndServe(":9090", middleware)
}
