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

const responseBody = `
{
  "plans": [
    {
      "id": 98183,
      "name": "Multi-AZ 5GB",
      "size": 5,
      "sizeMeasurementUnit": "GB",
      "provider": "AWS",
      "region": "us-east-1",
      "regionId": 1,
      "price": 100,
      "priceCurrency": "USD",
      "pricePeriod": "Month",
      "maximumDatabases": 1,
      "availability": "Multi-zone",
      "connections": "unlimited",
      "cidrAllowRules": 16,
      "supportDataPersistence": true,
      "supportInstantAndDailyBackups": true,
      "supportReplication": true,
      "supportClustering": false,
      "supportedAlerts": [
        "datasets-size",
        "latency",
        "throughput-higher-than",
        "throughput-lower-than"
      ],
      "customerSupport": "Standard",
      "links": []
    },
    {
      "id": 98181,
      "name": "Multi-AZ 1GB",
      "size": 1,
      "sizeMeasurementUnit": "GB",
      "provider": "AWS",
      "region": "us-east-1",
      "regionId": 1,
      "price": 22,
      "priceCurrency": "USD",
      "pricePeriod": "Month",
      "maximumDatabases": 1,
      "availability": "Multi-zone",
      "connections": "1024",
      "cidrAllowRules": 8,
      "supportDataPersistence": true,
      "supportInstantAndDailyBackups": true,
      "supportReplication": true,
      "supportClustering": false,
      "supportedAlerts": [
        "datasets-size",
        "throughput-higher-than",
        "throughput-lower-than",
        "latency",
        "connections-limit"
      ],
      "customerSupport": "Standard",
      "links": []
    }
  ],
  "links": [
    {
      "rel": "self",
      "href": "http://localhost:8081/v1/fixed/plans?cloud_provider=AWS",
      "type": "GET"
    }
  ]
}`

func Test_Plans_Subscriptions_List(t *testing.T) {
	s := httptest.NewServer(
		testServer("apiKey", "secret",
			getRequest(
				t,
				"/fixed/plans/subscriptions/98183",
				responseBody,
			),
		),
	)

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actualResponse, err := subject.FixedPlanSubscriptions.List(context.TODO(), 98183)

	require.NoError(t, err)

	expectedResponse := []*plans.GetPlanResponse{
		{
			ID:                            redis.Int(98183),
			Name:                          redis.String("Multi-AZ 5GB"),
			Size:                          redis.Float64(5),
			SizeMeasurementUnit:           redis.String("GB"),
			Provider:                      redis.String("AWS"),
			Region:                        redis.String("us-east-1"),
			RegionID:                      redis.Int(1),
			Price:                         redis.Int(100),
			PriceCurrency:                 redis.String("USD"),
			PricePeriod:                   redis.String("Month"),
			MaximumDatabases:              redis.Int(1),
			Availability:                  redis.String("Multi-zone"),
			Connections:                   redis.String("unlimited"),
			CidrAllowRules:                redis.Int(16),
			SupportDataPersistence:        redis.Bool(true),
			SupportInstantAndDailyBackups: redis.Bool(true),
			SupportReplication:            redis.Bool(true),
			SupportClustering:             redis.Bool(false),
			SupportedAlerts: redis.StringSlice(
				"datasets-size",
				"latency",
				"throughput-higher-than",
				"throughput-lower-than"),
			CustomerSupport: redis.String("Standard"),
		},
		{
			ID:                            redis.Int(98181),
			Name:                          redis.String("Multi-AZ 1GB"),
			Size:                          redis.Float64(1),
			SizeMeasurementUnit:           redis.String("GB"),
			Provider:                      redis.String("AWS"),
			Region:                        redis.String("us-east-1"),
			RegionID:                      redis.Int(1),
			Price:                         redis.Int(22),
			PriceCurrency:                 redis.String("USD"),
			PricePeriod:                   redis.String("Month"),
			MaximumDatabases:              redis.Int(1),
			Availability:                  redis.String("Multi-zone"),
			Connections:                   redis.String("1024"),
			CidrAllowRules:                redis.Int(8),
			SupportDataPersistence:        redis.Bool(true),
			SupportInstantAndDailyBackups: redis.Bool(true),
			SupportReplication:            redis.Bool(true),
			SupportClustering:             redis.Bool(false),
			SupportedAlerts: redis.StringSlice(
				"datasets-size",
				"throughput-higher-than",
				"throughput-lower-than",
				"latency",
				"connections-limit"),
			CustomerSupport: redis.String("Standard"),
		},
	}
	assert.Equal(t, expectedResponse, actualResponse)
}
