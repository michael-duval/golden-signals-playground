package server

import "github.com/prometheus/client_golang/prometheus"

var (
	reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests processed, labeled by status code and HTTP route.",
		},
		[]string{"route", "code"},
	)
	latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds_bucket",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"route"},
	)
	inflight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "inflight_requests",
			Help: "Currently in-flight requests.",
		},
	)
)

func init() {
	prometheus.MustRegister(reqs, latency, inflight)
}
