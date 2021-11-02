# YAAE (Yet Another AWS Exporter)

> A Prometheus metrics exporter for AWS that fills in gaps CloudWatch doesn't cover

## About

This exporter is meant to expose metrics that aren't included in CloudWatch for Prometheus to scrape. When run in conjuction with a CloudWatch exporter (such as [YACE](https://github.com/nerdswords/yet-another-cloudwatch-exporter) or the [Prometheus Community CloudWatch Exporter](https://github.com/prometheus/cloudwatch_exporter)), this exporter provides increased visbility into your AWS ecosystem.

A full list of scrapers and metrics is found below.

## Image

```
cashapp/yet-another-aws-exporter:<TAG>
```

- [Images are published to DockerHub](https://hub.docker.com/repository/docker/cashapp/yet-another-aws-exporter/tags), replace `<TAG>` with whichever version you're targeting (i.e. `v1.0.0`)
- Binaries can be found in [Releases](https://github.com/cashappyet-another-aws-exporter/releases)

## Exported Metrics

- [A full list of exported metrics can be found here](./docs/metrics.md)
- Run `yaae scrapers list` to see all of the scrapers installed by default

## Features

- Self-documenting scraper/metric structure architecture
- All metrics/scrapers can be disabled
- IAM permissions can be generated off only selected scrapers
- [Internal metrics](./docs/metrics.md) for tracking error rates and scrape duration
- Structured JSON logging

## Configuration File

YAAE looks for a config file next to the binary `yaae` binary named `yaae.yaml`. To pass a custom path to a config file, use the `-c` flag:

```
yaae -c /path/to/yaae.yaml
```

A config file example [can be found in the `examples` directory](./examples).

### AWS IAM Permissions

This exporter should never require anything more than read-only permissions. If you're running every scraper available, you'll need the following permissions on the role/user that the exporter is running under:

```
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "yaae",
            "Effect": "Allow",
            "Action": [
                "eks:DescribeCluster",
                "eks:ListClusters",
                "ec2:DescribeVpcs",
                "ec2:DescribeSubnets",
                "iam:GetAccountSummary"
            ],
            "Resource": "*"
        }
    ]
}
```

Not using every exporter? Run `yaae scrapers list-iam-permissions -c <PATH TO CONFIG FILE>` to generate a new list of permissions.

# Contributing

If you're interested in contributing, we'd love the help! [See the contributing guidelines for initial details](./CONTRIBUTING.md).

