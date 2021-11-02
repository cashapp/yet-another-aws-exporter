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

func errScrape(sess *session.Session) ([]*ScrapeResult, error) {
	var err = errors.New("Scrape error")
	return []*ScrapeResult{}, err
}

func successScrape(sess *session.Session) ([]*ScrapeResult, error) {
	return []*ScrapeResult{successResult}, nil
}

func TestScraper_Scrape(t *testing.T) {
	tests := []struct {
		name string
		s    *Scraper
		want []*ScrapeResult
	}{
		{
			name: "success",
			s: &Scraper{
				Name:        "test",
				Description: "Description",
				Fn:          successScrape,
			},
			want: []*ScrapeResult{successResult},
		},
		{
			name: "error",
			s: &Scraper{
				Name:        "test",
				Description: "Description",
				Fn:          errScrape,
			},
			want: []*ScrapeResult{},
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
				ID:          "foo",
				Name:        "foo",
				Description: "Foo description",
				Fn:          successScrape,
			},
			disabledScrapers: []string{"bar"},
			want:             true,
		},
		{
			name: "returns false for scraper in the disabled list",
			scraper: &Scraper{
				ID:          "foo",
				Name:        "foo",
				Description: "Foo description",
				Fn:          successScrape,
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
				ID:          "foo",
				Name:        "foo",
				Description: "Foo description",
				Fn:          successScrape,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.scraper.InitializeMetric()
			assert.NotNil(t, tt.scraper.Metric)
		})
	}
}

func TestScraper_PrefixMetricName(t *testing.T) {
	tests := []struct {
		name    string
		scraper *Scraper
		want    string
	}{
		{
			name: "Adds `aws_` prefix to scraper names",
			scraper: &Scraper{
				Name: "foo",
			},
			want: "aws_foo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.scraper.PrefixMetricName(); got != tt.want {
				t.Errorf("Scraper.PrefixMetricName() = %v, want %v", got, tt.want)
			}
		})
	}
}
