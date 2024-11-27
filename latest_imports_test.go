package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/latest_imports"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestGetLatestImportTooEarly(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/12/databases/34/import",
				`{
				  "taskId": "1dfd6084-21df-40c6-829c-e9b4790e207e",
				  "commandType": "databaseImportStatusRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-04-15T10:19:06.710686Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/1dfd6084-21df-40c6-829c-e9b4790e207e",
					  "type": "GET",
					  "rel": "task"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/1dfd6084-21df-40c6-829c-e9b4790e207e",
				`{
				  "taskId": "1dfd6084-21df-40c6-829c-e9b4790e207e",
				  "commandType": "databaseImportStatusRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2024-04-15T10:19:07.331898Z",
				  "response": {
					"error": {
					  "type": "SUBSCRIPTION_NOT_ACTIVE",
					  "status": "403 FORBIDDEN",
					  "description": "Cannot preform any actions for subscription that is not in an active state"
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/1dfd6084-21df-40c6-829c-e9b4790e207e",
					  "type": "GET",
					  "rel": "self"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.LatestImport.Get(context.TODO(), 12, 34)
	require.NoError(t, err)

	assert.Equal(t, &latest_imports.LatestImportStatus{
		CommandType: redis.String("databaseImportStatusRequest"),
		Description: redis.String("Task request failed during processing. See error information for failure details."),
		Status:      redis.String("processing-error"),
		ID:          redis.String("1dfd6084-21df-40c6-829c-e9b4790e207e"),
		Response: &latest_imports.Response{
			Error: &latest_imports.Error{
				Type:        redis.String("SUBSCRIPTION_NOT_ACTIVE"),
				Description: redis.String("Cannot preform any actions for subscription that is not in an active state"),
				Status:      redis.String("403 FORBIDDEN"),
			},
		},
	}, actual)
}

func TestGetFixedLatestImport(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/fixed/subscriptions/12/databases/34/import",
				`{
				  "taskId": "e9232e43-3781-4263-a38e-f4d150e03475",
				  "commandType": "databaseImportStatusRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-04-15T10:44:34.325298Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/e9232e43-3781-4263-a38e-f4d150e03475",
					  "type": "GET",
					  "rel": "task"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/e9232e43-3781-4263-a38e-f4d150e03475",
				`{
				  "taskId": "e9232e43-3781-4263-a38e-f4d150e03475",
				  "commandType": "databaseImportStatusRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-04-15T10:44:35.225468Z",
				  "response": {
					"resourceId": 51051302,
					"additionalResourceId": 110777,
					"resource": {
                      "status": "importing"
                	}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/e9232e43-3781-4263-a38e-f4d150e03475",
					  "type": "GET",
					  "rel": "self"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.LatestImport.GetFixed(context.TODO(), 12, 34)
	require.NoError(t, err)

	assert.Equal(t, &latest_imports.LatestImportStatus{
		CommandType: redis.String("databaseImportStatusRequest"),
		Description: redis.String("Request processing completed successfully and its resources are now being provisioned / de-provisioned."),
		Status:      redis.String("processing-completed"),
		ID:          redis.String("e9232e43-3781-4263-a38e-f4d150e03475"),
		Response: &latest_imports.Response{
			ID: redis.Int(51051302),
			Resource: &latest_imports.Resource{
				Status: redis.String("importing"),
			},
			Error: nil,
		},
	}, actual)
}

func TestGetLatestImport(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/12/databases/34/import",
				`{
				  "taskId": "e9232e43-3781-4263-a38e-f4d150e03475",
				  "commandType": "databaseImportStatusRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-04-15T10:44:34.325298Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/e9232e43-3781-4263-a38e-f4d150e03475",
					  "type": "GET",
					  "rel": "task"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/e9232e43-3781-4263-a38e-f4d150e03475",
				`{
				  "taskId": "e9232e43-3781-4263-a38e-f4d150e03475",
				  "commandType": "databaseImportStatusRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-04-15T10:44:35.225468Z",
				  "response": {
					"resourceId": 51051302,
					"additionalResourceId": 110777,
					"resource": {
					  "failureReason": "file-corrupted",
					  "failureReasonParams": [
						{
						  "key": "bytes_configured_bdb_limit",
						  "value": "1234"
						},
						{
						  "key": "bytes_of_expected_dataset",
						  "value": "5678"
						}
					  ],
					  "lastImportTime": "2024-05-21T10:36:26Z",
                      "status": "failed"
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/e9232e43-3781-4263-a38e-f4d150e03475",
					  "type": "GET",
					  "rel": "self"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.LatestImport.Get(context.TODO(), 12, 34)
	require.NoError(t, err)

	assert.Equal(t, &latest_imports.LatestImportStatus{
		CommandType: redis.String("databaseImportStatusRequest"),
		Description: redis.String("Request processing completed successfully and its resources are now being provisioned / de-provisioned."),
		Status:      redis.String("processing-completed"),
		ID:          redis.String("e9232e43-3781-4263-a38e-f4d150e03475"),
		Response: &latest_imports.Response{
			ID: redis.Int(51051302),
			Resource: &latest_imports.Resource{
				Status:         redis.String("failed"),
				LastImportTime: redis.Time(time.Date(2024, 5, 21, 10, 36, 26, 0, time.UTC)),
				FailureReason:  redis.String("file-corrupted"),
				FailureReasonParams: []*latest_imports.FailureReasonParam{
					{
						Key:   redis.String("bytes_configured_bdb_limit"),
						Value: redis.String("1234"),
					},
					{
						Key:   redis.String("bytes_of_expected_dataset"),
						Value: redis.String("5678"),
					},
				},
			},
			Error: nil,
		},
	}, actual)
}
