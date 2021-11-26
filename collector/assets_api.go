package collector

import (
	"context"
	"fmt"
	"github.com/artemlive/kubecost_exporter/kubecost_api"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"net/url"
)

const (
	// Scrape query.
	assetsEndpoint = "model/assets"

	// Subsystem.
	scrapeAssetsSubsystemName = "scrape_assets"
)

type ScrapeAssets struct{}

func (ScrapeAssets) Name() string {
	return scrapeAssetsSubsystemName
}

func (ScrapeAssets) Help() string {
	return "Scrapes the information about Assets API"
}

func (ScrapeAssets) Scrape(ctx context.Context, apiBaseUrl **url.URL, scraperParams []string, ch chan<- prometheus.Metric, logger log.Logger) error {
	level.Info(logger).Log("msg", scrapeAssetsSubsystemName, "scraperParams", fmt.Sprintf("%+v, len(%d)", scraperParams, len(scraperParams)))
	apiClient := kubecost_api.NewApiClient(*apiBaseUrl, namespace)
	assets, err := apiClient.ListAssets(scraperParams)
	if err != nil {
		return err
	}

	level.Debug(logger).Log("msg", "%+v", assets)
	return nil
}
