package types

import (
	"os"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

var (
	goodScraper = &Scraper{
		ID:          "good",
		Name:        "foo",
		Description: "Foo description",
		Fn:          successScrape,
	}
)

func restoreEnvFn(envVars map[string]string) {
	for k := range envVars {
		os.Unsetenv(k)
	}
}

func TestCollector_Describe(t *testing.T) {
	goodScraper.InitializeMetric()

	tests := []struct {
		name string
		c    *Collector
	}{
		{
			name: "sends the scraper's metric to the provided channel",
			c: &Collector{
				Scrapers: []*Scraper{goodScraper},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := make(chan *prometheus.Desc, len(tt.c.Scrapers))
			tt.c.Describe(metrics)
			assert.Equal(t, len(metrics), len(tt.c.Scrapers))
			close(metrics)
		})
	}
}

func TestCollector_Collect(t *testing.T) {
	tests := []struct {
		name    string
		c       *Collector
		envVars map[string]string
		want    int
	}{
		{
			name: "passes a recorded metric from a scraper",
			c: &Collector{
				Scrapers: []*Scraper{goodScraper},
			},
			envVars: map[string]string{},
			want:    1,
		},
		{
			name: "collects no metric if session creation errors",
			c: &Collector{
				Scrapers: []*Scraper{goodScraper},
			},
			envVars: map[string]string{
				"AWS_STS_REGIONAL_ENDPOINTS": "fake",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			metrics := make(chan prometheus.Metric, len(tt.c.Scrapers))
			defer close(metrics)
			defer restoreEnvFn(tt.envVars)

			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			tt.c.Collect(metrics)
			assert.Equal(t, len(metrics), tt.want)
		})
	}
}
