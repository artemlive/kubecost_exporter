package collector

import (
	"context"
	"fmt"
	"github.com/artemlive/kubecost_exporter/kubecost_api"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"net/url"
	"strings"
	"time"
)

const (
	// Subsystem for logging.
	scrapeAllocationSubsystemName = "scrape_allocation"
	// Subsystem for exporter metrics
	promDesc = "cost_cluster_allocation"
)

type ScrapeAllocation struct{}

func (ScrapeAllocation) Name() string {
	return scrapeAllocationSubsystemName
}

func (ScrapeAllocation) Help() string {
	return "Scrapes the information about Cost Allocation API"
}

func (s ScrapeAllocation) Scrape(ctx context.Context, apiBaseUrl **url.URL, scraperParams []string, ch chan<- prometheus.Metric, logger log.Logger, skipTLSVerify bool, offset int64) error {
	//2021-12-14T00:00:00Z,2021-12-15T00:00:00Z
	RFC3339local := "2006-01-02T15:04:05Z"
	now := time.Now()
	dateFrom := now.AddDate(0, 0, int(-offset))
	dateTo := now.AddDate(0, 0, int(-offset+1))
	scraperParams = append(scraperParams, fmt.Sprintf("window=%s,%s", TruncateDate(dateFrom).Format(RFC3339local), TruncateDate(dateTo).Format(RFC3339local)))
	scraperParams = append(scraperParams, "accumulate=true")
	level.Debug(logger).Log("msg", scrapeAllocationSubsystemName, "scraperParams", fmt.Sprintf("%+v, len(%d)", scraperParams, len(scraperParams)))
	apiClient := kubecost_api.NewApiClient(*apiBaseUrl, namespace, skipTLSVerify)
	costs, err := apiClient.GetAllocation(scraperParams)
	if err != nil {
		return err
	}
	if len(costs.Data) == 0 {
		return fmt.Errorf("empty allocations")
	}

	// weird response, that has map in a first element of an array
	for _,cost := range costs.Data[0]{
		s.generateMetric(cost, ch, logger)
	}
	return nil
}

// This function maps labels with their values for prometheus metric construction
func (s ScrapeAllocation) getDefaultLabels(allocation kubecost_api.Allocation) ([]string, []string, error){
	// if there is an idle resources allocation
	// we mark namespace as __idle__
	// so we can calculcate idle resources for cluster via that label
	if allocation.Name == "__idle__" {
		return []string{"property_cluster", "property_namespace"}, []string{allocation.Properties.Cluster, allocation.Name}, nil
	}
	// TODO: refactor this after POC testing
	var enabledPropsLabels []string
	var enabledPropsValues []string
	if len(allocation.Properties.Namespace) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_namespace")
		enabledPropsValues = append(enabledPropsValues, allocation.Properties.Namespace)
	}
	if len(allocation.Properties.Node) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_node")
		enabledPropsValues = append(enabledPropsValues, allocation.Properties.Node)
	}
	if len(allocation.Properties.Cluster) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_cluster")
		enabledPropsValues = append(enabledPropsValues, allocation.Properties.Cluster)
	}
	if len(allocation.Properties.ProviderID) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_provider_id")
		enabledPropsValues = append(enabledPropsValues, allocation.Properties.ProviderID)
	}
	if len(allocation.Properties.Container) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_container")
		enabledPropsValues = append(enabledPropsValues, allocation.Properties.Container)
	}
	if len(allocation.Properties.Controller) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_controller")
		enabledPropsValues = append(enabledPropsValues, allocation.Properties.Controller)
	}
	if len(allocation.Properties.Pod) > 0 {
		enabledPropsLabels = append(enabledPropsLabels, "property_pod")
		enabledPropsValues = append(enabledPropsValues, allocation.Properties.Pod)
	}
	s.getDefaultAllocationLabels(&enabledPropsLabels, &enabledPropsValues, allocation.Properties.Labels)
	return enabledPropsLabels,enabledPropsValues, nil
}

func (s ScrapeAllocation) getDefaultAllocationLabels(labelNames *[]string, labelValues *[]string, labels kubecost_api.AllocationLabels) error{
	for name, value := range labels {
		*labelNames = append(*labelNames, strings.ReplaceAll(name, "-", "_"))
		*labelValues = append(*labelValues, strings.ReplaceAll(value, "-", "_"))
	}
	return nil
}

func (s ScrapeAllocation) generateMetric(allocation kubecost_api.Allocation, ch chan<- prometheus.Metric, logger log.Logger) error {
	labelNames, labelValues, err := s.getDefaultLabels(allocation)
	if err != nil {
		return err
	}
	diskDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, promDesc, "total"),
		"k8s total cost from Kubecost Assets API",
		labelNames, nil,
	)
	ch <- prometheus.MustNewConstMetric(
		diskDesc, prometheus.GaugeValue, allocation.TotalCost, labelValues...,
	)
	return nil
}