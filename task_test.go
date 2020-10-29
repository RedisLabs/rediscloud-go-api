package rediscloud_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/service/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTask_Get(t *testing.T) {
	resourceId := 100556

	s := httptest.NewServer(testServer("/tasks/task-uuid", "key", "secret", fmt.Sprintf(`{
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
}`, resourceId)))

	subject := NewClient(BaseUrl(s.URL), Auth("key", "secret"), Transporter(s.Client().Transport))

	actual, err := subject.Task.Get(context.TODO(), "task-uuid")
	require.NoError(t, err)

	assert.Equal(t, &task.Task{
		CommandType: "subscriptionCreateRequest",
		Description: "Task request failed during processing. See error information for failure details.",
		Status:      "processing-error",
		Id:          "e02b40d6-1395-4861-a3b9-ecf829d835fd",
		Response: &task.Response{
			Id: &resourceId,
		},
	}, actual)
}

func TestTask_Get_UnwrapsTaskError(t *testing.T) {
	s := httptest.NewServer(testServer("/tasks/task-uuid", "key", "secret", `{
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
}`))

	subject := NewClient(BaseUrl(s.URL), Auth("key", "secret"), Transporter(s.Client().Transport))

	actual, err := subject.Task.Get(context.TODO(), "task-uuid")
	assert.Equal(t, &task.Error{
		Type:        "SUBSCRIPTION_PI_NOT_FOUND",
		Description: "Payment info was not found for subscription. Use 'GET /payment-methods' to lookup valid payment methods for current Account",
		Status:      "400 BAD_REQUEST",
	}, err)
	assert.Nil(t, actual)
}

func TestTask_WaitForIdFromTask(t *testing.T) {
	resourceId := 100556
	resource := "oiuygfcvbnmk"

	s := httptest.NewServer(testServer("/tasks/task-uuid", "key", "secret", `{
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
}`, `{
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
}`, fmt.Sprintf(`{
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
}`, resourceId, resource)))

	subject := NewClient(BaseUrl(s.URL), Auth("key", "secret"), Transporter(s.Client().Transport))

	actual, err := subject.Task.WaitForTaskToComplete(context.TODO(), "task-uuid")
	require.NoError(t, err)
	assert.Equal(t, resourceId, *actual.Response.Id)

	var actualResponse string
	err = json.Unmarshal(*actual.Response.Resource, &actualResponse)
	require.NoError(t, err)

	assert.Equal(t, resource, actualResponse)
}

func testServer(path, apiKey, secretKey string, responses ...string) http.HandlerFunc {
	responseCount := 0
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(500)
			return
		}
		if r.URL.Path != path {
			w.WriteHeader(501)
			return
		}
		if r.Header.Get("X-Api-Key") != apiKey {
			w.WriteHeader(502)
			return
		}
		if r.Header.Get("X-Api-Secret-Key") != secretKey {
			w.WriteHeader(503)
			return
		}

		body := responses[responseCount]
		responseCount++
		if responseCount > len(responses)-1 {
			responseCount = len(responses) - 2
		}
		w.WriteHeader(200)
		_, _ = w.Write([]byte(body))
	}
}
