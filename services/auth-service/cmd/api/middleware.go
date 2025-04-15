package main

import (
	"net/http"
	"strconv"
	"time"
)

func (app *application) metricsMiddleware(next http.Handler) http.Handler {
	// TODO: normalize paths (e.g., replace /v1/users/123 with /v1/users/:id) to reduce high cardinality in metrics.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		method := r.Method

		start := time.Now()

		// In-flight request increment
		app.metrics.inFlightRequests.Inc()
		defer app.metrics.inFlightRequests.Dec()

		// Capture request size (from Content-Length)
		reqSize := float64(0)
		if r.ContentLength > 0 {
			reqSize = float64(r.ContentLength)
		}
		app.metrics.requestSize.WithLabelValues(path, method).Observe(reqSize)

		// Wrap the ResponseWriter to capture status code and response size
		rw := &responseWriterWrapper{ResponseWriter: w, statusCode: 200}

		// Call next handler
		next.ServeHTTP(rw, r)

		// Observe duration
		duration := time.Since(start).Seconds()
		app.metrics.requestDuration.WithLabelValues(path, method).Observe(duration)

		// Count response size
		app.metrics.responseSize.WithLabelValues(path, method).Observe(float64(rw.size))

		// Count requests
		app.metrics.requestCount.WithLabelValues(path, method).Inc()

		// Count errors
		if rw.statusCode >= 400 {
			app.metrics.errorCount.WithLabelValues(path, method, strconv.Itoa(rw.statusCode)).Inc()
		}
	})
}

// Custom response writer to capture response size and status code
type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriterWrapper) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}
