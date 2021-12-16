package collector

import (
	"context"
	"fmt"
	"github.com/artemlive/kubecost_exporter/kubecost_api"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"net/url"
	"time"
)

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

// TODO: move to helpers
func TruncateDate(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}
func (s ScrapeAssets) Scrape(ctx context.Context, apiBaseUrl **url.URL, scraperParams []string, ch chan<- prometheus.Metric, logger log.Logger, skipTLSVerify bool, offset int64) error {
	//2021-12-14T00:00:00Z,2021-12-15T00:00:00Z
	RFC3339local := "2006-01-02T15:04:05Z"
	now := time.Now()
	dateFrom := now.AddDate(0, 0, int(-offset))
	dateTo := now.AddDate(0, 0, int(-offset+1))
	scraperParams = append(scraperParams, fmt.Sprintf("window=%s,%s", TruncateDate(dateFrom).Format(RFC3339local), TruncateDate(dateTo).Format(RFC3339local)))
	level.Debug(logger).Log("msg", scrapeAssetsSubsystemName, "scraperParams", fmt.Sprintf("%+v, len(%d)", scraperParams, len(scraperParams)))
	apiClient := kubecost_api.NewApiClient(*apiBaseUrl, namespace, skipTLSVerify)
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
	err = s.generateDisksMetrics(cloudAssetsMapper.GetDisks(), cloudAssetsMapper, ch, logger)
	if err != nil {
		return err
	}

	// Generate metrics for Clouds
	err = s.generateCloudMetrics(cloudAssetsMapper.GetClouds(), cloudAssetsMapper, ch, logger)
	if err != nil {
		return err
	}

	// Generate metrics for Nodes
	err = s.generateNodeMetrics(cloudAssetsMapper.GetNodes(), cloudAssetsMapper, ch, logger)
	if err != nil {
		return err
	}

	// Generate metrics for LoadBalancers
	err = s.generateLoadBalancerMetrics(cloudAssetsMapper.GetLoadBalancers(), cloudAssetsMapper, ch, logger)
	if err != nil {
		return err
	}
	return nil
}

// I know that all functions repeat same code, I'll refactor it a bit later
// TODO: move all common code to a separate function
func (ScrapeAssets) generateDisksMetrics(disks *[]kubecost_api.CloudAssetDisk, assetsMapper *CloudAssets, ch chan<- prometheus.Metric, logger log.Logger) error {
	if len(*disks) > 0 {
		for _, disk := range *disks {
			// maybe this is not the best idea to cast asset -> interface -> asset
			// TODO: refactor this to use common interface for all cloud assets
			labelNames, labelValues, err := assetsMapper.GetDefaultLabelsFromAssets(disk)
			if err != nil {
				return err
			}
			diskDesc := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, promDescSubsystem, "total"),
				"Assets total cost from Kubecost Assets API",
				labelNames, nil,
			)

			ch <- prometheus.MustNewConstMetric(
				diskDesc, prometheus.GaugeValue, disk.TotalCost, labelValues...,
			)
		}
	}
	return nil
}



func (ScrapeAssets) generateCloudMetrics(clouds *[]kubecost_api.CloudAssetCloud, assetsMapper *CloudAssets, ch chan<- prometheus.Metric, logger log.Logger) error {
	if len(*clouds) > 0 {
		for _, cloud := range *clouds {
			// maybe this is not the best idea to cast asset -> interface -> asset
			// TODO: refactor this to use common interface for all cloud assets
			labelNames, labelValues, err := assetsMapper.GetDefaultLabelsFromAssets(cloud)
			if err != nil {
				return err
			}
			diskDesc := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, promDescSubsystem, "total"),
				"Assets total cost from Kubecost Assets API",
				labelNames, nil,
			)
			ch <- prometheus.MustNewConstMetric(
				diskDesc, prometheus.GaugeValue, cloud.TotalCost, labelValues...,
			)
		}
	}
	return nil
}

func (ScrapeAssets) generateNodeMetrics(nodes *[]kubecost_api.CloudAssetNode, assetsMapper *CloudAssets, ch chan<- prometheus.Metric, logger log.Logger) error {
	if len(*nodes) > 0 {
		for _, node := range *nodes {
			// maybe this is not the best idea to cast asset -> interface -> asset
			// TODO: refactor this to use common interface for all cloud assets
			labelNames, labelValues, err := assetsMapper.GetDefaultLabelsFromAssets(node)
			if err != nil {
				return err
			}
			diskDesc := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, promDescSubsystem, "total"),
				"Assets total cost from Kubecost Assets API",
				labelNames, nil,
			)
			ch <- prometheus.MustNewConstMetric(
				diskDesc, prometheus.GaugeValue, node.TotalCost, labelValues...,
			)
		}
	}
	return nil
}

func (ScrapeAssets) generateLoadBalancerMetrics(nodes *[]kubecost_api.CloudAssetLoadBalancer, assetsMapper *CloudAssets, ch chan<- prometheus.Metric, logger log.Logger) error {
	if len(*nodes) > 0 {
		for _, node := range *nodes {
			// maybe this is not the best idea to cast asset -> interface -> asset
			// TODO: refactor this to use common interface for all cloud assets
			labelNames, labelValues, err := assetsMapper.GetDefaultLabelsFromAssets(node)
			if err != nil {
				return err
			}
			diskDesc := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, promDescSubsystem, "total"),
				"Assets total cost from Kubecost Assets API",
				labelNames, nil,
			)
			ch <- prometheus.MustNewConstMetric(
				diskDesc, prometheus.GaugeValue, node.TotalCost, labelValues...,
			)
		}
	}
	return nil
}

func (ScrapeAssets) generateClusterManagementMetrics(cm *[]kubecost_api.CloudAssetClusterManagement, assetsMapper *CloudAssets, ch chan<- prometheus.Metric, logger log.Logger) error {
	if len(*cm) > 0 {
		for _, c := range *cm {
			// maybe this is not the best idea to cast asset -> interface -> asset
			// TODO: refactor this to use common interface for all cloud assets
			labelNames, labelValues, err := assetsMapper.GetDefaultLabelsFromAssets(c)
			if err != nil {
				return err
			}
			diskDesc := prometheus.NewDesc(
				prometheus.BuildFQName(namespace, promDescSubsystem, "total"),
				"Assets total cost from Kubecost Assets API",
				labelNames, nil,
			)
			ch <- prometheus.MustNewConstMetric(
				diskDesc, prometheus.GaugeValue, c.TotalCost, labelValues...,
			)
		}
	}
	return nil
}
