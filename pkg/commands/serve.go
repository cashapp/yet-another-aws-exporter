package commands

import (
	"fmt"
	"net/http"

	"github.com/cashapp/yet-another-aws-exporter/pkg/config"
	"github.com/cashapp/yet-another-aws-exporter/pkg/globals"
	"github.com/cashapp/yet-another-aws-exporter/pkg/metrics"
	"github.com/cashapp/yet-another-aws-exporter/pkg/scrapers"
	"github.com/cashapp/yet-another-aws-exporter/pkg/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type ServeCmd struct {
	Port string `type:"string" short:"p" help:"The port to run the exporter on" default:":9100"`
}

func (cmd *ServeCmd) Run(globals *globals.Globals) error {
	log.SetFormatter(&log.JSONFormatter{})
	log.Infof("Serving on port: %s", cmd.Port)
	log.Infof("Serving Prometheus metrics on /metrics")

	log.Debugf("Loading config at path: %s", globals.Config)
	config := &config.Config{ConfigPath: globals.Config}
	err := config.Load()
	if err != nil {
		return err
	}

	log.Debug("Registering internal metrics")
	metrics.InitMetrics()

	log.Debug("Registering collector metrics")
	prometheus.MustRegister(&types.Collector{
		Scrapers: scrapers.Registry.GetActiveScrapers(config),
	})

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
	<head><title>Yet Another AWS Exporter</title></head>
	<body>
	<h1>Yet Another AWS Exporter</h1>
	<p><a href="/metrics">Metrics</a></p>
	</body>
	</html>`))
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	return http.ListenAndServe(cmd.Port, nil)
}
