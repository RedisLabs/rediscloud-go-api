package plans

import "github.com/RedisLabs/rediscloud-go-api/internal"

// Get

type ListPlansResponse struct {
	Plans []*GetPlanResponse `json:"plans,omitempty"`
}

func (o ListPlansResponse) String() string {
	return internal.ToString(o)
}

type GetPlanResponse struct {
	ID                            *int      `json:"id,omitempty"`
	Name                          *string   `json:"name,omitempty"`
	Size                          *float64  `json:"size,omitempty"`
	SizeMeasurementUnit           *string   `json:"sizeMeasurementUnit,omitempty"`
	Provider                      *string   `json:"provider,omitempty"`
	Region                        *string   `json:"region,omitempty"`
	RegionID                      *int      `json:"regionId,omitempty"`
	Price                         *int      `json:"price,omitempty"`
	PriceCurrency                 *string   `json:"priceCurrency,omitempty"`
	PricePeriod                   *string   `json:"pricePeriod,omitempty"`
	MaximumDatabases              *int      `json:"maximumDatabases,omitempty"`
	MaximumThroughput             *int      `json:"maximumThroughput,omitempty"`
	MaximumBandwidthGB            *int      `json:"maximumBandwidthGB,omitempty"`
	Availability                  *string   `json:"availability,omitempty"`
	Connections                   *string   `json:"connections,omitempty"` // Could be a number or 'unlimited'
	CidrAllowRules                *int      `json:"cidrAllowRules,omitempty"`
	SupportDataPersistence        *bool     `json:"supportDataPersistence,omitempty"`
	SupportInstantAndDailyBackups *bool     `json:"supportInstantAndDailyBackups,omitempty"`
	SupportReplication            *bool     `json:"supportReplication,omitempty"`
	SupportClustering             *bool     `json:"supportClustering,omitempty"`
	SupportedAlerts               []*string `json:"supportedAlerts,omitempty"`
	CustomerSupport               *string   `json:"customerSupport,omitempty"`
	SubscriptionID                *int      `json:"subscriptionId,omitempty"`
}

func (o GetPlanResponse) String() string {
	return internal.ToString(o)
}
