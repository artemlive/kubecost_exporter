package kubecost_api

import "time"

type ReservedInstanceData struct {
	ReservedCPU int64   `json:"reservedCPU"`
	ReservedRAM int64   `json:"reservedRAM"`
	CPUCost     float64 `json:"CPUHourlyCost"`
	RAMCost     float64 `json:"RAMHourlyCost"`
}

type PricingType string

type Node struct {
	Cost             string                `json:"hourlyCost"`
	VCPU             string                `json:"CPU"`
	VCPUCost         string                `json:"CPUHourlyCost"`
	RAM              string                `json:"RAM"`
	RAMBytes         string                `json:"RAMBytes"`
	RAMCost          string                `json:"RAMGBHourlyCost"`
	Storage          string                `json:"storage"`
	StorageCost      string                `json:"storageHourlyCost"`
	UsesBaseCPUPrice bool                  `json:"usesDefaultPrice"`
	BaseCPUPrice     string                `json:"baseCPUPrice"` // Used to compute an implicit RAM GB/Hr price when RAM pricing is not provided.
	BaseRAMPrice     string                `json:"baseRAMPrice"` // Used to compute an implicit RAM GB/Hr price when RAM pricing is not provided.
	BaseGPUPrice     string                `json:"baseGPUPrice"`
	UsageType        string                `json:"usageType"`
	GPU              string                `json:"gpu"` // GPU represents the number of GPU on the instance
	GPUName          string                `json:"gpuName"`
	GPUCost          string                `json:"gpuCost"`
	InstanceType     string                `json:"instanceType,omitempty"`
	Region           string                `json:"region,omitempty"`
	Reserved         *ReservedInstanceData `json:"reserved,omitempty"`
	ProviderID       string                `json:"providerID,omitempty"`
	PricingType      PricingType           `json:"pricingType,omitempty"`
}

type Vector struct {
	Timestamp float64 `json:"timestamp"`
	Value     float64 `json:"value"`
}

type Allocation struct {
	Name                       string                `json:"name"`
	Properties                 *AllocationProperties `json:"properties,omitempty"`
	Window                     Window                `json:"window"`
	Start                      time.Time             `json:"start"`
	End                        time.Time             `json:"end"`
	CPUCoreHours               float64               `json:"cpuCoreHours"`
	CPUCoreRequestAverage      float64               `json:"cpuCoreRequestAverage"`
	CPUCoreUsageAverage        float64               `json:"cpuCoreUsageAverage"`
	CPUCost                    float64               `json:"cpuCost"`
	CPUCostAdjustment          float64               `json:"cpuCostAdjustment"`
	GPUHours                   float64               `json:"gpuHours"`
	GPUCost                    float64               `json:"gpuCost"`
	GPUCostAdjustment          float64               `json:"gpuCostAdjustment"`
	NetworkTransferBytes       float64               `json:"networkTransferBytes"`
	NetworkReceiveBytes        float64               `json:"networkReceiveBytes"`
	NetworkCost                float64               `json:"networkCost"`
	NetworkCostAdjustment      float64               `json:"networkCostAdjustment"`
	LoadBalancerCost           float64               `json:"loadBalancerCost"`
	LoadBalancerCostAdjustment float64               `json:"loadBalancerCostAdjustment"`
	PVs                        PVAllocations         `json:"-"`
	PVCostAdjustment           float64               `json:"pvCostAdjustment"`
	RAMByteHours               float64               `json:"ramByteHours"`
	RAMBytesRequestAverage     float64               `json:"ramByteRequestAverage"`
	RAMBytesUsageAverage       float64               `json:"ramByteUsageAverage"`
	RAMCost                    float64               `json:"ramCost"`
	RAMCostAdjustment          float64               `json:"ramCostAdjustment"`
	SharedCost                 float64               `json:"sharedCost"`
	ExternalCost               float64               `json:"externalCost"`
	// RawAllocationOnly is a pointer so if it is not present it will be
	// marshalled as null rather than as an object with Go default values.
	RawAllocationOnly *RawAllocationOnlyData `json:"rawAllocationOnly"`
	TotalCost 				   float64               `json:"totalCost"`
}

type RawAllocationOnlyData struct {
	CPUCoreUsageMax  float64 `json:"cpuCoreUsageMax"`
	RAMBytesUsageMax float64 `json:"ramByteUsageMax"`
}

type PVAllocations map[PVKey]*PVAllocation

// PVKey for identifying Disk type assets
type PVKey struct {
	Cluster string `json:"cluster"`
	Name    string `json:"name"`
}

type PVAllocation struct {
	ByteHours float64 `json:"byteHours"`
	Cost      float64 `json:"cost"`
}


// AllocationProperties describes a set of Kubernetes objects.
type AllocationProperties struct {
	Cluster        string                `json:"cluster,omitempty"`
	Node           string                `json:"node,omitempty"`
	Container      string                `json:"container,omitempty"`
	Controller     string                `json:"controller,omitempty"`
	ControllerKind string                `json:"controllerKind,omitempty"`
	Namespace      string                `json:"namespace,omitempty"`
	Pod            string                `json:"pod,omitempty"`
	Services       []string              `json:"services,omitempty"`
	ProviderID     string                `json:"providerID,omitempty"`
	Labels         AllocationLabels      `json:"labels,omitempty"`
	Annotations    AllocationAnnotations `json:"annotations,omitempty"`
}

// AllocationLabels is a schema-free mapping of key/value pairs that can be
// attributed to an Allocation
type AllocationLabels map[string]string

// AllocationAnnotations is a schema-free mapping of key/value pairs that can be
// attributed to an Allocation
type AllocationAnnotations map[string]string


type CostDataResponse struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data []map[string]Allocation `json:"data"`
}

type PersistentVolumeClaimData struct {
	Class        string   `json:"class"`
	Claim        string   `json:"claim"`
	Namespace    string   `json:"namespace"`
	ClusterID    string   `json:"clusterId"`
	TimesClaimed int      `json:"timesClaimed"`
	VolumeName   string   `json:"volumeName"`
	Volume       PV       `json:"persistentVolume"`
	Values       []Vector `json:"values"`
}

// PV is the interface by which the provider and cost model communicate PV prices.
// The provider will best-effort try to fill out this struct.
type PV struct {
	Cost       string            `json:"hourlyCost"`
	CostPerIO  string            `json:"costPerIOOperation"`
	Class      string            `json:"storageClass"`
	Size       string            `json:"size"`
	Region     string            `json:"region"`
	ProviderID string            `json:"providerID,omitempty"`
	Parameters map[string]string `json:"parameters"`
}
