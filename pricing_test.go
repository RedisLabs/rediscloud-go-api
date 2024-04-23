package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/pricing"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/20000/pricing",
				`{
					"pricing": [
						{
							"type": "Shards",
							"typeDetails": "micro",
							"quantity": 1,
							"quantityMeasurement": "shards",
							"pricePerUnit": 0.027,
							"priceCurrency": "USD",
							"pricePeriod": "hour"
						}
					],
					"links": [
						{
							"rel": "self",
							"href": "https://api-staging.qa.redislabs.com/v1/subscriptions/110777/pricing",
							"type": "GET"
						}
					]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Pricing.List(context.TODO(), 20000)
	require.NoError(t, err)

	assert.Equal(t, []*pricing.Pricing{
		{
			Type:                redis.String("Shards"),
			TypeDetails:         redis.String("micro"),
			Quantity:            redis.Int(1),
			QuantityMeasurement: redis.String("shards"),
			PricePerUnit:        redis.Float64(0.027),
			PriceCurrency:       redis.String("USD"),
			PricePeriod:         redis.String("hour"),
		},
	}, actual)
}
