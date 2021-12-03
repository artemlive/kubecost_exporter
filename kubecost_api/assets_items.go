package kubecost_api

import "time"

// TODO: we may contribute to the cost-model repository to refactor their assets definitions
// I couldn't use their assets definitions, because it wasn't possible to include only assets, there were some prometheus things,
// that were conflicting with our current packages

type Window struct {
	Start *time.Time
	End   *time.Time
}

type AssetLabels map[string]interface{}

type Breakdown struct {
	Idle   float64 `json:"idle"`
	Other  float64 `json:"other"`
	System float64 `json:"system"`
	User   float64 `json:"user"`
}

type CloudAssetLoadBalancer struct {
	Type       string           `json:"type"`
	Properties *AssetProperties `json:"properties"`
	Labels     AssetLabels      `json:"labels"`
	Start      time.Time        `json:"start"`
	End        time.Time        `json:"end"`
	Window     Window           `json:"window"`
	Minutes    int              `json:"minutes"`
	Adjustment float64          `json:"adjustment"`
	TotalCost  float64          `json:"totalCost"`
}

type CloudAssetClusterManagement struct {
	Type       string           `json:"type"`
	Labels     AssetLabels      `json:"labels"`
	Properties *AssetProperties `json:"properties"`
	Window     Window           `json:"window"`
	Minutes    int              `json:"minutes"`
	Adjustment float64          `json:"adjustment"`
	TotalCost  float64          `json:"totalCost"`
}
type AssetProperties struct {
	Category   string `json:"category,omitempty"`
	Provider   string `json:"provider,omitempty"`
	Account    string `json:"account,omitempty"`
	Project    string `json:"project,omitempty"`
	Service    string `json:"service,omitempty"`
	Cluster    string `json:"cluster,omitempty"`
	Name       string `json:"name,omitempty"`
	ProviderID string `json:"providerID,omitempty"`
}

// Cloud assets api with type "Disk"
type CloudAssetDisk struct {
	Type       string           `json:"type"`
	Properties *AssetProperties `json:"properties"`
	Labels     AssetLabels      `json:"labels"`
	Window     Window           `json:"window"`
	Start      time.Time        `json:"start"`
	End        time.Time        `json:"end"`
	Minutes    int              `json:"minutes"`
	ByteHours  float64          `json:"byteHours"`
	Bytes      int64            `json:"bytes"`
	Breakdown  *Breakdown       `json:"breakdown"`
	Adjustment int              `json:"adjustment"`
	TotalCost  float64          `json:"totalCost"`
}

type CloudAssetCloud struct {
	Type       string           `json:"type"`
	Properties *AssetProperties `json:"properties"`
	Labels     AssetLabels      `json:"labels"`
	Window     Window           `json:"window"`
	Start      time.Time        `json:"start"`
	End        time.Time        `json:"end"`
	Minutes    float64          `json:"minutes"`
	Adjustment float64          `json:"adjustment"`
	Credit     float64          `json:"credit"`
	TotalCost  float64          `json:"totalCost"`
}

type CloudAssetNode struct {
	Type         string           `json:"type"`
	Properties   *AssetProperties `json:"properties"`
	Labels       AssetLabels      `json:"labels"`
	Start        time.Time        `json:"start"`
	End          time.Time        `json:"end"`
	Window       Window           `json:"window"`
	NodeType     string           `json:"nodeType"`
	CPUCoreHours float64          `json:"cpuCoreHours"`
	RAMByteHours float64          `json:"ramByteHours"`
	GPUHours     float64          `json:"GPUHours"`
	CPUBreakdown *Breakdown       `json:"cpuBreakdown"`
	RAMBreakdown *Breakdown       `json:"ramBreakdown"`
	CPUCost      float64          `json:"cpuCost"`
	GPUCost      float64          `json:"gpuCost"`
	GPUCount     float64          `json:"gpuCount"`
	RAMCost      float64          `json:"ramCost"`
	Discount     float64          `json:"discount"`
	Preemptible  float64          `json:"preemptible"`
	Adjustment   float64          `json:"adjustment"`
	Credit       float64          `json:"credit"`
	TotalCost    float64          `json:"totalCost"`
}
