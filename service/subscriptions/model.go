package subscriptions

import (
	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type CreateSubscription struct {
	Name                        *string                `json:"name,omitempty"`
	DryRun                      *bool                  `json:"dryRun,omitempty"`
	PaymentMethodID             *int                   `json:"paymentMethodId,omitempty"`
	MemoryStorage               *string                `json:"memoryStorage,omitempty"`
	PersistentStorageEncryption *bool                  `json:"persistentStorageEncryption,omitempty"`
	CloudProviders              []*CreateCloudProvider `json:"cloudProviders,omitempty"`
	Databases                   []*CreateDatabase      `json:"databases,omitempty"`
}

func (o CreateSubscription) String() string {
	return internal.ToString(o)
}

type CreateCloudProvider struct {
	Provider       *string         `json:"provider,omitempty"`
	CloudAccountID *int            `json:"cloudAccountId,omitempty"`
	Regions        []*CreateRegion `json:"regions,omitempty"`
}

func (o CreateCloudProvider) String() string {
	return internal.ToString(o)
}

type CreateRegion struct {
	Region                     *string           `json:"region,omitempty"`
	MultipleAvailabilityZones  *bool             `json:"multipleAvailabilityZones,omitempty"`
	PreferredAvailabilityZones []*string         `json:"preferredAvailabilityZones,omitempty"`
	Networking                 *CreateNetworking `json:"networking,omitempty"`
}

func (o CreateRegion) String() string {
	return internal.ToString(o)
}

type CreateNetworking struct {
	DeploymentCIDR *string `json:"deploymentCIDR,omitempty"`
	VPCId          *string `json:"vpcId,omitempty"`
}

func (o CreateNetworking) String() string {
	return internal.ToString(o)
}

type CreateDatabase struct {
	Name                   *string           `json:"name,omitempty"`
	Protocol               *string           `json:"protocol,omitempty"`
	MemoryLimitInGB        *float64          `json:"memoryLimitInGb,omitempty"`
	SupportOSSClusterAPI   *bool             `json:"supportOSSClusterApi,omitempty"`
	DataPersistence        *string           `json:"dataPersistence,omitempty"`
	Replication            *bool             `json:"replication,omitempty"`
	ThroughputMeasurement  *CreateThroughput `json:"throughputMeasurement,omitempty"`
	Modules                []*CreateModules  `json:"modules,omitempty"`
	Quantity               *int              `json:"quantity,omitempty"`
	AverageItemSizeInBytes *int              `json:"averageItemSizeInBytes,omitempty"`
}

func (o CreateDatabase) String() string {
	return internal.ToString(o)
}

type CreateThroughput struct {
	By    *string `json:"by,omitempty"`
	Value *int    `json:"value,omitempty"`
}

func (o CreateThroughput) String() string {
	return internal.ToString(o)
}

type CreateModules struct {
	Name       *string            `json:"name,omitempty"`
	Parameters map[string]*string `json:"parameters,omitempty"`
}

func (o CreateModules) String() string {
	return internal.ToString(o)
}

type taskResponse struct {
	ID *string `json:"taskId,omitempty"`
}

func (o taskResponse) String() string {
	return internal.ToString(o)
}
