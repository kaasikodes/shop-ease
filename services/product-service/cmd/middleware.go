package main

import (
	"net/http"
	"regexp"
	"strconv"
	"time"
)

func (app *application) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		normalizedPath := normalizePath(r.URL.Path)
		// path := r.URL.Path
		path := normalizedPath //normalize paths (e.g., replace /v1/users/123 with /v1/users/:id) to reduce high cardinality in metrics.
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
func normalizePath(path string) string {
	// You can add more normalization patterns as needed
	userIdPattern := regexp.MustCompile(`/v1/users/\d+`)
	orderIdPattern := regexp.MustCompile(`/v1/orders/\d+`)
	uuidPattern := regexp.MustCompile(`/[0-9a-fA-F\-]{36}`)

	// Apply them in order
	path = userIdPattern.ReplaceAllString(path, "/v1/users/:id")
	path = orderIdPattern.ReplaceAllString(path, "/v1/orders/:id")
	path = uuidPattern.ReplaceAllString(path, "/:uuid")

	return path
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
