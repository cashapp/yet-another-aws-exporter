package eksinfo

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/cashapp/yet-another-aws-exporter/pkg/metrics"
	"github.com/cashapp/yet-another-aws-exporter/pkg/types"
)

// New returns an instance of the Scraper.
func New() *types.Scraper {
	return &types.Scraper{
		ID: "eksInfo",
		Metrics: map[string]*types.Metric{
			"info": &types.Metric{
				Name:        "eks_cluster_info",
				Description: "The current running EKS clusters in a region",
				Labels:      []string{"cluster_name", "version", "status"},
			},
		},
		IamPermissions: []string{
			"eks:DescribeCluster",
			"eks:ListClusters",
		},
		Fn: EksClusterInfo,
	}
}

// EksClusterInfo scrapes the basic information about EKS clusters in a region. The value
// of the metric will always be 1, but the labels will include information about cluster name,
// Kubernetes versiona and whether or not the cluster is active.
func EksClusterInfo(sess *session.Session) (map[string][]*types.ScrapeResult, error) {
	scrapeResults := map[string][]*types.ScrapeResult{}
	info := []*types.ScrapeResult{}
	client := eks.New(sess)

	clusters, err := client.ListClusters(&eks.ListClustersInput{})
	if err != nil {
		log.Error(err)
		metrics.APICallErrorsTotal.WithLabelValues("eks", "ListClusters").Inc()
		return scrapeResults, err
	}

	for _, cluster := range clusters.Clusters {
		log.Debugf("Retrieving EKS clusters information for: %s", *cluster)
		c, err := client.DescribeCluster(&eks.DescribeClusterInput{
			Name: cluster,
		})
		if err != nil {
			log.Error(err)
			metrics.APICallErrorsTotal.WithLabelValues("eks", "DescribeCluster").Inc()
			return scrapeResults, err
		}

		metricVal := 1
		if *c.Cluster.Status != "ACTIVE" {
			metricVal = 0
		}

		info = append(info, &types.ScrapeResult{
			Labels: []string{*c.Cluster.Name, *c.Cluster.Version, strings.ToLower(*c.Cluster.Status)},
			Value:  float64(metricVal),
			Type:   prometheus.GaugeValue,
		})
	}

	scrapeResults["info"] = info

	return scrapeResults, nil
}
