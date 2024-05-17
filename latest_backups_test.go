package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLatestBackup(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/12/databases/34/backup",
				`{
				  "taskId": "50ec6172-8475-4ef6-8b3c-d61e688d8fe5",
				  "commandType": "databaseBackupStatusRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-04-15T09:08:04.222268Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/50ec6172-8475-4ef6-8b3c-d61e688d8fe5",
					  "type": "GET",
					  "rel": "task"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/50ec6172-8475-4ef6-8b3c-d61e688d8fe5",
				`{
				  "taskId": "50ec6172-8475-4ef6-8b3c-d61e688d8fe5",
				  "commandType": "databaseBackupStatusRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-04-15T09:08:07.537915Z",
				  "response": {
					"resourceId": 51051292,
					"additionalResourceId": 12,
					"resource": {}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/50ec6172-8475-4ef6-8b3c-d61e688d8fe5",
					  "type": "GET",
					  "rel": "self"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	_, err = subject.LatestBackup.Get(context.TODO(), 12, 34)
	require.NoError(t, err)
}

func TestGetFixedLatestBackup(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/fixed/subscriptions/12/databases/34/backup",
				`{
				  "taskId": "ce2cbfea-9b15-4250-a516-f014161a8dd3",
				  "commandType": "databaseBackupStatusRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-04-15T09:52:23.963337Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/ce2cbfea-9b15-4250-a516-f014161a8dd3",
					  "type": "GET",
					  "rel": "task"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/ce2cbfea-9b15-4250-a516-f014161a8dd3",
				`{
				  "taskId": "ce2cbfea-9b15-4250-a516-f014161a8dd3",
				  "commandType": "databaseBackupStatusRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2024-04-15T09:52:26.101936Z",
				  "response": {
					"error": {
					  "type": "DATABASE_BACKUP_DISABLED",
					  "status": "400 BAD_REQUEST",
					  "description": "Database backup is disabled"
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/ce2cbfea-9b15-4250-a516-f014161a8dd3",
					  "type": "GET",
					  "rel": "self"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	_, err = subject.LatestBackup.GetFixed(context.TODO(), 12, 34)
	require.NoError(t, err)
}

func TestGetAALatestBackup(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/12/databases/34/backup?regionName=eu-west-2",
				`{
				  "taskId": "ce2cbfea-9b15-4250-a516-f014161a8dd3",
				  "commandType": "databaseBackupStatusRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-04-15T09:52:23.963337Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/ce2cbfea-9b15-4250-a516-f014161a8dd3",
					  "type": "GET",
					  "rel": "task"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/ce2cbfea-9b15-4250-a516-f014161a8dd3",
				`{
				  "taskId": "ce2cbfea-9b15-4250-a516-f014161a8dd3",
				  "commandType": "databaseBackupStatusRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2024-04-15T09:52:26.101936Z",
				  "response": {
					"error": {
					  "type": "DATABASE_BACKUP_DISABLED",
					  "status": "400 BAD_REQUEST",
					  "description": "Database backup is disabled"
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/ce2cbfea-9b15-4250-a516-f014161a8dd3",
					  "type": "GET",
					  "rel": "self"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	_, err = subject.LatestBackup.GetActiveActive(context.TODO(), 12, 34, "eu-west-2")
	require.NoError(t, err)
}
