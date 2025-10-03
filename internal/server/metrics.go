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
			Name: "http_request_duration_seconds",
			Help: "HTTP request latency in seconds.",
			Buckets: []float64{
				0.005, 0.010, 0.020, 0.030, 0.040, 0.050,
				0.060, 0.070, 0.080, 0.090, 0.100,
				0.110, 0.120, 0.130, 0.140, 0.150,
				0.160, 0.170, 0.180, 0.190, 0.200,
				0.210, 0.220, 0.230, 0.240, 0.250,
				0.300, 0.400, 0.500, 0.750, 1.000, 1.500, 2.000, 3.000, 5.000,
			},
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
