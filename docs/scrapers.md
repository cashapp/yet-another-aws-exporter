# Scrapers

This exporter is a simple [Prometheus exporter](https://prometheus.io/docs/instrumenting/writing_exporters/) but centralizes around the idea of a "Scraper" being the source of truth for an individual timeseries. The docs below contain information about Scraper basics, [but looking at a real scraper might be more helpful](https://github.com/cashapp/yet-another-aws-exporter/blob/main/pkg/scrapers/eks-info/awseks.go).

## What Is A Scraper?

A Scraper is a struct which contains metadata about a Prometheus metric and a single function that will be invoked to determine that metric's value.The properties of a Scraper do the following:

- `ID`: the identifier of the scraper that is the value used for disabling scrapers in the YAAE config file
- `Name`: the name of the Prometheus metrics; `aws_` will be prepended to this value automatically 
- `Description`: maps to the Prometheus metric description
- `Labels`: maps to Prometheus metric labels
- `IamPermissions`: a list of IAM permissions the individual scraper requires
- `Fn`: the scrape function that the collector will invoke when `/metrics` endpoint is hit 

### The Scrape Function

The Scrape Function is the method that will be invoked each time the `/metrics` endpoint receives a request. This function has a simple signature: 

- Receives an AWS Session (`\*session.Session`)
- [Returns a `ScrapeResult` struct](https://github.com/cashapp/yet-another-aws-exporter/blob/7899bc86586b2bda032e895f756cd248d8dd4bb4/pkg/types/scraper.go#L13-L19)

Inside a Scrape Function you can do anything you want with the session, but the operation should run as quickly as possible so that the default Prometheus scrape duration is not exceeded.

## Naming Metrics

When creating a metric, be sure to follow the [Prometheus Metric & Label Naming](https://prometheus.io/docs/practices/naming/) guidelines!

And keep in mind, all metrics will be prefixed with `aws_`!
