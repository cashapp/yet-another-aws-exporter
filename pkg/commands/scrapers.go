package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/cashapp/yet-another-aws-exporter/pkg/config"
	"github.com/cashapp/yet-another-aws-exporter/pkg/globals"
	"github.com/cashapp/yet-another-aws-exporter/pkg/scrapers"
	"github.com/olekukonko/tablewriter"
	log "github.com/sirupsen/logrus"
)

func wrapInTicks(s string) string {
	return fmt.Sprintf("`%s`", s)
}

// ScrapersCmd is a struct contianining all of the commands for working
// with scrapers configured in this project.
type ScrapersCmd struct {
	List ScrapersListCmd        `cmd help:"Output a markdown table of all scrapers"`
	Iam  ScrapersPermissionsCmd `cmd name:"list-iam-permissions" help:"Output the IAM permissions required for scrapers to run"`
}

// ScrapersListCmd outputs all the registered scrapers.
type ScrapersListCmd struct{}

func (cmd *ScrapersListCmd) Run(globals *globals.Globals) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Scraper ID", "Metric", "Description", "Labels"})
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoWrapText(false)

	for _, scraper := range scrapers.Registry.Scrapers {
		// Wrap labels in ticks
		labels := []string{}
		for _, l := range scraper.Labels {
			labels = append(labels, wrapInTicks(l))
		}

		// Append to table.
		table.Append([]string{
			wrapInTicks(scraper.ID),
			wrapInTicks(scraper.PrefixMetricName()),
			scraper.Description,
			strings.Join(labels, ","),
		})
	}
	table.Render() // Send Output

	return nil
}

type ScrapersPermissionsCmd struct{}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (cmd *ScrapersPermissionsCmd) Run(globals *globals.Globals) error {
	log.Debugf("Loading config at path: %s", globals.Config)
	config := &config.Config{ConfigPath: globals.Config}
	err := config.Load()
	if err != nil {
		return err
	}

	allPerms := []string{}
	for _, scraper := range scrapers.Registry.GetActiveScrapers(config) {
		for _, perm := range scraper.IamPermissions {
			if !contains(allPerms, perm) {
				allPerms = append(allPerms, perm)
			}
		}
	}

	log.Info("Ensure the following permissions are present on the role/user YAAE is running as:")
	log.SetOutput(os.Stdout)                   // Swap to stdout for perms listing
	fmt.Println(strings.Join(allPerms, ",\n")) //nolint
	return nil
}
