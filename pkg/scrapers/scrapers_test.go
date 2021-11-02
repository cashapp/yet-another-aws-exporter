package scrapers

import (
	"os"
	"path"
	"testing"

	"github.com/cashapp/yet-another-aws-exporter/pkg/config"
	"github.com/cashapp/yet-another-aws-exporter/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestScraperRegistry_Add(t *testing.T) {
	tests := []struct {
		name    string
		scraper *types.Scraper
	}{
		{
			name: "adds a scraper to the global registry and calls InitializeMetric",
			scraper: &types.Scraper{
				ID:          "foo",
				Name:        "foo",
				Description: "Foo description",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sr := &ScraperRegistry{}
			sr.Add(tt.scraper)

			assert.Equal(t, len(sr.Scrapers), 1)
			// Assert that the initialize function was called during adding
			assert.NotNil(t, sr.Scrapers[0].Metric)
		})
	}
}

func TestScraperRegistry_GetActiveScrapers(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
	}{
		{
			name:       "returns the scrapers that are currently active",
			configPath: "../../examples/disable-example.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := os.Getwd()
			if err != nil {
				t.Fatal(err)
			}

			config := &config.Config{
				ConfigPath: path.Join(p, tt.configPath),
			}
			err = config.Load()
			if err != nil {
				t.Fatal(err)
			}

			// get active count from global modules
			active := Registry.GetActiveScrapers(config)
			// we want the full registry length minus the disabled length so that
			// updating the examples doesn't cause someone to constantly update a static
			// value in this test file
			assert.Equal(t, len(active), len(Registry.Scrapers)-len(config.DisabledScrapers))
		})
	}
}
