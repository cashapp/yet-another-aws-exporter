package main

import (
	"strings"

	"github.com/alecthomas/kong"
	"github.com/cashapp/yet-another-aws-exporter/pkg/commands"
	"github.com/cashapp/yet-another-aws-exporter/pkg/globals"
	log "github.com/sirupsen/logrus"
)

type CLI struct {
	globals.Globals

	Serve    commands.ServeCmd    `cmd help:"Start the exporter server"`
	Scrapers commands.ScrapersCmd `cmd help:"Commands to find retrieve information about scrapers"`
}

func main() {
	logLevels := []string{
		log.ErrorLevel.String(),
		log.WarnLevel.String(),
		log.InfoLevel.String(),
		log.DebugLevel.String(),
	}
	cli := CLI{
		Globals: globals.Globals{
			Version: globals.VersionFlag(""),
		},
	}
	ctx := kong.Parse(&cli,
		kong.Name("yet-another-aws-exporter"),
		kong.Description("A Prometheus metrics exporter that fills in gaps of AWS metrics"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
		kong.Vars{
			"version":     "0.0.2",
			"log_levels":  strings.Join(logLevels, ","),
			"config_file": "yaae.yaml",
		})

	level, err := log.ParseLevel(cli.LogLevel)
	ctx.FatalIfErrorf(err)
	// Set the global log level
	log.SetLevel(level)

	err = ctx.Run(&cli.Globals)
	ctx.FatalIfErrorf(err)
}
