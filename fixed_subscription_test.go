package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	fixedSubscriptions "github.com/RedisLabs/rediscloud-go-api/service/fixed/subscriptions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixedSubscription_Create(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey", "secret",
			postRequest(
				t,
				"/fixed/subscriptions",
				`{
					"name": "My test fixed subscription",
					"planId": 34858,
					"paymentMethodId": 30949
				}`,
				`{
					"taskId": "107ca763-da2e-4a23-9558-10153015a010",
					"commandType": "fixedSubscriptionCreateRequest",
					"status": "received",
					"description": "Task request received and is being queued for processing.",
					"timestamp": "2024-05-09T09:36:16.122289471Z",
					"links": [
						{
							"rel": "task",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/107ca763-da2e-4a23-9558-10153015a010"
						}
					]
				}`,
			),
			getRequest(
				t,
				"/tasks/107ca763-da2e-4a23-9558-10153015a010",
				`{
					"taskId": "107ca763-da2e-4a23-9558-10153015a010",
					"commandType": "fixedSubscriptionCreateRequest",
					"status": "initialized",
					"timestamp": "2020-10-28T09:58:16.798Z",
					"response": {},
					"_links": {
						"self": {
							"href": "https://example.com",
							"type": "GET"
						}
					}
				}`,
			),
			getRequest(
				t,
				"/tasks/107ca763-da2e-4a23-9558-10153015a010",
				`{
					"taskId": "107ca763-da2e-4a23-9558-10153015a010",
					"commandType": "fixedSubscriptionCreateRequest",
					"status": "processing-completed",
					"description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					"timestamp": "2024-05-09T09:36:35.177603409Z",
					"response": {
						"resourceId": 111614
					},
					"links": [
						{
							"rel": "resource",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/111614"
						},
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/107ca763-da2e-4a23-9558-10153015a010"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.FixedSubscriptions.Create(
		context.TODO(),
		fixedSubscriptions.FixedSubscriptionRequest{
			Name:            redis.String("My test fixed subscription"),
			PlanId:          redis.Int(34858),
			PaymentMethodID: redis.Int(30949),
		},
	)

	require.NoError(t, err)
	assert.Equal(t, 111614, actual)
}

func TestFixedSubscription_Create_Marketplace(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey", "secret",
			postRequest(
				t,
				"/fixed/subscriptions",
				`{
					"name": "My test fixed subscription with marketplace payments",
					"planId": 34811,
					"paymentMethod": "marketplace"
				}`,
				`{
					"taskId": "2a6c6c5b-a16a-4f19-a803-17c1013a5888",
					"commandType": "fixedSubscriptionCreateRequest",
					"status": "received",
					"description": "Task request received and is being queued for processing.",
					"timestamp": "2024-05-09T09:36:16.122289471Z",
					"links": [
						{
							"rel": "task",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/2a6c6c5b-a16a-4f19-a803-17c1013a5888"
						}
					]
				}`,
			),
			getRequest(
				t,
				"/tasks/2a6c6c5b-a16a-4f19-a803-17c1013a5888",
				`{
					"taskId": "2a6c6c5b-a16a-4f19-a803-17c1013a5888",
					"commandType": "fixedSubscriptionCreateRequest",
					"status": "initialized",
					"timestamp": "2020-10-28T09:58:16.798Z",
					"response": {},
					"_links": {
						"self": {
							"href": "https://example.com",
							"type": "GET"
						}
					}
				}`,
			),
			getRequest(
				t,
				"/tasks/2a6c6c5b-a16a-4f19-a803-17c1013a5888",
				`{
					"taskId": "2a6c6c5b-a16a-4f19-a803-17c1013a5888",
					"commandType": "fixedSubscriptionCreateRequest",
					"status": "processing-completed",
					"description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					"timestamp": "2024-05-09T09:36:35.177603409Z",
					"response": {
						"resourceId": 111191
					},
					"links": [
						{
							"rel": "resource",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/111191"
						},
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/2a6c6c5b-a16a-4f19-a803-17c1013a5888"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.FixedSubscriptions.Create(
		context.TODO(),
		fixedSubscriptions.FixedSubscriptionRequest{
			Name:          redis.String("My test fixed subscription with marketplace payments"),
			PlanId:        redis.Int(34811),
			PaymentMethod: redis.String("marketplace"),
		},
	)

	require.NoError(t, err)
	assert.Equal(t, 111191, actual)
}

func TestFixedSubscription_List(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			getRequest(
				t,
				"/fixed/subscriptions",
				`{
					"accountId": 69369,
					"subscriptions": [
						{
							"id": 111614,
							"name": "My test fixed subscription",
							"status": "active",
							"paymentMethodId": 30949,
							"paymentMethodType": "credit-card",
							"planId": 34858,
							"planName": "250MB",
							"size": 250,
							"sizeMeasurementUnit": "MB",
							"provider": "AWS",
							"region": "us-east-1",
							"price": 5,
							"pricePeriod": "Month",
							"priceCurrency": "USD",
							"maximumDatabases": 1,
							"availability": "No replication",
							"connections": "256",
							"cidrAllowRules": 4,
							"supportDataPersistence": false,
							"supportInstantAndDailyBackups": true,
							"supportReplication": false,
							"supportClustering": false,
							"customerSupport": "Standard",
							"creationDate": "2024-05-09T09:36:18Z",
							"links": []
						},
						{
							"id": 111615,
							"name": "Another test fixed subscription",
							"status": "active",
							"paymentMethodId": 30949,
							"paymentMethodType": "credit-card",
							"planId": 34858,
							"planName": "250MB",
							"size": 250,
							"sizeMeasurementUnit": "MB",
							"provider": "AWS",
							"region": "us-east-1",
							"price": 5,
							"pricePeriod": "Month",
							"priceCurrency": "USD",
							"maximumDatabases": 1,
							"availability": "No replication",
							"connections": "256",
							"cidrAllowRules": 4,
							"supportDataPersistence": false,
							"supportInstantAndDailyBackups": true,
							"supportReplication": false,
							"supportClustering": false,
							"customerSupport": "Standard",
							"creationDate": "2024-05-09T10:49:52Z",
							"links": []
						}
					],
					"links": [
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.FixedSubscriptions.List(context.TODO())
	require.NoError(t, err)

	assert.ElementsMatch(t, []*fixedSubscriptions.FixedSubscriptionResponse{
		{
			ID:              redis.Int(111614),
			Name:            redis.String("My test fixed subscription"),
			Status:          redis.String("active"),
			PlanId:          redis.Int(34858),
			PaymentMethodID: redis.Int(30949),
			PaymentMethod:   redis.String("credit-card"),
			CreationDate:    redis.Time(time.Date(2024, 5, 9, 9, 36, 18, 0, time.UTC)),
		},
		{
			ID:              redis.Int(111615),
			Name:            redis.String("Another test fixed subscription"),
			Status:          redis.String("active"),
			PlanId:          redis.Int(34858),
			PaymentMethodID: redis.Int(30949),
			PaymentMethod:   redis.String("credit-card"),
			CreationDate:    redis.Time(time.Date(2024, 5, 9, 10, 49, 52, 0, time.UTC)),
		},
	}, actual)
}

func TestFixedSubscription_Get(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			getRequest(
				t,
				"/fixed/subscriptions/111614",
				`{
					"id": 111614,
					"name": "My test fixed subscription",
					"status": "active",
					"paymentMethodId": 30949,
					"paymentMethodType": "credit-card",
					"planId": 34858,
					"planName": "250MB",
					"size": 250,
					"sizeMeasurementUnit": "MB",
					"provider": "AWS",
					"region": "us-east-1",
					"price": 5,
					"pricePeriod": "Month",
					"priceCurrency": "USD",
					"maximumDatabases": 1,
					"availability": "No replication",
					"connections": "256",
					"cidrAllowRules": 4,
					"supportDataPersistence": false,
					"supportInstantAndDailyBackups": true,
					"supportReplication": false,
					"supportClustering": false,
					"customerSupport": "Standard",
					"creationDate": "2024-05-09T09:36:18Z",
					"links": [
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/111614"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.FixedSubscriptions.Get(context.TODO(), 111614)
	require.NoError(t, err)

	assert.Equal(t, &fixedSubscriptions.FixedSubscriptionResponse{
		ID:              redis.Int(111614),
		Name:            redis.String("My test fixed subscription"),
		Status:          redis.String("active"),
		PlanId:          redis.Int(34858),
		PaymentMethod:   redis.String("credit-card"),
		PaymentMethodID: redis.Int(30949),
		CreationDate:    redis.Time(time.Date(2024, 5, 9, 9, 36, 18, 0, time.UTC)),
	}, actual)
}

func TestFixedSubscription_Get_wraps404Error(t *testing.T) {
	s := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			getRequestWithStatus(
				t,
				"/fixed/subscriptions/123",
				404,
				"",
			),
		),
	)

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.FixedSubscriptions.Get(context.TODO(), 123)

	assert.Nil(t, actual)
	assert.IsType(t, &fixedSubscriptions.NotFound{}, err)
}

func TestFixedSubscription_Update(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey", "secret",
			putRequest(
				t,
				"/fixed/subscriptions/111614",
				`{
					"name": "My renamed fixed subscription",
					"planId": 34853
				}`,
				`{
					"taskId": "99e15c6e-cf4f-46d4-99b5-5eae9e7b9d9c",
					"commandType": "fixedSubscriptionUpdateRequest",
					"status": "received",
					"description": "Task request received and is being queued for processing.",
					"timestamp": "2024-05-09T10:08:20.912550659Z",
					"links": [
						{
							"rel": "task",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/99e15c6e-cf4f-46d4-99b5-5eae9e7b9d9c"
						}
					]
				}`,
			),
			getRequest(
				t,
				"/tasks/99e15c6e-cf4f-46d4-99b5-5eae9e7b9d9c",
				`{
						"taskId": "99e15c6e-cf4f-46d4-99b5-5eae9e7b9d9c",
						"commandType": "fixedSubscriptionUpdateRequest",
						"status": "processing-completed",
						"description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
						"timestamp": "2024-05-09T10:08:24.54464946Z",
						"response": {
						"resourceId": 111614
					},
					"links": [
						{
							"rel": "resource",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/fixed/subscriptions/111614"
						},
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/99e15c6e-cf4f-46d4-99b5-5eae9e7b9d9c"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.FixedSubscriptions.Update(
		context.TODO(),
		111614,
		fixedSubscriptions.FixedSubscriptionRequest{
			Name:   redis.String("My renamed fixed subscription"),
			PlanId: redis.Int(34853),
		},
	)
	require.NoError(t, err)
}

func TestFixedSubscription_Delete(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"apiKey",
			"secret",
			deleteRequest(
				t,
				"/fixed/subscriptions/111614",
				`{
					"taskId": "819d87c3-9a9f-4813-ba31-9d9938397541",
					"commandType": "fixedSubscriptionDeleteRequest",
					"status": "received",
					"description": "Task request received and is being queued for processing.",
					"timestamp": "2024-05-09T10:17:34.937624078Z",
					"links": [
						{
							"rel": "task",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/819d87c3-9a9f-4813-ba31-9d9938397541"
						}
					]
				}`,
			),
			getRequest(
				t,
				"/tasks/819d87c3-9a9f-4813-ba31-9d9938397541",
				`{
					"taskId": "819d87c3-9a9f-4813-ba31-9d9938397541",
					"commandType": "fixedSubscriptionDeleteRequest",
					"status": "processing-completed",
					"description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					"timestamp": "2024-05-09T10:17:48.619671336Z",
					"response": {
					"resourceId": 111614
					},
					"links": [
						{
							"rel": "self",
							"type": "GET",
							"href": "https://api-staging.qa.redislabs.com/v1/tasks/819d87c3-9a9f-4813-ba31-9d9938397541"
						}
					]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.FixedSubscriptions.Delete(context.TODO(), 111614)
	require.NoError(t, err)
}
