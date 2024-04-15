package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/12/databases/34/pricing",
				`[
				  {
					"type": "Shards",
					"typeDetails": "micro",
					"quantity": 1,
					"quantityMeasurement": "shards",
					"pricePerUnit": 0.027,
					"priceCurrency": "USD",
					"pricePeriod": "hour"
				  },
				  [
					{
					  "rel": "self",
					  "href": "https://api-staging.qa.redislabs.com/v1/subscriptions/110777/pricing",
					  "type": "GET"
					}
				  ]
				]`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	_, err = subject.Pricing.List(context.TODO(), 20000)
	require.NoError(t, err)

	//assert.Equal(t, &users.GetUserResponse{
	//	ID:     redis.Int(20000),
	//	Name:   redis.String("test-user"),
	//	Role:   redis.String("test-role"),
	//	Status: redis.String(users.StatusPending),
	//}, actual)
}
