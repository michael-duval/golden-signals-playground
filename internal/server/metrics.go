package server

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reqs = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests processed, labeled by status code and HTTP route.",
		},
		[]string{"method", "code"},
	)
	latency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
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

func InstrumentedHandler(routeName string, h http.Handler) http.Handler {
	return promhttp.InstrumentHandlerInFlight(inflight,
		promhttp.InstrumentHandlerDuration(latency,
			promhttp.InstrumentHandlerCounter(reqs, h)))
}
