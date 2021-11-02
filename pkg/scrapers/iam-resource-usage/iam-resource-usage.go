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
		ID:          "iamResourceUsage",
		Name:        "iam_resource_usage_total",
		Description: "The number of IAM resources being used by resource type",
		Labels:      []string{"resource"},
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
func IamResourceUsageScrape(sess *session.Session) ([]*types.ScrapeResult, error) {
	client := iam.New(sess)
	scrapeResults := []*types.ScrapeResult{}

	summary, err := client.GetAccountSummary(&iam.GetAccountSummaryInput{})
	if err != nil {
		log.Error(err)
		return scrapeResults, err
	}

	// Iterate through all the resources above and grab their usage
	for _, resource := range resources {
		if val, ok := summary.SummaryMap[resource]; ok {
			scrapeResults = append(scrapeResults, &types.ScrapeResult{
				Labels: []string{strings.ToLower(resource)},
				Value:  float64(*val),
				Type:   prometheus.GaugeValue,
			})
		}
	}

	return scrapeResults, nil
}
