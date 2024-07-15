package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/maintenance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetMaintenanceAutomatic(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/113569/maintenance-windows",
				`{
				  "mode": "automatic"
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Maintenance.Get(context.TODO(), 113569)
	require.NoError(t, err)

	assert.Equal(t, &maintenance.Maintenance{
		Mode: redis.String("automatic"),
	}, actual)
}

func TestGetMaintenanceManual(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/113569/maintenance-windows",
				`{
				  "mode": "manual",
				  "timeZone": "string",
				  "windows": [
					{
					  "days": [
						"Monday"
					  ],
					  "startHour": 3,
					  "durationInHours": 8
					}
				  ],
				  "skipStatus": {
					"remainingSkips": 0,
					"currentSkipEnd": "string"
				  }
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Maintenance.Get(context.TODO(), 113569)
	require.NoError(t, err)

	assert.Equal(t, &maintenance.Maintenance{
		Mode: redis.String("manual"),
		Windows: []*maintenance.Window{
			{
				StartHour:       redis.Int(3),
				DurationInHours: redis.Int(8),
				Days:            []*string{redis.String("Monday")},
			},
		},
	}, actual)
}

func TestUpdateMaintenanceAutomatic(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey", "secret",
			putRequest(
				t,
				"/subscriptions/113569/maintenance-windows",
				`{
				  "mode": "automatic"
				}`,
				`{
				  "taskId" : "7e7b57f4-70f3-47f3-b5a4-0b3c270a9117",
				  "commandType" : "subscriptionMaintenanceWindowsUpdateRequest",
				  "status" : "received",
				  "description" : "Task request received and is being queued for processing.",
				  "timestamp" : "2024-07-15T13:25:39.304430279Z",
				  "links" : [ {
					"href" : "https://api-staging.qa.redislabs.com/v1/tasks/7e7b57f4-70f3-47f3-b5a4-0b3c270a9117",
					"rel" : "task",
					"type" : "GET"
				  } ]
				}`,
			),
			getRequest(
				t,
				"/tasks/7e7b57f4-70f3-47f3-b5a4-0b3c270a9117",
				`{
				  "taskId" : "7e7b57f4-70f3-47f3-b5a4-0b3c270a9117",
				  "commandType" : "subscriptionMaintenanceWindowsUpdateRequest",
				  "status" : "processing-completed",
				  "description" : "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp" : "2024-07-15T13:25:41.067889681Z",
				  "response" : {
					"resourceId" : 113972
				  },
				  "links" : [ {
					"href" : "https://api-staging.qa.redislabs.com/v1/subscriptions/113972",
					"rel" : "resource",
					"type" : "GET"
				  }, {
					"href" : "https://api-staging.qa.redislabs.com/v1/tasks/7e7b57f4-70f3-47f3-b5a4-0b3c270a9117",
					"rel" : "self",
					"type" : "GET"
				  } ]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.Maintenance.Update(
		context.TODO(),
		113569,
		maintenance.Maintenance{
			Mode: redis.String("automatic"),
		},
	)

	require.NoError(t, err)
}

func TestUpdateMaintenanceManual(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey", "secret",
			putRequest(
				t,
				"/subscriptions/113569/maintenance-windows",
				`{
				  "mode": "manual",
				  "windows": [
					{
					  "startHour": 12,
					  "durationInHours": 8,
					  "days": [
						"Monday",
						"Wednesday"
					  ]
					}
				  ]
				}`,
				`{
				  "taskId" : "b6e0b40f-be10-4dce-8481-f4c4812855bc",
				  "commandType" : "subscriptionMaintenanceWindowsUpdateRequest",
				  "status" : "received",
				  "description" : "Task request received and is being queued for processing.",
				  "timestamp" : "2024-07-15T13:22:29.483556324Z",
				  "links" : [ {
					"href" : "https://api-staging.qa.redislabs.com/v1/tasks/b6e0b40f-be10-4dce-8481-f4c4812855bc",
					"rel" : "task",
					"type" : "GET"
				  } ]
				}`,
			),
			getRequest(
				t,
				"/tasks/b6e0b40f-be10-4dce-8481-f4c4812855bc",
				`{
				  "taskId" : "b6e0b40f-be10-4dce-8481-f4c4812855bc",
				  "commandType" : "subscriptionMaintenanceWindowsUpdateRequest",
				  "status" : "processing-completed",
				  "description" : "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp" : "2024-07-15T13:22:32.934703954Z",
				  "response" : {
					"resourceId" : 113972
				  },
				  "links" : [ {
					"href" : "https://api-staging.qa.redislabs.com/v1/subscriptions/113972",
					"rel" : "resource",
					"type" : "GET"
				  }, {
					"href" : "https://api-staging.qa.redislabs.com/v1/tasks/b6e0b40f-be10-4dce-8481-f4c4812855bc",
					"rel" : "self",
					"type" : "GET"
				  } ]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.Maintenance.Update(
		context.TODO(),
		113569,
		maintenance.Maintenance{
			Mode: redis.String("manual"),
			Windows: []*maintenance.Window{
				{
					StartHour:       redis.Int(12),
					DurationInHours: redis.Int(8),
					Days: []*string{
						redis.String("Monday"),
						redis.String("Wednesday"),
					},
				},
			},
		},
	)

	require.NoError(t, err)
}
