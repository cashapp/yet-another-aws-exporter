package vpcinfo

import (
	"net"

	"github.com/apparentlymart/go-cidr/cidr"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/cashapp/yet-another-aws-exporter/pkg/metrics"
	"github.com/cashapp/yet-another-aws-exporter/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Metric aliases for consistent naming
const infoMetric = "info"
const subnetAvailableIpsMetric = "subnetAvailableIps"
const subnetTotalIpsMetric = "subnetTotalIps"

// New returns an instance of the Scraper.
func New() *types.Scraper {
	return &types.Scraper{
		ID: "vpcInfo",
		Metrics: map[string]*types.Metric{
			infoMetric: &types.Metric{
				Name:        "vpc_info",
				Description: "The current running VPCs in a region",
				Labels:      []string{"vpc_id", "state"},
			},
			subnetAvailableIpsMetric: &types.Metric{
				Name:        "vpc_subnet_ips_available",
				Description: "The number of IPs available in a subnet",
				Labels:      []string{"vpc_id", "subnet_id"},
			},
			subnetTotalIpsMetric: &types.Metric{
				Name:        "vpc_subnet_ips_capacity",
				Description: "The total number of available IP addresses in a subnet CIDR",
				Labels:      []string{"vpc_id", "subnet_id"},
			},
		},
		IamPermissions: []string{
			"ec2:DescribeVpcs",
			"ec2:DescribeSubnets",
		},
		Fn: VpcInfoScrape,
	}
}

// VpcInfoScrape scrapes the EC2 API for information about each VPC in a region
// and publishes a metric that will aways be 1.
func VpcInfoScrape(sess *session.Session) (map[string][]*types.ScrapeResult, error) {
	scrapeResults := map[string][]*types.ScrapeResult{}
	info := []*types.ScrapeResult{}
	availableIps := []*types.ScrapeResult{}
	totalIps := []*types.ScrapeResult{}
	client := ec2.New(sess)

	vpcs, err := client.DescribeVpcs(&ec2.DescribeVpcsInput{})
	if err != nil {
		metrics.APICallErrorsTotal.WithLabelValues("vpc", "DescribeVpcs").Inc()
		return scrapeResults, err
	}

	// Iterate through all VPCs
	for _, v := range vpcs.Vpcs {
		// Append to the info slice
		info = append(info, &types.ScrapeResult{
			Labels: []string{*v.VpcId, *v.State},
			Value:  1.0,
			Type:   prometheus.CounterValue,
		})

		// Query for all the subnets in the same VPC
		subs, err := client.DescribeSubnets(&ec2.DescribeSubnetsInput{
			Filters: []*ec2.Filter{
				{
					Name: aws.String("vpc-id"),
					Values: []*string{
						v.VpcId, // Filter by VPC ID
					},
				},
			},
		})
		if err != nil {
			log.Error(err)
			metrics.APICallErrorsTotal.WithLabelValues("vpc", "DescribeSubnets").Inc()
			return scrapeResults, err
		}

		// Iterate through all subnets and append results
		for _, s := range subs.Subnets {
			// Return the number of available IPs
			availableIps = append(availableIps, &types.ScrapeResult{
				Labels: []string{*s.VpcId, *s.SubnetId},
				Value:  float64(*s.AvailableIpAddressCount),
				Type:   prometheus.GaugeValue,
			})

			// Find the total number of IPs in the subnet CIDR range
			totalIps = append(totalIps, &types.ScrapeResult{
				Labels: []string{*s.VpcId, *s.SubnetId},
				Value:  float64(getTotalIPCount(*s.CidrBlock)),
				Type:   prometheus.GaugeValue,
			})
		}
	}

	// Append results
	scrapeResults[infoMetric] = info
	scrapeResults[subnetAvailableIpsMetric] = availableIps
	scrapeResults[subnetTotalIpsMetric] = totalIps

	return scrapeResults, nil
}

// getTotalIPCount returns the max address count in a CIDR.
func getTotalIPCount(cidrString string) uint64 {
	_, ipnet, _ := net.ParseCIDR(cidrString)
	return cidr.AddressCount(ipnet)
}
