package regions

import (
	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type Regions struct {
	SubscriptionId *int      `json:"subscriptionId,omitempty"`
	Regions        []*Region `json:"regions,omitempty"`
}

func (o Regions) String() string {
	return internal.ToString(o)
}

type Region struct {
	RegionId       *int        `json:"regionId,omitempty"`
	Region         *string     `json:"region,omitempty"`
	RecreateRegion *bool       `json:"-"`
	DeploymentCIDR *string     `json:"deploymentCIDR,omitempty"`
	VpcId          *string     `json:"vpcId,omitempty"`
	Databases      []*Database `json:"databases,omitempty"`
}

func (o Region) String() string {
	return internal.ToString(o)
}

type Database struct {
	DatabaseId               *int    `json:"databaseId,omitempty"`
	DatabaseName             *string `json:"DatabaseName,omitempty"`
	ReadOperationsPerSecond  *int    `json:"readOperationsPerSecond,omitempty"`
	WriteOperationsPerSecond *int    `json:"writeOperationsPerSecond,omitempty"`
}

type CreateRegion struct {
	Region         *string           `json:"region,omitempty"`
	DeploymentCIDR *string           `json:"deploymentCIDR,omitempty"`
	DryRun         *bool             `json:"dryRun,omitempty"`
	Databases      []*CreateDatabase `json:"databases,omitempty"`
}

type DeleteRegion struct {
	Region *string `json:"region,omitempty"`
}
type DeleteRegions struct {
	Regions []*DeleteRegion `json:"regions,omitempty"`
}

type CreateLocalThroughput struct {
	Region                   *string `json:"region,omitempty"`
	WriteOperationsPerSecond *int    `json:"writeOperationsPerSecond"`
	ReadOperationsPerSecond  *int    `json:"readOperationsPerSecond"`
}

type CreateDatabase struct {
	Name                       *string                `json:"name,omitempty"`
	LocalThroughputMeasurement *CreateLocalThroughput `json:"localThroughputMeasurement,omitempty"`
}
