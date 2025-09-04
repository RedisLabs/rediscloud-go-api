package rediscloud_api

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	pl "github.com/RedisLabs/rediscloud-go-api/service/privatelink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPrivateLinkConfig(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *pl.PrivateLink
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a privatelink config",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/private-link",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "privateLinkGetRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-16T09:26:40.929904847Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
				),
				getRequest(
					t,
					"/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "privatelinkGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "subscriptionId": 114019
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
				),
			},
			expectedResult: &pl.PrivateLink{
				SubscriptionId: redis.Int(114019),
			},
		},
		{
			description: "should fail when private link is not found",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/private-link",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "privatelinkGetRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-16T09:26:40.929904847Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
				),
				getRequest(
					t,
					"/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "privatelinkGetRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2025-01-13T11:22:51.204189721Z",
				  "response": {
					"error": {
					  "type": "PRIVATELINK_SERVICE_NOT_FOUND",
					  "status": "404 NOT_FOUND",
					  "description": "Private Service Connect service not found"
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
				),
			},
			expectedError:   errors.New("resource not found - subscription 114019"),
			expectedErrorAs: &pl.NotFound{},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithStatus(
					t,
					"/subscriptions/114019/private-link",
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/private-link"
					}`),
			},
			expectedError:   errors.New("resource not found - subscription 114019"),
			expectedErrorAs: &pl.NotFound{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateLink.GetPrivateLink(context.TODO(), 114019)
			if testCase.expectedError == nil {
				assert.NoError(t, err)
				assert.Equal(t, testCase.expectedResult, actual)
			} else {
				assert.IsType(t, err, testCase.expectedErrorAs)
				assert.EqualError(t, err, testCase.expectedError.Error())
			}
		})
	}
}
