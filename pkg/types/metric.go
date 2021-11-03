package types

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// Metric is an internal wrapper around Prometheus metrics. A Scraper
// can produce multiple timeseries, and this struct is meant to provide
// a wrapper around each time series so we can treat them all uniformly.
type Metric struct {
	Name        string
	Description string
	Labels      []string
	metric      *prometheus.Desc
}

// PrefixMetricName adds `aws_` to the beginning of every metric name.
func (m *Metric) PrefixMetricName() string {
	return fmt.Sprintf("aws_%s", m.Name)
}
