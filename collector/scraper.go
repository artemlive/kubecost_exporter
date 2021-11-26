package collector

import (
	"context"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus"
	"net/url"
)

type Scraper interface {
	// Name of the Scraper. Should be unique.
	Name() string

	// Help describes the role of the Scraper.
	// Example: "Collect from KubeCost Assets API"
	Help() string

	// Scrape collects data from the KubeCost Assets API and sends it over channel as prometheus metric.
	Scrape(ctx context.Context, apiBaseUrl **url.URL, scraperParams []string, ch chan<- prometheus.Metric, logger log.Logger) error
}
