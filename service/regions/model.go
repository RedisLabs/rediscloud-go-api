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
