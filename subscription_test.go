package rediscloud_api

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/subscriptions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscription_Create(t *testing.T) {
	expected := 1235
	s := httptest.NewServer(testServer("key", "secret", postRequest(t, "/subscriptions", `{
  "name": "Test subscription",
  "dryRun": false,
  "paymentMethodId": 2,
  "memoryStorage": "ram",
  "persistentStorageEncryption": false,
  "cloudProviders": [
    {
      "provider": "AWS",
      "cloudAccountId": 1,
      "regions": [
        {
          "region": "eu-west-1"
        }
      ]
    }
  ],
  "databases": [
    {
      "name": "example",
      "protocol": "redis",
      "memoryLimitInGb": 1,
      "supportOSSClusterApi": true,
      "dataPersistence": "none",
      "replication": false,
      "throughputMeasurement": {
        "by": "operations-per-second",
        "value": 10000
      },
      "quantity": 1
    }
  ]
}`, `{
  "taskId": "task-id",
  "commandType": "subscriptionCreateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task-id", fmt.Sprintf(`{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "subscriptionCreateRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
    "resourceId": %d
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`, expected))))

	subject, err := NewClient(BaseUrl(s.URL), Auth("key", "secret"), Transporter(s.Client().Transport))
	require.NoError(t, err)

	actual, err := subject.Subscription.Create(context.TODO(), subscriptions.CreateSubscription{
		Name:                        redis.String("Test subscription"),
		DryRun:                      redis.Bool(false),
		PaymentMethodID:             redis.Int(2),
		MemoryStorage:               redis.String("ram"),
		PersistentStorageEncryption: redis.Bool(false),
		CloudProviders: []*subscriptions.CreateCloudProvider{
			{
				Provider:       redis.String("AWS"),
				CloudAccountID: redis.Int(1),
				Regions: []*subscriptions.CreateRegion{
					{
						Region: redis.String("eu-west-1"),
					},
				},
			},
		},
		Databases: []*subscriptions.CreateDatabase{
			{
				Name:                 redis.String("example"),
				Protocol:             redis.String("redis"),
				MemoryLimitInGB:      redis.Float64(1),
				SupportOSSClusterAPI: redis.Bool(true),
				DataPersistence:      redis.String("none"),
				Replication:          redis.Bool(false),
				ThroughputMeasurement: &subscriptions.CreateThroughput{
					By:    redis.String("operations-per-second"),
					Value: redis.Int(10000),
				},
				Quantity: redis.Int(1),
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSubscription_Delete(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", deleteRequest(t, "/subscriptions/12356", `{
  "taskId": "task",
  "commandType": "subscriptionDeleteRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "subscriptionDeleteRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := NewClient(BaseUrl(s.URL), Auth("apiKey", "secret"), Transporter(s.Client().Transport))
	require.NoError(t, err)

	err = subject.Subscription.Delete(context.TODO(), 12356)
	require.NoError(t, err)
}
