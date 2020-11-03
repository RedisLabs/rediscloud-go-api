package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/service/databases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDatabase_List(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequest(t, "/subscriptions/23456/databases", `{
  "accountId": 2,
  "subscription": [
    {
      "subscriptionId": 23456,
      "databases": [
        {
          "databaseId": 42,
          "name": "first-example",
          "protocol": "redis",
          "provider": "AWS",
          "region": "eu-west-1"
        },
        {
          "databaseId": 43,
          "name": "second-example",
          "protocol": "redis",
          "provider": "AWS",
          "region": "eu-west-1"
        }
      ]
    }
  ],
  "_links": {
    "self": {
      "href": "https://example.org",
      "type": "GET"
    }
  }
}`)))

	subject, err := NewClient(BaseUrl(s.URL), Auth("apiKey", "secret"), Transporter(s.Client().Transport))
	require.NoError(t, err)

	actual, err := subject.Database.List(context.TODO(), 23456)
	require.NoError(t, err)

	assert.ElementsMatch(t, []*databases.Database{
		{
			ID:       42,
			Name:     "first-example",
			Protocol: "redis",
			Provider: "AWS",
			Region:   "eu-west-1",
		},
		{
			ID:       43,
			Name:     "second-example",
			Protocol: "redis",
			Provider: "AWS",
			Region:   "eu-west-1",
		},
	}, actual)
}

func TestDatabase_Get(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequest(t, "/subscriptions/23456/databases/98765", `{
  "databaseId": 98765,
  "name": "Example",
  "protocol": "redis",
  "provider": "AWS",
  "region": "eu-west-1",
  "redisVersionCompliance": "6.0.5",
  "status": "active",
  "memoryLimitInGb": 7,
  "memoryUsedInMb": 5,
  "memoryStorage": "ram",
  "supportOSSClusterApi": true,
  "dataPersistence": "none",
  "replication": false,
  "dataEvictionPolicy": "volatile-random",
  "throughputMeasurement": {
    "by": "operations-per-second",
    "value": 10000
  },
  "activatedOn": "2020-11-03T09:03:30Z",
  "lastModified": "2020-11-03T09:03:30Z",
  "publicEndpoint": "example.com:16668",
  "privateEndpoint": "example.net:16668",
  "replicaOf": null,
  "clustering": {
    "numberOfShards": 1,
    "regexRules": [
      {
        "ordinal": 1,
        "pattern": "(?<tag>.*)"
      },
      {
        "ordinal": 0,
        "pattern": ".*\\{(?<tag>.*)\\}.*"
      }
    ],
    "hashingPolicy": "custom"
  },
  "security": {
    "password": "test",
    "sslClientAuthentication": false,
    "sourceIps": [
      "0.0.0.0/0"
    ]
  },
  "modules": [],
  "alerts": [],
  "_links": {
    "self": {
      "href": "https://example.org",
      "type": "GET"
    }
  }
}`)))

	subject, err := NewClient(BaseUrl(s.URL), Auth("apiKey", "secret"), Transporter(s.Client().Transport))
	require.NoError(t, err)

	actual, err := subject.Database.Get(context.TODO(), 23456, 98765)
	require.NoError(t, err)

	assert.Equal(t, &databases.Database{
		ID:                     98765,
		Name:                   "Example",
		Protocol:               "redis",
		Provider:               "AWS",
		Region:                 "eu-west-1",
		Status:                 "active",
		MemoryLimitInGb:        7,
		MemoryUsedInMb:         5,
		SupportOSSClusterApi:   true,
		DataPersistence:        "none",
		Replication:            false,
		DataEvictionPolicy:     "volatile-random",
		ActivatedOn:            time.Date(2020, 11, 3, 9, 3, 30, 0, time.UTC),
		LastModified:           time.Date(2020, 11, 3, 9, 3, 30, 0, time.UTC),
		MemoryStorage:          "ram",
		PrivateEndpoint:        "example.net:16668",
		PublicEndpoint:         "example.com:16668",
		RedisVersionCompliance: "6.0.5",
	}, actual)
}

func TestDatabase_Delete(t *testing.T) {
	s := httptest.NewServer(testServer("key", "secret", deleteRequest(t, "/subscriptions/42/databases/4291", `{
  "taskId": "task",
  "commandType": "databaseDeleteRequest",
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
  "commandType": "databaseDeleteRequest",
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

	subject, err := NewClient(BaseUrl(s.URL), Auth("key", "secret"), Transporter(s.Client().Transport))
	require.NoError(t, err)

	err = subject.Database.Delete(context.TODO(), 42, 4291)
	require.NoError(t, err)
}
