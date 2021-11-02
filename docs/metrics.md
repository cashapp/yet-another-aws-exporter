# Exported Metrics

Information around the metrics exported can be found below. 

- All YAAE internal metrics are prefixed with `yaae_`
- All AWS metrics are prefixed with `aws_`

## AWS Metrics

AWS metrics come from [individually configured Scrapers](./scrapers.md). A table of Scrapers and their corresponding timeseries are documented below. 

|      SCRAPER ID      |             METRIC             |                       DESCRIPTION                       |              LABELS               |
|----------------------|--------------------------------|---------------------------------------------------------|-----------------------------------|
| `eksInfo`            | `aws_eks_cluster_info`         | The current running EKS clusters in a region            | `cluster_name`,`version`,`status` |
| `vpcInfo`            | `aws_vpc_info`                 | The current running VPCs in a region                    | `vpc_id`,`state`                  |
| `subnetAvailableIps` | `aws_vpc_subnet_available_ips` | The number of IPs available in a subnet                 | `vpc_id`,`subnet_id`              |
| `iamResourceUsage`   | `aws_iam_resource_usage_total` | The number of IAM resources being used by resource type | `resource`                        |
| `iamQuotas`          | `aws_iam_resource_quota`       | The service quota cap for IAM resources                 | `resource`                        |


Not using all of the default scrapers?

## Internal Metrics

Initially there are only two internal metrics ([both of which can be found in the `metrics.go` package](https://github.com/cashapp/yet-another-aws-exporter/blob/main/pkg/metrics/metrics.go)):

- `yaae_scrape_duration_seconds`: a **histogram** of the duration it takes each scraper to run. Uses [Prometheus default histogram buckets](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus#DefBuckets).
- `yaae_api_call_errors_total`: a **counter** that tracks the rate of errors when calling AWS APIs. Useful for tracking down errors in IAM configuration issues.
