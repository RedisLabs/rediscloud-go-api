package rediscloud_api

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskErrorsGetUnwrapped(t *testing.T) {
	s := httptest.NewServer(testServer("key", "secret", deleteRequest(t, "/cloud-accounts/1", `{
  "taskId": "task",
  "commandType": "cloudAccountDeleteRequest",
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
  "commandType": "cloudAccountDeleteRequest",
  "status": "processing-error",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
    "error": {
      "type": "SUBSCRIPTION_PI_NOT_FOUND",
      "status": "400 BAD_REQUEST",
      "description": "Payment info was not found for subscription. Use 'GET /payment-methods' to lookup valid payment methods for current Account"
    }
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	err = subject.CloudAccount.Delete(context.TODO(), 1)
	assert.Equal(t, &internal.Error{
		Type:        redis.String("SUBSCRIPTION_PI_NOT_FOUND"),
		Description: redis.String("Payment info was not found for subscription. Use 'GET /payment-methods' to lookup valid payment methods for current Account"),
		Status:      redis.String("400 BAD_REQUEST"),
	}, errors.Unwrap(err))
}

func TestTask_Handles404Eventually(t *testing.T) {
	s := httptest.NewServer(testServer("key", "secret", deleteRequest(t, "/cloud-accounts/1", `{
  "taskId": "task",
  "commandType": "cloudAccountDeleteRequest",
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
  "commandType": "cloudAccountDeleteRequest",
  "status": "initialized",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {},
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`), getRequestWithStatus(t, "/tasks/task", 404, ""),
		getRequestWithStatus(t, "/tasks/task", 404, ""),
		getRequestWithStatus(t, "/tasks/task", 404, ""),
		getRequestWithStatus(t, "/tasks/task", 404, ""),
		getRequestWithStatus(t, "/tasks/task", 404, "")))

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	err = subject.CloudAccount.Delete(context.TODO(), 1)
	assert.Error(t, err)
}
