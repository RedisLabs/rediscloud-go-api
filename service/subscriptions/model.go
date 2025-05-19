package subscriptions

import (
	"fmt"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type CreateSubscription struct {
	Name            *string                `json:"name,omitempty"`
	DeploymentType  *string                `json:"deploymentType,omitempty"`
	DryRun          *bool                  `json:"dryRun,omitempty"`
	PaymentMethodID *int                   `json:"paymentMethodId,omitempty"`
	PaymentMethod   *string                `json:"paymentMethod,omitempty"`
	MemoryStorage   *string                `json:"memoryStorage,omitempty"`
	CloudProviders  []*CreateCloudProvider `json:"cloudProviders,omitempty"`
	Databases       []*CreateDatabase      `json:"databases,omitempty"`
	RedisVersion    *string                `json:"redisVersion,omitempty"`
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
	Name                       *string                  `json:"name,omitempty"`
	Protocol                   *string                  `json:"protocol,omitempty"`
	MemoryLimitInGB            *float64                 `json:"memoryLimitInGb,omitempty"`
	DatasetSizeInGB            *float64                 `json:"datasetSizeInGb,omitempty"`
	SupportOSSClusterAPI       *bool                    `json:"supportOSSClusterApi,omitempty"`
	DataPersistence            *string                  `json:"dataPersistence,omitempty"`
	Replication                *bool                    `json:"replication,omitempty"`
	ThroughputMeasurement      *CreateThroughput        `json:"throughputMeasurement,omitempty"`
	LocalThroughputMeasurement []*CreateLocalThroughput `json:"localThroughputMeasurement,omitempty"`
	Modules                    []*CreateModules         `json:"modules,omitempty"`
	Quantity                   *int                     `json:"quantity,omitempty"`
	AverageItemSizeInBytes     *int                     `json:"averageItemSizeInBytes,omitempty"`
	QueryPerformanceFactor     *string                  `json:"queryPerformanceFactor,omitempty"`
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

type CreateLocalThroughput struct {
	Region                   *string `json:"region,omitempty"`
	WriteOperationsPerSecond *int    `json:"writeOperationsPerSecond"`
	ReadOperationsPerSecond  *int    `json:"readOperationsPerSecond"`
}

func (o CreateLocalThroughput) String() string {
	return internal.ToString(o)
}

type CreateModules struct {
	Name *string `json:"name,omitempty"`
}

func (o CreateModules) String() string {
	return internal.ToString(o)
}

type UpdateSubscription struct {
	Name            *string `json:"name,omitempty"`
	PaymentMethodID *int    `json:"paymentMethodId,omitempty"`
}

func (o UpdateSubscription) String() string {
	return internal.ToString(o)
}

type Subscription struct {
	ID                *int           `json:"id,omitempty"`
	Name              *string        `json:"name,omitempty"`
	Status            *string        `json:"status,omitempty"`
	DeploymentType    *string        `json:"deploymentType,omitempty"`
	PaymentMethod     *string        `json:"paymentMethodType,omitempty"`
	PaymentMethodID   *int           `json:"paymentMethodId,omitempty"`
	MemoryStorage     *string        `json:"memoryStorage,omitempty"`
	StorageEncryption *bool          `json:"storageEncryption,omitempty"`
	NumberOfDatabases *int           `json:"numberOfDatabases,omitempty"`
	CloudDetails      []*CloudDetail `json:"cloudDetails,omitempty"`
}

func (o Subscription) String() string {
	return internal.ToString(o)
}

type CloudDetail struct {
	Provider       *string   `json:"provider,omitempty"`
	CloudAccountID *int      `json:"cloudAccountId,omitempty"`
	TotalSizeInGB  *float64  `json:"totalSizeInGb,omitempty"`
	Regions        []*Region `json:"regions,omitempty"`
}

func (o CloudDetail) String() string {
	return internal.ToString(o)
}

type Region struct {
	Region                     *string       `json:"region,omitempty"`
	Networking                 []*Networking `json:"networking,omitempty"`
	PreferredAvailabilityZones []*string     `json:"preferredAvailabilityZones,omitempty"`
	MultipleAvailabilityZones  *bool         `json:"multipleAvailabilityZones,omitempty"`
}

func (o Region) String() string {
	return internal.ToString(o)
}

type Networking struct {
	DeploymentCIDR *string `json:"deploymentCIDR,omitempty"`
	VPCId          *string `json:"vpcId,omitempty"`
	SubnetID       *string `json:"subnetId,omitempty"`
}

func (o Networking) String() string {
	return internal.ToString(o)
}

type CIDRAllowlist struct {
	CIDRIPs          []*string   `json:"cidr_ips,omitempty"`
	SecurityGroupIDs []*string   `json:"security_group_ids,omitempty"`
	Errors           interface{} `json:"errors,omitempty"`
}

func (o CIDRAllowlist) String() string {
	return internal.ToString(o)
}

type UpdateCIDRAllowlist struct {
	CIDRIPs          []*string `json:"cidrIps,omitempty"`
	SecurityGroupIDs []*string `json:"securityGroupIds,omitempty"`
}

func (o UpdateCIDRAllowlist) String() string {
	return internal.ToString(o)
}

type CreateVPCPeering struct {
	Region         *string   `json:"region,omitempty"`
	AWSAccountID   *string   `json:"awsAccountId,omitempty"`
	VPCId          *string   `json:"vpcId,omitempty"`
	VPCCidr        *string   `json:"vpcCidr,omitempty"`
	VPCCidrs       []*string `json:"vpcCidrs,omitempty"`
	Provider       *string   `json:"provider,omitempty"`
	VPCProjectUID  *string   `json:"vpcProjectUid,omitempty"`
	VPCNetworkName *string   `json:"vpcNetworkName,omitempty"`
}

func (o CreateVPCPeering) String() string {
	return internal.ToString(o)
}

type CreateActiveActiveVPCPeering struct {
	SourceRegion      *string   `json:"sourceRegion,omitempty"`
	DestinationRegion *string   `json:"destinationRegion,omitempty"`
	AWSAccountID      *string   `json:"awsAccountId,omitempty"`
	VPCId             *string   `json:"vpcId,omitempty"`
	VPCCidr           *string   `json:"vpcCidr,omitempty"`
	VPCCidrs          []*string `json:"vpcCidrs,omitempty"`
	Provider          *string   `json:"provider,omitempty"`
	VPCProjectUID     *string   `json:"vpcProjectUid,omitempty"`
	VPCNetworkName    *string   `json:"vpcNetworkName,omitempty"`
}

func (o CreateActiveActiveVPCPeering) String() string {
	return internal.ToString(o)
}

type listVpcPeering struct {
	Peerings []*VPCPeering `json:"peerings"`
}

type VPCPeering struct {
	ID               *int    `json:"vpcPeeringId,omitempty"`
	Status           *string `json:"status,omitempty"`
	AWSAccountID     *string `json:"awsAccountId,omitempty"`
	AWSPeeringID     *string `json:"awsPeeringUid,omitempty"`
	VPCId            *string `json:"vpcUid,omitempty"`
	VPCCidr          *string `json:"vpcCidr,omitempty"`
	VPCCidrs         []*CIDR `json:"vpcCidrs,omitempty"`
	GCPProjectUID    *string `json:"projectUid,omitempty"`
	NetworkName      *string `json:"networkName,omitempty"`
	RedisProjectUID  *string `json:"redisProjectUid,omitempty"`
	RedisNetworkName *string `json:"redisNetworkName,omitempty"`
	CloudPeeringID   *string `json:"cloudPeeringId,omitempty"`
	Region           *string `json:"regionName,omitempty"`
}

func (o VPCPeering) String() string {
	return internal.ToString(o)
}

type listActiveActiveVpcPeering struct {
	SubscriptionId *int                     `json:"subscriptionId,omitempty"`
	Regions        []*ActiveActiveVpcRegion `json:"regions,omitempty"`
}

type ActiveActiveVpcRegion struct {
	ID           *int                      `json:"id,omitempty"`
	SourceRegion *string                   `json:"region,omitempty"`
	VPCPeerings  []*ActiveActiveVPCPeering `json:"vpcPeerings,omitempty"`
}

type ActiveActiveVPCPeering struct {
	ID                *int    `json:"id,omitempty"`
	Status            *string `json:"status,omitempty"`
	RegionId          *int    `json:"regionId,omitempty"`
	RegionName        *string `json:"regionName,omitempty"`
	AWSAccountID      *string `json:"awsAccountId,omitempty"`
	AWSPeeringID      *string `json:"awsPeeringUid,omitempty"`
	VPCId             *string `json:"vpcUid,omitempty"`
	VPCCidr           *string `json:"vpcCidr,omitempty"`
	VPCCidrs          []*CIDR `json:"vpcCidrs,omitempty"`
	GCPProjectUID     *string `json:"vpcProjectUid,omitempty"`
	NetworkName       *string `json:"vpcNetworkName,omitempty"`
	RedisProjectUID   *string `json:"redisProjectUid,omitempty"`
	RedisNetworkName  *string `json:"redisNetworkName,omitempty"`
	CloudPeeringID    *string `json:"cloudPeeringId,omitempty"`
	SourceRegion      *string `json:"sourceRegion,omitempty"`
	DestinationRegion *string `json:"destinationRegion,omitempty"`
}

func (o ActiveActiveVPCPeering) String() string {
	return internal.ToString(o)
}

type CIDR struct {
	VPCCidr *string `json:"vpcCidr,omitempty"`
	Status  *string `json:"active,omitempty"`
}

func (o CIDR) String() string {
	return internal.ToString(o)
}

type listSubscriptionResponse struct {
	Subscriptions []*Subscription `json:"subscriptions"`
}

type ListAASubscriptionRegionsResponse struct {
	SubscriptionId *int                  `json:"subscriptionId,omitempty"`
	Regions        []*ActiveActiveRegion `json:"regions"`
}

// have to redeclare these here (copied from regions model) to avoid an import cycle
type ActiveActiveRegion struct {
	//RegionId       *int `json:"regionId,omitempty"` // not populated by the API
	Region         *string                `json:"region,omitempty"`
	DeploymentCIDR *string                `json:"deploymentCidr,omitempty"`
	VpcId          *string                `json:"vpcId,omitempty"`
	Databases      []ActiveActiveDatabase `json:"databases,omitempty"`
}
type ActiveActiveDatabase struct {
	DatabaseId               *int    `json:"databaseId,omitempty"`
	DatabaseName             *string `json:"databaseName,omitempty"`
	ReadOperationsPerSecond  *int    `json:"readOperationsPerSecond,omitempty"`
	WriteOperationsPerSecond *int    `json:"writeOperationsPerSecond,omitempty"`
}

type NotFound struct {
	ID int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("subscription %d not found", f.ID)
}

const (
	// SubscriptionStatusActive is the active value of the `Status` field in `Subscription`
	SubscriptionStatusActive = "active"
	// SubscriptionStatusPending is the pending value of the `Status` field in `Subscription`
	SubscriptionStatusPending = "pending"
	// SubscriptionStatusError is the error value of the `Status` field in `Subscription`
	SubscriptionStatusError = "error"
	// SubscriptionStatusDeleting is the deleting value of the `Status` field in `Subscription`
	SubscriptionStatusDeleting = "deleting"

	// VPCPeeringStatusInitiatingRequest is the initiating request value of the `Status` field in `VPCPeering`
	VPCPeeringStatusInitiatingRequest = "initiating-request"
	// VPCPeeringStatusActive is the active value of the `Status` field in `VPCPeering`
	VPCPeeringStatusActive = "active"
	// VPCPeeringStatusInactive is the inactive value of the `Status` field in `VPCPeering`
	VPCPeeringStatusInactive = "inactive"
	// VPCPeeringStatusPendingAcceptance is the pending acceptance value of the `Status` field in `VPCPeering`
	VPCPeeringStatusPendingAcceptance = "pending-acceptance"
	// VPCPeeringStatusFailed is the failed value of the `Status` field in `VPCPeering`
	VPCPeeringStatusFailed = "failed"

	SubscriptionDeploymentTypeSingleRegion = "single-region"
	SubscriptionDeploymentTypeActiveActive = "active-active"
)
