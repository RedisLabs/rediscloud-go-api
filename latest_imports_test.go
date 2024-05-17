package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

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

	_, err = subject.LatestImport.Get(context.TODO(), 12, 34)
	require.NoError(t, err)
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
					"resource": {}
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

	_, err = subject.LatestImport.GetFixed(context.TODO(), 12, 34)
	require.NoError(t, err)
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
					"resource": {}
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

	_, err = subject.LatestImport.Get(context.TODO(), 12, 34)
	require.NoError(t, err)
}
