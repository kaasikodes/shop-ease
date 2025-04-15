package main

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	requestCount     *prometheus.CounterVec
	errorCount       *prometheus.CounterVec   //Track how many requests failed (e.g. 4xx or 5xx).
	requestDuration  *prometheus.HistogramVec //Measure how long each HTTP request takes. This helps identify slow endpoints.
	requestSize      *prometheus.HistogramVec //Tracks the size of incoming requests. Useful for spotting unexpectedly large payloads.
	responseSize     *prometheus.HistogramVec //Track size of outgoing responses. Helps monitor compression and heavy responses.
	inFlightRequests prometheus.Gauge         //Track the number of currently active requests being processed.
}

func NewMetrics(reg *prometheus.Registry) *metrics {
	m := &metrics{
		requestCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests.",
			},
			[]string{"path", "method"},
		),
		errorCount: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_error_count_total",
				Help: "Total number of HTTP error responses.",
			},
			[]string{"path", "method", "status"},
		),
		requestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path", "method"},
		),
		requestSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_size_bytes",
				Help:    "Size of HTTP requests in bytes.",
				Buckets: prometheus.ExponentialBuckets(100, 2, 10),
			},
			[]string{"path", "method"},
		),
		responseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "Size of HTTP responses in bytes.",
				Buckets: prometheus.ExponentialBuckets(100, 2, 10),
			},
			[]string{"path", "method"},
		),
		inFlightRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "in_flight_requests",
				Help: "Number of in-flight HTTP requests.",
			},
		),
	}

	reg.MustRegister(
		m.requestCount,
		m.errorCount,
		m.requestDuration,
		m.requestSize,
		m.responseSize,
		m.inFlightRequests,
	)

	return m
}

// gauge
// counter
// histogram
