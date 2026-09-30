package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/service/agentmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentMemory_Create(t *testing.T) {
	s := httptest.NewServer(testServer("key", "secret", postRequest(t, "/memory-stores", `{
  "name": "store",
  "databaseId": 123
}`, `{
  "taskId": "task",
  "commandType": "memoryStoreCreateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2026-09-22T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", `{
  "taskId": "task",
  "commandType": "memoryStoreCreateRequest",
  "status": "processing-completed",
  "timestamp": "2026-09-22T09:05:35.3Z",
  "response": {
    "resource": {
      "storeId": "store-1"
    }
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))
	defer s.Close()

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.AgentMemory.Create(context.TODO(), agentmemory.CreateStore{
		Name:       "store",
		DatabaseID: 123,
	})

	require.NoError(t, err)
	assert.Equal(t, "store-1", actual)
}

func TestAgentMemory_Get(t *testing.T) {
	s := httptest.NewServer(testServer("key", "secret", getRequest(t, "/memory-stores/store-1", `{
  "storeId": "store-1",
  "name": "store",
  "databaseId": "123",
  "endpoint": "https://aws-us-east-1.memory.redis.io",
  "endpoints": [
    {
      "url": "https://aws-us-east-1.memory.redis.io",
      "provider": "AWS",
      "region": "us-east-1",
      "egressIps": ["10.0.0.1"],
      "isAccessible": true
    },
    {
      "url": "https://gcp-us-east4.memory.redis.io",
      "provider": "GCP",
      "region": "us-east4",
      "egressIps": [],
      "isAccessible": false
    }
  ],
  "status": "READY"
}`)))
	defer s.Close()

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.AgentMemory.Get(context.TODO(), "store-1")

	require.NoError(t, err)
	require.Len(t, actual.Endpoints, 2)
	assert.Equal(t, "https://aws-us-east-1.memory.redis.io", actual.Endpoint)
	assert.Equal(t, "AWS", actual.Endpoints[0].Provider)
	assert.Equal(t, []string{"10.0.0.1"}, actual.Endpoints[0].EgressIPs)
	assert.Equal(t, "GCP", actual.Endpoints[1].Provider)
	assert.Empty(t, actual.Endpoints[1].EgressIPs)
	assert.False(t, actual.Endpoints[1].IsAccessible)
}

func TestAgentMemory_CreateAPIKey(t *testing.T) {
	s := httptest.NewServer(testServer("key", "secret", getRequest(t, "/memory-stores/store-1/api-keys", `{
  "apiKeys": [
    {
      "apiKeyId": "existing-id",
      "keyName": "existing",
      "obfuscatedToken": "mem1...old",
      "createdAt": 1234
    }
  ]
}`), postRequest(t, "/memory-stores/store-1/api-keys", `{
  "keyName": "terraform"
}`, `{
  "apiKey": "secret-token"
}`), getRequest(t, "/memory-stores/store-1/api-keys", `{
  "apiKeys": [
    {
      "apiKeyId": "existing-id",
      "keyName": "existing",
      "obfuscatedToken": "mem1...old",
      "createdAt": 1234
    },
    {
      "apiKeyId": "new-id",
      "keyName": "terraform",
      "obfuscatedToken": "mem1...new",
      "createdAt": 5678
    }
  ]
}`)))
	defer s.Close()

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.AgentMemory.CreateAPIKey(context.TODO(), "store-1", "terraform")

	require.NoError(t, err)
	assert.Equal(t, "secret-token", actual.APIKey)
	assert.Equal(t, "new-id", actual.Info.APIKeyID)
}
