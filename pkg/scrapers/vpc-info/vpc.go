package vpcinfo

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cashapp/yet-another-aws-exporter/pkg/metrics"
	"github.com/cashapp/yet-another-aws-exporter/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
)

// New returns an instance of the Scraper.
func New() *types.Scraper {
	return &types.Scraper{
		ID:          "vpcInfo",
		Name:        "vpc_info",
		Description: "The current running VPCs in a region",
		Labels:      []string{"vpc_id", "state"},
		IamPermissions: []string{
			"ec2:DescribeVpcs",
		},
		Fn: VpcInfoScrape,
	}
}

// VpcInfoScrape scrapes the EC2 API for information about each VPC in a region
// and publishes a metric that will aways be 1.
func VpcInfoScrape(sess *session.Session) ([]*types.ScrapeResult, error) {
	scrapeResults := []*types.ScrapeResult{}
	client := ec2.New(sess)

	vpcs, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		metrics.APICallErrorsTotal.WithLabelValues("vpc", "DescribeVpcs").Inc()
		return scrapeResults, err
	}

	for _, v := range vpcs.Vpcs {
		scrapeResults = append(scrapeResults, &types.ScrapeResult{
			Labels: []string{*v.VpcId, *v.State},
			Value:  1.0,
			Type:   prometheus.CounterValue,
		})
	}

	return scrapeResults, nil
}
