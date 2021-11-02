package types

import (
	"sync"

	"github.com/cashapp/yet-another-aws-exporter/pkg/sessions"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Collector implements the prometheus.Collector interface and wraps internal
// scrapers. During collection, scrapers are invoked and pass their results back
// so that they can be exported.
type Collector struct {
	Scrapers []*Scraper
}

// Describe implements the prometheus.Collector Describe method
// https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Collector
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	for _, scraper := range c.Scrapers {
		ch <- scraper.Metric
	}
}

// Collect implements the prometheus.Collector Collect method
// https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#Collector
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	sess, err := sessions.CreateAWSSession()
	if err != nil {
		log.Error(err)
		return
	}

	// Init WaitGroup. Without a WaitGroup the channel we write
	// results to will close before the goroutines finish
	var wg sync.WaitGroup
	wg.Add(len(c.Scrapers))

	// Iterate through all scrapers and invoke the scrape
	for _, scraper := range c.Scrapers {
		// Wrape the scrape invocation in a goroutine, but we need to pass
		// the scraper into the function explicitly to re-scope the variable
		// the goroutine accesses. If we don't do this, we can sometimes hit
		// a case where the scraper reports results twice and the collector panics
		go func(scraper *Scraper) {
			// Done call deferred until end of the scrape
			defer wg.Done()

			log.Debugf("Running scrape: %s", scraper.ID)
			scrapeResults := scraper.Scrape(sess)

			// Iterate through scrape results and send the metric
			for _, result := range scrapeResults {
				ch <- prometheus.MustNewConstMetric(scraper.Metric, result.Type, result.Value, result.Labels...)
			}
			log.Debugf("Scrape completed: %s", scraper.ID)
		}(scraper)
	}
	// Wait
	wg.Wait()
}
