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

// TODO Replace with a working response!
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
				  "taskId": "1f19cc39-5346-4629-8dc7-f0dc9ad63720",
				  "commandType": "subscriptionMaintenanceWindowsUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-08T09:13:21.84935915Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/1f19cc39-5346-4629-8dc7-f0dc9ad63720"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/1f19cc39-5346-4629-8dc7-f0dc9ad63720",
				`{
				  "taskId": "1f19cc39-5346-4629-8dc7-f0dc9ad63720",
				  "commandType": "subscriptionMaintenanceWindowsUpdateRequest",
				  "status": "processing-completed",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2024-07-08T09:13:22.748959662Z",
				  "response": {
				  },
				  "links": [
					{
					  "rel": "self",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/1f19cc39-5346-4629-8dc7-f0dc9ad63720"
					}
				  ]
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

// TODO Replace with a working response!
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
				  "taskId": "1f19cc39-5346-4629-8dc7-f0dc9ad63720",
				  "commandType": "subscriptionMaintenanceWindowsUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-08T09:13:21.84935915Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/1f19cc39-5346-4629-8dc7-f0dc9ad63720"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/1f19cc39-5346-4629-8dc7-f0dc9ad63720",
				`{
				  "taskId": "1f19cc39-5346-4629-8dc7-f0dc9ad63720",
				  "commandType": "subscriptionMaintenanceWindowsUpdateRequest",
				  "status": "processing-completed",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2024-07-08T09:13:22.748959662Z",
				  "response": {
				  },
				  "links": [
					{
					  "rel": "self",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/1f19cc39-5346-4629-8dc7-f0dc9ad63720"
					}
				  ]
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
