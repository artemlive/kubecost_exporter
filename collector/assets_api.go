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
	"strings"
)

type CloudAssets struct {
	logger log.Logger
	Cloud  []kubecost_api.CloudAssetCloud
	Disk   []kubecost_api.CloudAssetDisk
}

type CloudAssetsPropertiesToPrometheus map[string]string

func NewCloudAssets(logger log.Logger) *CloudAssets {
	return &CloudAssets{
		logger: logger,
		Cloud:  []kubecost_api.CloudAssetCloud{},
		Disk:   []kubecost_api.CloudAssetDisk{},
	}
}
func (c *CloudAssets) GetDisks() *[]kubecost_api.CloudAssetDisk {
	return &c.Disk
}

func (c *CloudAssets) AddDisk(disk kubecost_api.CloudAssetDisk) {
	level.Debug(c.logger).Log("AddDisk", fmt.Sprintf("%+v", disk))
	c.Disk = append(c.Disk, disk)
}

func (c *CloudAssets) AddCloud(cloud kubecost_api.CloudAssetCloud) {
	level.Debug(c.logger).Log("AddCloud", fmt.Sprintf("%+v", cloud))
	c.Cloud = append(c.Cloud, cloud)
}

const (
	// Subsystem for logging.
	scrapeAssetsSubsystemName = "scrape_assets"
	// Subsystem for exporter metrics
	promDescSubsystem = "cost"
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
	// to avoid duplication
	// if don't use accumulate, it would duplicate resources usage for multiple time windows
	scraperParams = append(scraperParams, "accumulate=true")
	assets, err := apiClient.ListAssets(scraperParams)
	if err != nil {
		return err
	}
	cloudAssetsMapper := NewCloudAssets(logger)
	err = cloudAssetsMapper.MapAssets(assets)

	// Generate metrics for Disks
	// TODO: move that to another function
	disks := *cloudAssetsMapper.GetDisks()

	if len(disks) > 0 {
		for _, disk := range disks {
			// maybe this is not the best idea to cast asset -> interface -> asset
			// TODO: refactor this to use common interface for all cloud assets
			labelNames, labelValues, err := cloudAssetsMapper.GetDefaultLabelsFromAssets(disk)
			if err != nil {
				return err
			}
			diskDesc := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, promDescSubsystem, "total"),
				"Disk total cost from Kubecost Assets API",
				labelNames, nil,
			)
			ch <- prometheus.MustNewConstMetric(
				diskDesc, prometheus.GaugeValue, disk.TotalCost, labelValues...,
			)
		}
	}

	return err
}

//type AssetProperties struct {
//	Category   string `json:"category,omitempty"`
//	Provider   string `json:"provider,omitempty"`
//	Account    string `json:"account,omitempty"`
//	Project    string `json:"project,omitempty"`
//	Service    string `json:"service,omitempty"`
//	Cluster    string `json:"cluster,omitempty"`
//	Name       string `json:"name,omitempty"`
//	ProviderID string `json:"providerID,omitempty"`
//}

// Sets and return the default labels set for each assets
func (c *CloudAssets) GetDefaultLabelsFromAssets(asset interface{}) ([]string, []string, error) {
	switch asset.(type) {
	case kubecost_api.CloudAssetDisk:
		disk, ok := asset.(kubecost_api.CloudAssetDisk)
		if !ok {
			fmt.Errorf("couldn't cast interface to CloudAssetDisk: %+v", asset)
		}
		labels, labelsVals, err := c.getLabelsFromAsset(disk.Labels)
		if err != nil {
			return []string{}, []string{}, err
		}
		// we have to create values list according to the defaultDiskLabels
		propertiesLabels, propertiesLabelsVals := c.getEnabledProperties(disk.Properties)
		diskLabels := append(propertiesLabels, "type")
		diskLabelsValues := append(propertiesLabelsVals, disk.Type)
		// concat array of properties labels/values with actual labels/values
		outLabels := append(diskLabels, labels...)
		outValues := append(diskLabelsValues, labelsVals...)
		return outLabels, outValues, err
	}
	return []string{}, []string{}, nil
}

// mapping default properties from assets api to the corresponding prometheus labels
// not all of these fields are set for assets, so we have to check all of them, to understand which labels we have to export
func (c *CloudAssets) getEnabledProperties(properties *kubecost_api.AssetProperties) ([]string, []string) {
	// I didn't want to make a lot of if statements
	// my other attempts to rewrite this code had failed and I didn't want to waste time
	// there were a lot of reflect code, which is not readable and efficient
	// TODO: refactor this after POC testing
	var enabledPropsLabels []string
	var enabledPropsValues []string
	if len(properties.Category) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_category")
		enabledPropsValues = append(enabledPropsValues, properties.Category)
	}
	if len(properties.Name) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_name")
		enabledPropsValues = append(enabledPropsValues, properties.Name)
	}
	if len(properties.Cluster) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_cluster")
		enabledPropsValues = append(enabledPropsValues, properties.Cluster)
	}
	if len(properties.Service) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_service")
		enabledPropsValues = append(enabledPropsValues, properties.Service)
	}
	if len(properties.Account) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_account")
		enabledPropsValues = append(enabledPropsValues, properties.Account)
	}
	if len(properties.Project) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_project")
		enabledPropsValues = append(enabledPropsValues, properties.Project)
	}
	if len(properties.Provider) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_provider")
		enabledPropsValues = append(enabledPropsValues, properties.Provider)
	}
	if len(properties.ProviderID) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_provider_id")
		enabledPropsValues = append(enabledPropsValues, properties.ProviderID)
	}
	return enabledPropsLabels, enabledPropsValues
}

// Generates two arrays, first of them is slice of labels, second one is slice of corresponding values
// used for prometheus metric
func (c *CloudAssets) getLabelsFromAsset(labels kubecost_api.AssetLabels) ([]string, []string, error) {
	var outLabels []string
	var outValues []string
	for k, v := range labels {
		outLabels = append(outLabels, strings.ReplaceAll(k, "-", "_"))
		val, ok := v.(string)
		if !ok {
			return nil, nil, fmt.Errorf("couldn't process label value to string: %+v", val)
		}
		outValues = append(outValues, strings.ReplaceAll(val, "-", "_"))
	}
	return outLabels, outValues, nil
}

// the function that maps different resources types eg Cloud/Disk/Node to according Cloud Assets instance
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
		c.AddDiskFromMap(asset)
	}
	return nil
}

// this is the abstraction above AddDisk
// this function maps the map[string]interface from Api response to concrete CloudAssetDisk
// and adds it to the disks list
func (c *CloudAssets) AddDiskFromMap(asset map[string]interface{}) error {
	disk := kubecost_api.CloudAssetDisk{}
	// convert json to struct
	// Convert map to json string
	jsonStr, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonStr, &disk)
	c.AddDisk(disk)
	return err
}

// this is the abstraction above AddCloud
// this function maps the map[string]interface from Api response to concrete CloudAssetCloud
// and adds it to the Cloud assets list
func (c *CloudAssets) AddCloudFromMap(asset map[string]interface{}) error {
	cloud := kubecost_api.CloudAssetCloud{}
	// convert json to struct
	// Convert map to json string
	jsonStr, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonStr, &cloud)
	c.AddCloud(cloud)
	return err
}
