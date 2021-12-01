package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/artemlive/kubecost_exporter/kubecost_api"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"net/url"
)

type CloudAssets struct {
	logger log.Logger
	Cloud []kubecost_api.CloudAssetOther
	Disk  []kubecost_api.CloudAssetDisk
}

func NewCloudAssets(logger log.Logger) *CloudAssets{
	return &CloudAssets{
		logger: logger,
		Cloud: []kubecost_api.CloudAssetOther{},
		Disk: []kubecost_api.CloudAssetDisk{},
	}
}
func (c *CloudAssets) GetDisks() *[]kubecost_api.CloudAssetDisk {
	return &c.Disk
}

func (c *CloudAssets) AddDisk(disk kubecost_api.CloudAssetDisk)  {
	level.Debug(c.logger).Log("AddDisk", fmt.Sprintf("%+v", disk))
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
	//level.Debug(logger).Log("msg", fmt.Sprintf("%+v", assets))
	if err != nil {
		return err
	}
	cloudAssetsMapper := NewCloudAssets(logger)
	err = cloudAssetsMapper.MapAssets(assets)
	return err
}

// function that maps different resources types eg Cloud/Disk/Node to a CloudAssets instance
func (c *CloudAssets) MapAssets(value interface{}) error {
	switch value.(type) {
	case []interface{}:
		for _, v := range value.([]interface{}) {
			c.MapAssets(v)
		}
	case map[string]interface{}:
		for k, v := range value.(map[string]interface{}) {
			//if key is data, we have to go deeper
			if k == "data" {
				level.Debug(c.logger).Log("msg", "structure seems to be correct, keyword 'data' found")
				c.MapAssets(v)
			} else {
				// according to current structure we get the asset itself
				// the format is: asset_uniq_key => asset_config interface{}
				level.Debug(c.logger).Log("msg", fmt.Sprintf("k = %s", k))
				concreteType, ok := v.(map[string]interface{})
				if ok {
					if err := c.addAccordingType(concreteType); err != nil {
						return err
					}
				}

			}
		}
	default:
		return fmt.Errorf("unknown map type")
	}
	return nil
}

func (c *CloudAssets) addAccordingType(asset map[string]interface{}) error {
	// ignore fields that doesn't match string => interface
	valType, ok := asset["type"]
	if !ok {
		return fmt.Errorf("asset %+v, doesn't have \"type\" field", asset)
	}
	switch valType {
		case "Disk":
			disk := kubecost_api.CloudAssetDisk{}
			// convert json to struct
			// Convert map to json string
			jsonStr, err := json.Marshal(asset)
			if err != nil {
				return err
			}
			json.Unmarshal(jsonStr, &disk)
			c.AddDisk(disk)
	}

	return nil
}