package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/fixed/plans"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const response = `{
	"plans": [
		{
			"id": 34944,
			"name": "Single-Zone_1GB",
			"size": 1,
			"sizeMeasurementUnit": "GB",
			"provider": "AWS",
			"region": "us-west-1",
			"regionId": 2,
			"price": 23,
			"priceCurrency": "USD",
			"pricePeriod": "Month",
			"maximumDatabases": 1,
			"maximumThroughput": 2000,
			"maximumBandwidthGB": 200,
			"availability": "Single-zone",
			"connections": "1024",
			"cidrAllowRules": 8,
			"supportDataPersistence": false,
			"supportInstantAndDailyBackups": true,
			"supportReplication": true,
			"supportClustering": false,
			"supportedAlerts": [
				"datasets-size",
				"throughput-higher-than",
				"latency",
				"throughput-lower-than",
				"connections-limit"
			],
			"customerSupport": "Standard",
			"links": []
		},
		{
			"id": 34947,
			"name": "Single-Zone_1GB",
			"size": 1,
			"sizeMeasurementUnit": "GB",
			"provider": "AWS",
			"region": "eu-west-1",
			"regionId": 4,
			"price": 23,
			"priceCurrency": "USD",
			"pricePeriod": "Month",
			"maximumDatabases": 1,
			"maximumThroughput": 2000,
			"maximumBandwidthGB": 200,
			"availability": "Single-zone",
			"connections": "1024",
			"cidrAllowRules": 8,
			"supportDataPersistence": false,
			"supportInstantAndDailyBackups": true,
			"supportReplication": true,
			"supportClustering": false,
			"supportedAlerts": [
				"datasets-size",
				"throughput-higher-than",
				"latency",
				"throughput-lower-than",
				"connections-limit"
			],
			"customerSupport": "Standard",
			"links": []
		}
	],
	"links": [
		{
			"href": "https://api-staging.qa.redislabs.com/v1/fixed/plans",
			"rel": "self",
			"type": "GET"
		}
	]
}`

func Test_List(t *testing.T) {
	s := httptest.NewServer(
		testServer("apiKey", "secret",
			getRequest(
				t,
				"/fixed/plans",
				response,
			),
		),
	)

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	list, err := subject.FixedPlans.List(context.TODO())

	require.NoError(t, err)

	assert.Equal(t, []*plans.GetPlanResponse{
		{
			ID:                            redis.Int(34944),
			Name:                          redis.String("Single-Zone_1GB"),
			Size:                          redis.Float64(1),
			SizeMeasurementUnit:           redis.String("GB"),
			Provider:                      redis.String("AWS"),
			Region:                        redis.String("us-west-1"),
			RegionID:                      redis.Int(2),
			Price:                         redis.Int(23),
			PriceCurrency:                 redis.String("USD"),
			PricePeriod:                   redis.String("Month"),
			MaximumDatabases:              redis.Int(1),
			MaximumThroughput:             redis.Int(2000),
			MaximumBandwidthGB:            redis.Int(200),
			Availability:                  redis.String("Single-zone"),
			Connections:                   redis.String("1024"),
			CidrAllowRules:                redis.Int(8),
			SupportDataPersistence:        redis.Bool(false),
			SupportInstantAndDailyBackups: redis.Bool(true),
			SupportReplication:            redis.Bool(true),
			SupportClustering:             redis.Bool(false),
			SupportedAlerts:               redis.StringSlice("datasets-size", "throughput-higher-than", "latency", "throughput-lower-than", "connections-limit"),
			CustomerSupport:               redis.String("Standard"),
		},
		{
			ID:                            redis.Int(34947),
			Name:                          redis.String("Single-Zone_1GB"),
			Size:                          redis.Float64(1),
			SizeMeasurementUnit:           redis.String("GB"),
			Provider:                      redis.String("AWS"),
			Region:                        redis.String("eu-west-1"),
			RegionID:                      redis.Int(4),
			Price:                         redis.Int(23),
			PriceCurrency:                 redis.String("USD"),
			PricePeriod:                   redis.String("Month"),
			MaximumDatabases:              redis.Int(1),
			MaximumThroughput:             redis.Int(2000),
			MaximumBandwidthGB:            redis.Int(200),
			Availability:                  redis.String("Single-zone"),
			Connections:                   redis.String("1024"),
			CidrAllowRules:                redis.Int(8),
			SupportDataPersistence:        redis.Bool(false),
			SupportInstantAndDailyBackups: redis.Bool(true),
			SupportReplication:            redis.Bool(true),
			SupportClustering:             redis.Bool(false),
			SupportedAlerts:               redis.StringSlice("datasets-size", "throughput-higher-than", "latency", "throughput-lower-than", "connections-limit"),
			CustomerSupport:               redis.String("Standard"),
		},
	}, list)
}

func Test_ListWithProvider(t *testing.T) {
	s := httptest.NewServer(
		testServer("apiKey", "secret",
			getRequestWithQuery(
				t,
				"/fixed/plans",
				map[string][]string{"provider": {"AWS"}},
				response,
			),
		),
	)

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	list, err := subject.FixedPlans.ListWithProvider(context.TODO(), "AWS")

	require.NoError(t, err)

	assert.Equal(t, []*plans.GetPlanResponse{
		{
			ID:                            redis.Int(34944),
			Name:                          redis.String("Single-Zone_1GB"),
			Size:                          redis.Float64(1),
			SizeMeasurementUnit:           redis.String("GB"),
			Provider:                      redis.String("AWS"),
			Region:                        redis.String("us-west-1"),
			RegionID:                      redis.Int(2),
			Price:                         redis.Int(23),
			PriceCurrency:                 redis.String("USD"),
			PricePeriod:                   redis.String("Month"),
			MaximumDatabases:              redis.Int(1),
			MaximumThroughput:             redis.Int(2000),
			MaximumBandwidthGB:            redis.Int(200),
			Availability:                  redis.String("Single-zone"),
			Connections:                   redis.String("1024"),
			CidrAllowRules:                redis.Int(8),
			SupportDataPersistence:        redis.Bool(false),
			SupportInstantAndDailyBackups: redis.Bool(true),
			SupportReplication:            redis.Bool(true),
			SupportClustering:             redis.Bool(false),
			SupportedAlerts:               redis.StringSlice("datasets-size", "throughput-higher-than", "latency", "throughput-lower-than", "connections-limit"),
			CustomerSupport:               redis.String("Standard"),
		},
		{
			ID:                            redis.Int(34947),
			Name:                          redis.String("Single-Zone_1GB"),
			Size:                          redis.Float64(1),
			SizeMeasurementUnit:           redis.String("GB"),
			Provider:                      redis.String("AWS"),
			Region:                        redis.String("eu-west-1"),
			RegionID:                      redis.Int(4),
			Price:                         redis.Int(23),
			PriceCurrency:                 redis.String("USD"),
			PricePeriod:                   redis.String("Month"),
			MaximumDatabases:              redis.Int(1),
			MaximumThroughput:             redis.Int(2000),
			MaximumBandwidthGB:            redis.Int(200),
			Availability:                  redis.String("Single-zone"),
			Connections:                   redis.String("1024"),
			CidrAllowRules:                redis.Int(8),
			SupportDataPersistence:        redis.Bool(false),
			SupportInstantAndDailyBackups: redis.Bool(true),
			SupportReplication:            redis.Bool(true),
			SupportClustering:             redis.Bool(false),
			SupportedAlerts:               redis.StringSlice("datasets-size", "throughput-higher-than", "latency", "throughput-lower-than", "connections-limit"),
			CustomerSupport:               redis.String("Standard"),
		},
	}, list)
}
