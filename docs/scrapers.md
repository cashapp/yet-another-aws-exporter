# Scrapers

This exporter is a simple [Prometheus exporter](https://prometheus.io/docs/instrumenting/writing_exporters/) but centralizes around the idea of a "Scraper" being the source of truth for timeseries. The docs below contain information about Scraper basics, [but looking at a real scraper might be more helpful](https://github.com/cashapp/yet-another-aws-exporter/blob/main/pkg/scrapers/vpc-info/vpc.go).

## What Is A Scraper?

A Scraper is a struct which contains metadata about a Prometheus metric and a single function that will be invoked to determine that metric's value.The properties of a Scraper do the following:

- `ID`: the identifier of the scraper that is the value used for disabling scrapers in the YAAE config file
- `Metrics`: a map of key/value pairs, where the key is a string and the value is of type `Metric`
  - `Name`: the name of the Prometheus metrics; `aws_` will be prepended to this value automatically
  - `Description`: maps to the Prometheus metric description
  - `Labels`: maps to Prometheus metric labels
- `IamPermissions`: a list of IAM permissions the individual scraper requires
- `Fn`: the scrape function that the collector will invoke when `/metrics` endpoint is hit

### The Scrape Function

The Scrape Function is the method that will be invoked each time the `/metrics` endpoint receives a request. Here is a simple sample:

```go
// Assume the following Scraper is initialized
var (
	myScraper = &types.Scraper{
		ID: "vpcInfo",
		Metrics: map[string]*types.Metric{
			"myMetricAlias": &types.Metric{
				Name:        "my_cool_metric_total",
				Description: "An example",
				Labels:      []string{"label", "another_label"},
			},
		IamPermissions: []string{
      ...IAM permissions based on scrape function needs
		},
		Fn: ExampleScrape,
	}
}

// The scrape function receives an AWS Session with which a client for an AWS service can be initialized
func ExampleScrape(sess *session.Session) (map[string][]*ScrapeResult, error){
  ...Do some logic...

  return map[string][]*ScrapeResult{
    "myMetricAlias": []*ScrapeResult{
      Type: prometheus.CounterValue,
      Value: float64(1),
      Labels: []string{"labelValue1", "labelValue2"},
    }
  }, nil
}
```

Inside your function, you can do anything you want, but there are only a few guidelines to follow:

1. The `map` returned by the function **must have the same keys as the `Metrics` map in the Scraper**
2. Your function should execute quickly. Default scrape timeout for Prometheus is 10s, your scraper should be much faster than that
3. You should perform read-only operations, this exporter should never modify any AWS resources

## Naming Metrics

When creating a metric, be sure to follow the [Prometheus Metric & Label Naming](https://prometheus.io/docs/practices/naming/) guidelines!

And keep in mind, all metrics will be prefixed with `aws_`!
