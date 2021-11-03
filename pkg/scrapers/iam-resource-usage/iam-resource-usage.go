package iamresourceusage

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/cashapp/yet-another-aws-exporter/pkg/types"
)

// New returns an instance of the Scraper.
func New() *types.Scraper {
	return &types.Scraper{
		ID: "iamResourceUsage",
		Metrics: map[string]*types.Metric{
			"usage": &types.Metric{
				Name:        "iam_resource_usage_total",
				Description: "The number of IAM resources being used by resource type",
				Labels:      []string{"resource"},
			},
			"quotas": &types.Metric{
				Name:        "iam_resource_quota",
				Description: "The service quota cap for IAM resources",
				Labels:      []string{"resource"},
			},
		},
		IamPermissions: []string{
			"iam:GetAccountSummary",
		},
		Fn: IamResourceUsageScrape,
	}
}

var (
	resources = []string{
		"Roles",
		"Users",
		"Groups",
		"Policies",
	}
)

// IamRoleCountScrape queries the AWS IAM API for all of the roles in an account using the
// Account Summary endpoint.
// https://docs.aws.amazon.com/IAM/latest/APIReference/API_GetAccountSummary.html
func IamResourceUsageScrape(sess *session.Session) (map[string][]*types.ScrapeResult, error) {
	client := iam.New(sess)
	scrapeResults := map[string][]*types.ScrapeResult{}
	usage := []*types.ScrapeResult{}
	quotas := []*types.ScrapeResult{}

	summary, err := client.GetAccountSummary(&iam.GetAccountSummaryInput{})
	if err != nil {
		log.Error(err)
		return scrapeResults, err
	}

	// Iterate through all the resources above and grab their usage
	// and quotas and append to the proper list of results
	for _, resource := range resources {
		// Capture usage info
		if val, ok := summary.SummaryMap[resource]; ok {
			usage = append(usage, &types.ScrapeResult{
				Labels: []string{strings.ToLower(resource)},
				Value:  float64(*val),
				Type:   prometheus.GaugeValue,
			})
		}
		// Capture quota info
		if val, ok := summary.SummaryMap[resource+"Quota"]; ok {
			quotas = append(quotas, &types.ScrapeResult{
				Labels: []string{strings.ToLower(resource)},
				Value:  float64(*val),
				Type:   prometheus.GaugeValue,
			})
		}
	}

	// Add to the return struct
	scrapeResults["usage"] = usage
	scrapeResults["quotas"] = quotas

	return scrapeResults, nil
}
