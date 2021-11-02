package scrapers

import (
	"github.com/cashapp/yet-another-aws-exporter/pkg/config"
	eksinfo "github.com/cashapp/yet-another-aws-exporter/pkg/scrapers/eks-info"
	iamquotas "github.com/cashapp/yet-another-aws-exporter/pkg/scrapers/iam-quotas"
	iamresourceusage "github.com/cashapp/yet-another-aws-exporter/pkg/scrapers/iam-resource-usage"
	subnetavailableips "github.com/cashapp/yet-another-aws-exporter/pkg/scrapers/subnet-available-ips"
	vpcinfo "github.com/cashapp/yet-another-aws-exporter/pkg/scrapers/vpc-info"
	"github.com/cashapp/yet-another-aws-exporter/pkg/types"
	log "github.com/sirupsen/logrus"
)

// On module load, create a registry of all scrapers. When adding a new
// scraper, it will need to be added here.
func init() {
	Registry.Add(eksinfo.New())
	Registry.Add(vpcinfo.New())
	Registry.Add(subnetavailableips.New())
	Registry.Add(iamresourceusage.New())
	Registry.Add(iamquotas.New())
}

var (
	// Registry is where all scrapers will be registered to.
	Registry ScraperRegistry
)

// ScraperRegistry is a struct which which contains a collection of Scrapers.
// This collection can then be filtered/iterated over to retrieve only the
// scrapers that should be active for the exporter.
type ScraperRegistry struct {
	Scrapers []*types.Scraper
}

// Add appends a Scraper to the Scrapers slice on the registry.
func (sr *ScraperRegistry) Add(s *types.Scraper) {
	// Initialize the Prometheus metric pointer
	s.InitializeMetric()
	// Append initialized scraper to the slice of all scrapers
	sr.Scrapers = append(sr.Scrapers, s)
}

// GetActiveScrapers iterates through the Scrapers slice and returns only those
// which are not disabled by the config passed to the exporter at initialization.
func (sr *ScraperRegistry) GetActiveScrapers(config *config.Config) []*types.Scraper {
	activeScrapers := []*types.Scraper{}

	for _, s := range sr.Scrapers {
		if s.IsEnabled(config.DisabledScrapers) {
			log.Debugf("Scraper enabled: %s", s.ID)
			activeScrapers = append(activeScrapers, s)
		} else {
			log.Debugf("Scraper disabled: %s", s.ID)
		}
	}

	return activeScrapers
}
