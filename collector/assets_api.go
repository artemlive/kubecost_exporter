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

type CloudAssets struct {
	Cloud []kubecost_api.CloudAssetOther
	Disk  []kubecost_api.CloudAssetDisk
}

func NewCloudAssets() *CloudAssets{
	return &CloudAssets{
		Cloud: []kubecost_api.CloudAssetOther{},
		Disk: []kubecost_api.CloudAssetDisk{},
	}
}
func (c *CloudAssets) GetDisks() *[]kubecost_api.CloudAssetDisk {
	return &c.Disk
}

func (c *CloudAssets) AddDisk(disk kubecost_api.CloudAssetDisk)  {
	c.Disk = append(c.Disk, disk)
}

const (
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
	level.Debug(logger).Log("msg", scrapeAssetsSubsystemName, "scraperParams", fmt.Sprintf("%+v, len(%d)", scraperParams, len(scraperParams)))
	apiClient := kubecost_api.NewApiClient(*apiBaseUrl, namespace)
	assets, err := apiClient.ListAssets(scraperParams)
	level.Debug(logger).Log("msg", fmt.Sprintf("%+v", assets))
	if err != nil {
		return err
	}

	return nil
}
