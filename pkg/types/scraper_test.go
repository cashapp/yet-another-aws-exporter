package types

import (
	"errors"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/cashapp/yet-another-aws-exporter/pkg/sessions"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

var (
	successResult = &ScrapeResult{
		Labels: []string{},
		Value:  1.0,
		Type:   prometheus.GaugeValue,
	}
)

func errScrape(sess *session.Session) (map[string][]*ScrapeResult, error) {
	var err = errors.New("Scrape error")
	return map[string][]*ScrapeResult{}, err
}

func successScrape(sess *session.Session) (map[string][]*ScrapeResult, error) {
	return map[string][]*ScrapeResult{
		"example": []*ScrapeResult{successResult},
	}, nil
}

func TestScraper_Scrape(t *testing.T) {
	tests := []struct {
		name string
		s    *Scraper
		want map[string][]*ScrapeResult
	}{
		{
			name: "success",
			s: &Scraper{
				ID: "success",
				Fn: successScrape,
			},
			want: map[string][]*ScrapeResult{
				"example": []*ScrapeResult{successResult},
			},
		},
		{
			name: "error",
			s: &Scraper{
				ID: "error",
				Fn: errScrape,
			},
			want: map[string][]*ScrapeResult{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sess, _ := sessions.CreateAWSSession()
			if got := tt.s.Scrape(sess); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Scraper.Scrape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScraper_IsEnabled(t *testing.T) {
	tests := []struct {
		name             string
		scraper          *Scraper
		disabledScrapers []string
		want             bool
	}{
		{
			name: "returns true for scraper not in the disabled list",
			scraper: &Scraper{
				ID: "foo",
				Fn: successScrape,
			},
			disabledScrapers: []string{"bar"},
			want:             true,
		},
		{
			name: "returns false for scraper in the disabled list",
			scraper: &Scraper{
				ID: "foo",
				Fn: successScrape,
			},
			disabledScrapers: []string{"foo"},
			want:             false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scraper.IsEnabled(tt.disabledScrapers); got != tt.want {
				t.Errorf("Scraper.IsEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScraper_InitializeMetric(t *testing.T) {
	tests := []struct {
		name    string
		scraper *Scraper
	}{
		{
			name: "creates a Prometheus metric for a configured scraper",
			scraper: &Scraper{
				ID: "foo",
				Metrics: map[string]*Metric{
					"example": &Metric{
						Name:        "foo_metric",
						Description: "Foo description",
						Labels:      nil,
					},
				},
				Fn: successScrape,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.scraper.InitializeMetrics()

			for _, m := range tt.scraper.Metrics {
				assert.NotNil(t, m.metric)
			}
		})
	}
}
