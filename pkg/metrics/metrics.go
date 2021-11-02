package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	APICallErrorsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "yaae_api_call_errors_total",
		Help: "Count of errors calling APIs",
	}, []string{"service", "api"})
	ScrapeDurationHistogram = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "yaae_scrape_duration_seconds",
		Help:    "Histogram of response time for handler in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"scraper", "status"})
)

func InitMetrics() {
	prometheus.MustRegister(APICallErrorsTotal)
	prometheus.MustRegister(ScrapeDurationHistogram)
}
