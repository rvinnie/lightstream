package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
)

const (
	defaultHttpConcurrentRequestsMax = 20
)

type Metrics struct {
	RequestsTotal         *prometheus.CounterVec
	RequestsDuration      *prometheus.HistogramVec
	ConcurrentRequests    prometheus.Gauge
	ConcurrentRequestsMax prometheus.Gauge
}

func CreateMetrics(httpConcurrentRequestsMax float64) *Metrics {
	if httpConcurrentRequestsMax == 0 {
		httpConcurrentRequestsMax = defaultHttpConcurrentRequestsMax
	}

	metrics := Metrics{}

	metrics.RequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gateway_requests_total",
			Help: "The number of all requests to the service",
		},
		[]string{"method", "url", "code"},
	)

	metrics.RequestsDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gateway_requests_duration",
			Help:    "Request processing time",
			Buckets: prometheus.LinearBuckets(0.020, 0.020, 5),
		},
		[]string{"method", "url", "code"},
	)

	metrics.ConcurrentRequests = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_concurrent_requests",
			Help: "The number of inflight requests",
		},
	)

	metrics.ConcurrentRequestsMax = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_concurrent_requests_max",
		},
	)

	metrics.ConcurrentRequestsMax.Set(httpConcurrentRequestsMax)

	prometheus.MustRegister(metrics.RequestsTotal)
	prometheus.MustRegister(metrics.RequestsDuration)
	prometheus.MustRegister(metrics.ConcurrentRequests)
	prometheus.MustRegister(metrics.ConcurrentRequestsMax)

	return &metrics
}

func (m *Metrics) CollectMetrics(method string, url string, statusCode int, duration float64) {
	m.RequestsTotal.With(
		prometheus.Labels{
			"method": method,
			"url":    url,
			"code":   strconv.Itoa(statusCode),
		}).Inc()

	m.RequestsDuration.With(
		prometheus.Labels{
			"method": method,
			"url":    url,
			"code":   strconv.Itoa(statusCode),
		}).Observe(duration)
}
