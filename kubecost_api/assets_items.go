package kubecost_api

import "time"

// Cloud assets api with type "Disk"
type CloudAssetDisk struct {
	Type       string `json:"type"`
	Properties struct {
		Category string `json:"category"`
		Service  string `json:"service"`
		Cluster  string `json:"cluster"`
		Name     string `json:"name"`
	} `json:"properties"`
	Labels map[string]interface {
	} `json:"labels"`
	Window struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"window"`
	Start     time.Time `json:"start"`
	End       time.Time `json:"end"`
	Minutes   int       `json:"minutes"`
	ByteHours float64   `json:"byteHours"`
	Bytes     int64     `json:"bytes"`
	Breakdown struct {
		Idle   int `json:"idle"`
		Other  int `json:"other"`
		System int `json:"system"`
		User   int `json:"user"`
	} `json:"breakdown"`
	Adjustment int     `json:"adjustment"`
	TotalCost  float64 `json:"totalCost"`
}

type CloudAssetOther struct {
	Type       string `json:"type"`
	Properties struct {
		Category   string `json:"category"`
		Provider   string `json:"provider"`
		Account    string `json:"account"`
		Project    string `json:"project"`
		Service    string `json:"service"`
		Name       string `json:"name"`
		ProviderID string `json:"providerID"`
	} `json:"properties"`
	Labels map[string]interface {} `json:"labels"`
	Window struct {
		Start time.Time `json:"start"`
		End   time.Time `json:"end"`
	} `json:"window"`
	Start      time.Time `json:"start"`
	End        time.Time `json:"end"`
	Minutes    float64   `json:"minutes"`
	Adjustment float64   `json:"adjustment"`
	Credit     float64   `json:"credit"`
	TotalCost  float64   `json:"totalCost"`
}
