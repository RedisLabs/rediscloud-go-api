package subscriptions

type CreateSubscription struct {
	Name                        string                `json:"name,omitempty"`
	DryRun                      bool                  `json:"dryRun"`
	PaymentMethodId             int                   `json:"paymentMethodId,omitempty"`
	MemoryStorage               string                `json:"memoryStorage,omitempty"`
	PersistentStorageEncryption bool                  `json:"persistentStorageEncryption,omitempty"`
	CloudProviders              []CreateCloudProvider `json:"cloudProviders,omitempty"`
	Databases                   []CreateDatabase      `json:"databases,omitempty"`
}

type CreateCloudProvider struct {
	Provider       string         `json:"provider,omitempty"`
	CloudAccountId int            `json:"cloudAccountId,omitempty"`
	Regions        []CreateRegion `json:"regions,omitempty"`
}

type CreateRegion struct {
	Region                     string            `json:"region,omitempty"`
	MultipleAvailabilityZones  bool              `json:"multipleAvailabilityZones,omitempty"`
	PreferredAvailabilityZones []string          `json:"preferredAvailabilityZones,omitempty"`
	Networking                 *CreateNetworking `json:"networking,omitempty"`
}

type CreateNetworking struct {
	DeploymentCIDR string `json:"deploymentCIDR,omitempty"`
	VPCId          string `json:"vpcId,omitempty"`
}

type CreateDatabase struct {
	Name                   string            `json:"name,omitempty"`
	Protocol               string            `json:"protocol,omitempty"`
	MemoryLimitInGb        float64           `json:"memoryLimitInGb,omitempty"`
	SupportOSSClusterApi   bool              `json:"supportOSSClusterApi,omitempty"`
	DataPersistence        string            `json:"dataPersistence,omitempty"`
	Replication            bool              `json:"replication,omitempty"`
	ThroughputMeasurement  *CreateThroughput `json:"throughputMeasurement,omitempty"`
	Modules                []CreateModules   `json:"modules,omitempty"`
	Quantity               int               `json:"quantity,omitempty"`
	AverageItemSizeInBytes int               `json:"averageItemSizeInBytes,omitempty"`
}

type CreateThroughput struct {
	By    string `json:"by"`
	Value int    `json:"value"`
}

type CreateModules struct {
	Name       string            `json:"name"`
	Parameters map[string]string `json:"parameters"`
}

type taskResponse struct {
	TaskId string `json:"taskId"`
}
