package types

import (
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cashapp/yet-another-aws-exporter/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// ScrapeResult is the struct that must be returned from a scraper.
// The order of Labels should match the order defined in the Prometheus metric.
type ScrapeResult struct {
	Labels []string
	Value  float64
	Type   prometheus.ValueType
}

// Scraper contains all the logic for a single Prometheus metric
// and the logic for how to scrape that metric. There is a one-to-one
// relationship between the metrics exported and scrapers.
type Scraper struct {
	Metrics        map[string]*Metric
	ID             string
	IamPermissions []string
	Fn             func(*session.Session) (map[string][]*ScrapeResult, error)
}

// Scrape invokes invokes the func that retrieves information from AWS
// and returns a response.
func (scraper *Scraper) Scrape(sess *session.Session) map[string][]*ScrapeResult {
	status := "success"
	start := time.Now()

	response, err := scraper.Fn(sess)
	if err != nil {
		log.Error(err)
		status = "error"
	}

	duration := time.Since(start)
	metrics.ScrapeDurationHistogram.WithLabelValues(scraper.ID, status).Observe(duration.Seconds())

	return response
}

// IsEnabled checks each scraper against a slice of strings. If a value
// in the slice matches the Id of a scraper then the scraper will be excluded
// from the returned Scraper slice.
func (scraper *Scraper) IsEnabled(disabledScrapers []string) bool {
	for _, s := range disabledScrapers {
		if s == scraper.ID {
			return false
		}
	}

	return true
}

// InitializeMetrics assigns the Prometheus.Desc pointer to the Metric property.
// We do this because once a Desc has been created, all the values are private
// and we can't render the metadata behind a metric easily. Therefore, a Scraper
// is created with all the metadata defined and then the Prometheus.Desc is created
// one the Scraper is registered into the ScrapeRegistry.
func (scraper *Scraper) InitializeMetrics() {
	for _, metric := range scraper.Metrics {
		metric.metric = prometheus.NewDesc(
			metric.PrefixMetricName(),
			metric.Description,
			metric.Labels,
			nil,
		)
	}
}
