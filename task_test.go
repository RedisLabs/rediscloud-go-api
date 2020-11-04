package rediscloud_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTask_Get(t *testing.T) {
	resourceId := 100556

	s := httptest.NewServer(testServer("key", "secret", getRequest(t, "/tasks/task-uuid", fmt.Sprintf(`{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "subscriptionCreateRequest",
  "status": "processing-error",
  "description": "Task request failed during processing. See error information for failure details.",
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
}`, resourceId))))

	subject, err := NewClient(BaseUrl(s.URL), Auth("key", "secret"), Transporter(s.Client().Transport))
	require.NoError(t, err)

	actual, err := subject.Task.Get(context.TODO(), "task-uuid")
	require.NoError(t, err)

	assert.Equal(t, &task.Task{
		CommandType: redis.String("subscriptionCreateRequest"),
		Description: redis.String("Task request failed during processing. See error information for failure details."),
		Status:      redis.String("processing-error"),
		ID:          redis.String("e02b40d6-1395-4861-a3b9-ecf829d835fd"),
		Response: &task.Response{
			ID: &resourceId,
		},
	}, actual)
}

func TestTask_Get_UnwrapsTaskError(t *testing.T) {
	s := httptest.NewServer(testServer("key", "secret", getRequest(t, "/tasks/task-uuid", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "subscriptionCreateRequest",
  "status": "processing-error",
  "description": "Task request failed during processing. See error information for failure details.",
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

	subject, err := NewClient(BaseUrl(s.URL), Auth("key", "secret"), Transporter(s.Client().Transport))
	require.NoError(t, err)

	actual, err := subject.Task.Get(context.TODO(), "task-uuid")
	assert.Equal(t, &task.Error{
		Type:        redis.String("SUBSCRIPTION_PI_NOT_FOUND"),
		Description: redis.String("Payment info was not found for subscription. Use 'GET /payment-methods' to lookup valid payment methods for current Account"),
		Status:      redis.String("400 BAD_REQUEST"),
	}, err)
	assert.Nil(t, actual)
}

func TestTask_WaitForTaskToComplete(t *testing.T) {
	resourceId := 100556
	resource := "oiuygfcvbnmk"

	s := httptest.NewServer(testServer("key", "secret", getRequest(t, "/tasks/task-uuid", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "subscriptionCreateRequest",
  "status": "initialized",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {},
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task-uuid", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "subscriptionCreateRequest",
  "status": "processing-in-progress",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {},
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task-uuid", fmt.Sprintf(`{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "subscriptionCreateRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
    "resourceId": %d,
    "resource": "%s"
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`, resourceId, resource))))

	subject, err := NewClient(BaseUrl(s.URL), Auth("key", "secret"), Transporter(s.Client().Transport))
	require.NoError(t, err)

	actual, err := subject.Task.WaitForTaskToComplete(context.TODO(), "task-uuid")
	require.NoError(t, err)
	assert.Equal(t, resourceId, *actual.Response.ID)

	var actualResponse string
	err = json.Unmarshal(*actual.Response.Resource, &actualResponse)
	require.NoError(t, err)

	assert.Equal(t, resource, actualResponse)
}
