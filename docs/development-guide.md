# Development Guide

This exporter relies on AWS user or role credentials in an environment. When beginning local development, ensure you can assume a role or have access tokens for a user with read-only IAM permissions for the resources you're looking to scrape.

## AWS Configuration

The exporter will use the environment variables present in your shell, with sessions initialized with [Shared Configuration Fields](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/sessions.html#shared-configuration-fields). 

## Running Locally

Right now, the easiest development flow is to run the exporter on your host machine with the following command:

```
go run ./cmd/yaae serve
```

This will start the exporter's server on port `9100`. To test your changes, send a request to `localhost:9100/metrics` to see the metrics that are scraped (`curl localhost:9100/metrics`). If you're working on a single scraper, it might be easier to disable all the other scrapers. 

### Containerized Development

Right now there's no workflow for containerized development, but that would be a great contribution!
