package subnetavailableips

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cashapp/yet-another-aws-exporter/pkg/metrics"
	"github.com/cashapp/yet-another-aws-exporter/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// New returns an instance of the Scraper.
func New() *types.Scraper {
	return &types.Scraper{
		ID:          "subnetAvailableIps",
		Name:        "vpc_subnet_available_ips",
		Description: "The number of IPs available in a subnet",
		Labels:      []string{"vpc_id", "subnet_id"},
		IamPermissions: []string{
			"ec2:DescribeVpcs",
			"ec2:DescribeSubnets",
		},
		Fn: SubnetAvailableIpsScrape,
	}
}

// SubnetAvailableIpsScrape scrapes all subnets in a region to return the number of available
// IP addresses in the subnet.
func SubnetAvailableIpsScrape(sess *session.Session) ([]*types.ScrapeResult, error) {
	scrapeResults := []*types.ScrapeResult{}
	client := ec2.New(sess)

	subs, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{})
	if err != nil {
		log.Error(err)
		metrics.APICallErrorsTotal.WithLabelValues("vpc", "DescribeVpcs").Inc()
		return scrapeResults, err
	}

	for _, s := range subs.Subnets {
		scrapeResults = append(scrapeResults, &types.ScrapeResult{
			Labels: []string{*s.VpcId, *s.SubnetId},
			Value:  float64(*s.AvailableIpAddressCount),
			Type:   prometheus.GaugeValue,
		})
	}

	return scrapeResults, nil
}
