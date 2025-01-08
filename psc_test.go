package rediscloud_api

import (
	"context"
	"errors"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/psc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPrivateServiceConnectService(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *psc.PrivateServiceConnectService
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a private service connect service",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/private-service-connect",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscServiceGetRequest",
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
				  "commandType": "pscServiceGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "id": 40,
					  "connectionHostName": "psc.mc2018-0.us-central1-mz.gcp.sdk-cloud.rlrcp.com",
					  "serviceAttachmentName": "service-attachment-mc2018-0-us-central1-mz-rlrcp",
					  "status": "active"
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
			expectedResult: &psc.PrivateServiceConnectService{
				ID:                    redis.Int(40),
				ConnectionHostName:    redis.String("psc.mc2018-0.us-central1-mz.gcp.sdk-cloud.rlrcp.com"),
				ServiceAttachmentName: redis.String("service-attachment-mc2018-0-us-central1-mz-rlrcp"),
				Status:                redis.String("active"),
			},
		},
		{
			description: "should fail when private service connect is not found",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/private-service-connect",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscServiceGetRequest",
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
				  "commandType": "pscServiceGetRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2025-01-13T11:22:51.204189721Z",
				  "response": {
					"error": {
					  "type": "PSC_SERVICE_NOT_FOUND",
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
			expectedErrorAs: &psc.NotFound{},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithStatus(
					t,
					"/subscriptions/114019/private-service-connect",
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/private-service-connect"
					}`),
			},
			expectedError:   errors.New("resource not found - subscription 114019"),
			expectedErrorAs: &psc.NotFound{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateServiceConnect.GetService(context.TODO(), 114019)
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

func TestGetPrivateServiceConnectServiceActiveActive(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *psc.PrivateServiceConnectService
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a private service connect service",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/regions/1/private-service-connect",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscServiceGetRequest",
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
				  "commandType": "activeActivePscServiceGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "id": 40,
					  "connectionHostName": "psc.mc2018-0.us-central1-mz.gcp.sdk-cloud.rlrcp.com",
					  "serviceAttachmentName": "service-attachment-mc2018-0-us-central1-mz-rlrcp",
					  "status": "active"
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
			expectedResult: &psc.PrivateServiceConnectService{
				ID:                    redis.Int(40),
				ConnectionHostName:    redis.String("psc.mc2018-0.us-central1-mz.gcp.sdk-cloud.rlrcp.com"),
				ServiceAttachmentName: redis.String("service-attachment-mc2018-0-us-central1-mz-rlrcp"),
				Status:                redis.String("active"),
			},
		},
		{
			description: "should fail when private service connect is not found",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/regions/1/private-service-connect",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscServiceGetRequest",
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
				  "commandType": "activeActivePscServiceGetRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2025-01-13T11:22:51.204189721Z",
				  "response": {
					"error": {
					  "type": "PSC_SERVICE_NOT_FOUND",
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
			expectedError:   errors.New("resource not found - subscription 114019 and region 1"),
			expectedErrorAs: &psc.NotFoundActiveActive{},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithStatus(
					t,
					"/subscriptions/114019/regions/1/private-service-connect",
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/regions/1/private-service-connect"
					}`),
			},
			expectedError:   errors.New("resource not found - subscription 114019 and region 1"),
			expectedErrorAs: &psc.NotFoundActiveActive{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateServiceConnect.GetActiveActiveService(context.TODO(), 114019, 1)
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

func TestCreatePrivateServiceConnectService(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequestWithNoRequest(
				t,
				"/subscriptions/114019/private-service-connect",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "pscServiceCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-01-13T11:58:49.673306547Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "pscServiceCreateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 40
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.PrivateServiceConnect.CreateService(context.TODO(), 114019)
	assert.NoError(t, err)
	assert.Equal(t, 40, actual)
}

func TestCreatePrivateServiceConnectEndpoint(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequest(
				t,
				"/subscriptions/114019/private-service-connect/40",
				`{
				  "gcpProjectId": "my-gcp-project",
				  "gcpVpcName": "my-vpc",
				  "gcpVpcSubnetName": "my-vpc-subnet",
				  "endpointConnectionName": "my-endpoint-connection"
				}`,
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "pscEndpointCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-01-13T11:58:49.673306547Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "pscEndpointCreateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 39
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.PrivateServiceConnect.CreateEndpoint(context.TODO(), 114019, 40, psc.CreatePrivateServiceConnectEndpoint{
		GCPProjectID:           redis.String("my-gcp-project"),
		GCPVPCName:             redis.String("my-vpc"),
		GCPVPCSubnetName:       redis.String("my-vpc-subnet"),
		EndpointConnectionName: redis.String("my-endpoint-connection"),
	})
	assert.NoError(t, err)
	assert.Equal(t, 39, actual)
}

func TestCreatePrivateServiceConnectEndpointActiveActive(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequest(
				t,
				"/subscriptions/114019/regions/1/private-service-connect/40",
				`{
				  "gcpProjectId": "my-gcp-project",
				  "gcpVpcName": "my-vpc",
				  "gcpVpcSubnetName": "my-vpc-subnet",
				  "endpointConnectionName": "my-endpoint-connection"
				}`,
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "activeActivePscEndpointCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-01-13T11:58:49.673306547Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "activeActivePscEndpointCreateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 39
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.PrivateServiceConnect.CreateActiveActiveEndpoint(context.TODO(), 114019, 1, 40, psc.CreatePrivateServiceConnectEndpoint{
		GCPProjectID:           redis.String("my-gcp-project"),
		GCPVPCName:             redis.String("my-vpc"),
		GCPVPCSubnetName:       redis.String("my-vpc-subnet"),
		EndpointConnectionName: redis.String("my-endpoint-connection"),
	})
	assert.NoError(t, err)
	assert.Equal(t, 39, actual)
}

func TestGetPrivateServiceConnectEndpoints(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *psc.PrivateServiceConnectEndpoints
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a private service connect service endpoints",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/private-service-connect/40",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscServiceEndpointsGetRequest",
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
				}`),
				getRequest(
					t,
					"/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
					`{
					  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
					  "commandType": "pscServiceEndpointsGetRequest",
					  "status": "processing-completed",
					  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					  "timestamp": "2025-01-14T16:18:42.048633704Z",
					  "response": {
						"resourceId": 114019,
						"resource": {
						  "pscServiceId": 40,
						  "endpoints": [
							{
							  "id": 39,
							  "gcpProjectId": "my-gcp-project",
							  "gcpVpcName": "my-vpc",
							  "gcpVpcSubnetName": "my-vpc-subnet",
							  "endpointConnectionName": "my-endpoint-connection",
							  "status": "initialized"
							}
						  ]
						}
					  },
					  "links": [
						{
						  "rel": "self",
						  "type": "GET",
						  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971"
						}
					  ]
					}`),
			},
			expectedResult: &psc.PrivateServiceConnectEndpoints{
				PSCServiceID: redis.Int(40),
				Endpoints: []*psc.PrivateServiceConnectEndpoint{
					{
						ID:                     redis.Int(39),
						GCPProjectID:           redis.String("my-gcp-project"),
						GCPVPCName:             redis.String("my-vpc"),
						GCPVPCSubnetName:       redis.String("my-vpc-subnet"),
						EndpointConnectionName: redis.String("my-endpoint-connection"),
						Status:                 redis.String("initialized"),
					},
				},
			},
		},
		{
			description: "should fail when private service connect is not found",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/private-service-connect/40",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscServiceEndpointsGetRequest",
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
				  "commandType": "pscServiceEndpointsGetRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2025-01-13T11:22:51.204189721Z",
				  "response": {
					"error": {
					  "type": "PSC_SERVICE_NOT_FOUND",
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
			expectedErrorAs: &psc.NotFound{},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithStatus(
					t,
					"/subscriptions/114019/private-service-connect/40",
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/private-service-connect/40"
					}`,
				),
			},
			expectedError:   errors.New("resource not found - subscription 114019"),
			expectedErrorAs: &psc.NotFound{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateServiceConnect.GetEndpoints(context.TODO(), 114019, 40)
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

func TestGetPrivateServiceConnectEndpointsActiveActive(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *psc.PrivateServiceConnectEndpoints
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a private service connect service endpoints",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/regions/1/private-service-connect/40",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscServiceEndpointsGetRequest",
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
				}`),
				getRequest(
					t,
					"/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
					`{
					  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
					  "commandType": "activeActivePscServiceEndpointsGetRequest",
					  "status": "processing-completed",
					  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					  "timestamp": "2025-01-14T16:18:42.048633704Z",
					  "response": {
						"resourceId": 114019,
						"resource": {
						  "pscServiceId": 40,
						  "endpoints": [
							{
							  "id": 39,
							  "gcpProjectId": "my-gcp-project",
							  "gcpVpcName": "my-vpc",
							  "gcpVpcSubnetName": "my-vpc-subnet",
							  "endpointConnectionName": "my-endpoint-connection",
							  "status": "initialized"
							}
						  ]
						}
					  },
					  "links": [
						{
						  "rel": "self",
						  "type": "GET",
						  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971"
						}
					  ]
					}`),
			},
			expectedResult: &psc.PrivateServiceConnectEndpoints{
				PSCServiceID: redis.Int(40),
				Endpoints: []*psc.PrivateServiceConnectEndpoint{
					{
						ID:                     redis.Int(39),
						GCPProjectID:           redis.String("my-gcp-project"),
						GCPVPCName:             redis.String("my-vpc"),
						GCPVPCSubnetName:       redis.String("my-vpc-subnet"),
						EndpointConnectionName: redis.String("my-endpoint-connection"),
						Status:                 redis.String("initialized"),
					},
				},
			},
		},
		{
			description: "should fail when private service connect is not found",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/regions/1/private-service-connect/40",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscServiceEndpointsGetRequest",
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
				  "commandType": "activeActivePscServiceEndpointsGetRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2025-01-13T11:22:51.204189721Z",
				  "response": {
					"error": {
					  "type": "PSC_SERVICE_NOT_FOUND",
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
			expectedError:   errors.New("resource not found - subscription 114019 and region 1"),
			expectedErrorAs: &psc.NotFoundActiveActive{},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithStatus(
					t,
					"/subscriptions/114019/regions/1/private-service-connect/40",
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/regions/1/private-service-connect/40"
					}`,
				),
			},
			expectedError:   errors.New("resource not found - subscription 114019 and region 1"),
			expectedErrorAs: &psc.NotFoundActiveActive{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateServiceConnect.GetActiveActiveEndpoints(context.TODO(), 114019, 1, 40)
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

func TestGetPrivateServiceConnectEndpointCreationScript(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *psc.CreationScript
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a private service connect service endpoint creation script",
			mockedResponse: []endpointRequest{
				getRequestWithQuery(
					t,
					"/subscriptions/114019/private-service-connect/40/endpoints/52/creationScripts",
					url.Values{
						"includeTerraformGcpScript": []string{"true"},
					},
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscServiceEndpointScriptGetRequest",
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
				  "commandType": "pscServiceEndpointScriptGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2025-01-14T16:18:42.048633704Z",
				  "response": {
					"resourceId": 42,
					"resource": {
					  "script": {
						"bash": "bash script",
						"powershell": "powershell script",
						"terraformGcp": {
						  "serviceAttachments": [
							{
							  "name": "projects/s169ae38adc5a3c69/regions/us-central1/serviceAttachments/service-attachment1-mc2025-0-us-central1-mz-rlrcp",
							  "dnsRecord": "psc1.mc2025-0.us-central1-mz.gcp.sdk-cloud.rlrcp.com.",
							  "ipAddressName": "redis-1111-psc-static-ip1",
							  "forwardingRuleName": "redis-1111"
							}
						  ]
						}
					  }
					}
				  },
				  "links": [
					{
					  "rel": "self",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971"
					}
				  ]
				}`,
				),
			},
			expectedResult: &psc.CreationScript{
				Script: &psc.GCPCreationScript{
					Bash:       redis.String("bash script"),
					Powershell: redis.String("powershell script"),
					TerraformGcp: &psc.TerraformGCP{
						ServiceAttachments: []psc.TerraformGCPServiceAttachment{
							{
								Name:               redis.String("projects/s169ae38adc5a3c69/regions/us-central1/serviceAttachments/service-attachment1-mc2025-0-us-central1-mz-rlrcp"),
								DNSRecord:          redis.String("psc1.mc2025-0.us-central1-mz.gcp.sdk-cloud.rlrcp.com."),
								IPAddressName:      redis.String("redis-1111-psc-static-ip1"),
								ForwardingRuleName: redis.String("redis-1111"),
							},
						},
					},
				},
			},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithQueryAndStatus(
					t,
					"/subscriptions/114019/private-service-connect/40/endpoints/52/creationScripts",
					url.Values{
						"includeTerraformGcpScript": []string{"true"},
					},
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/private-service-connect/40/endpoints/52/creationScripts"
					}`,
				),
			},
			expectedError:   errors.New("resource not found - subscription 114019"),
			expectedErrorAs: &psc.NotFound{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateServiceConnect.GetEndpointCreationScripts(context.TODO(),
				114019, 40, 52, true)
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

func TestGetPrivateServiceConnectEndpointCreationScriptActiveActive(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *psc.CreationScript
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a private service connect service endpoint creation script",
			mockedResponse: []endpointRequest{
				getRequestWithQuery(
					t,
					"/subscriptions/114019/regions/1/private-service-connect/40/endpoints/52/creationScripts",
					url.Values{
						"includeTerraformGcpScript": []string{"true"},
					},
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscServiceEndpointScriptGetRequest",
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
				  "commandType": "activeActivePscServiceEndpointScriptGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2025-01-14T16:18:42.048633704Z",
				  "response": {
					"resourceId": 42,
					"resource": {
					  "script": {
						"bash": "bash script",
						"powershell": "powershell script",
						"terraformGcp": {
						  "serviceAttachments": [
							{
							  "name": "projects/s169ae38adc5a3c69/regions/us-central1/serviceAttachments/service-attachment1-mc2025-0-us-central1-mz-rlrcp",
							  "dnsRecord": "psc1.mc2025-0.us-central1-mz.gcp.sdk-cloud.rlrcp.com.",
							  "ipAddressName": "redis-1111-psc-static-ip1",
							  "forwardingRuleName": "redis-1111"
							}
						  ]
						}
					  }
					}
				  },
				  "links": [
					{
					  "rel": "self",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971"
					}
				  ]
				}`,
				),
			},
			expectedResult: &psc.CreationScript{
				Script: &psc.GCPCreationScript{
					Bash:       redis.String("bash script"),
					Powershell: redis.String("powershell script"),
					TerraformGcp: &psc.TerraformGCP{
						ServiceAttachments: []psc.TerraformGCPServiceAttachment{
							{
								Name:               redis.String("projects/s169ae38adc5a3c69/regions/us-central1/serviceAttachments/service-attachment1-mc2025-0-us-central1-mz-rlrcp"),
								DNSRecord:          redis.String("psc1.mc2025-0.us-central1-mz.gcp.sdk-cloud.rlrcp.com."),
								IPAddressName:      redis.String("redis-1111-psc-static-ip1"),
								ForwardingRuleName: redis.String("redis-1111"),
							},
						},
					},
				},
			},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithQueryAndStatus(
					t,
					"/subscriptions/114019/regions/1/private-service-connect/40/endpoints/52/creationScripts",
					url.Values{
						"includeTerraformGcpScript": []string{"true"},
					},
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/regions/1/private-service-connect/40/endpoints/52/creationScripts"
					}`,
				),
			},
			expectedError:   errors.New("resource not found - subscription 114019 and region 1"),
			expectedErrorAs: &psc.NotFoundActiveActive{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateServiceConnect.GetActiveActiveEndpointCreationScripts(context.TODO(),
				114019, 1, 40, 52, true)
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

func TestGetPrivateServiceConnectEndpointDeletionScript(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *psc.DeletionScript
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a private service connect service endpoint deletion script",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/private-service-connect/40/endpoints/52/deletionScripts",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscServiceEndpointScriptGetRequest",
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
				  "commandType": "pscServiceEndpointScriptGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2025-01-14T16:18:42.048633704Z",
				  "response": {
					"resourceId": 42,
					"resource": {
					  "script": {
						"bash": "bash script",
						"powershell": "powershell script"
					  }
					}
				  },
				  "links": [
					{
					  "rel": "self",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971"
					}
				  ]
				}`,
				),
			},
			expectedResult: &psc.DeletionScript{
				Script: &psc.GCPDeletionScript{
					Bash:       redis.String("bash script"),
					Powershell: redis.String("powershell script"),
				},
			},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithStatus(
					t,
					"/subscriptions/114019/private-service-connect/40/endpoints/52/deletionScripts",
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/private-service-connect/40/endpoints/52/deletionScripts"
					}`,
				),
			},
			expectedError:   errors.New("resource not found - subscription 114019"),
			expectedErrorAs: &psc.NotFound{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateServiceConnect.GetEndpointDeletionScripts(context.TODO(),
				114019, 40, 52)
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

func TestGetPrivateServiceConnectEndpointDeletionScriptActiveActive(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *psc.DeletionScript
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a private service connect service endpoint deletion script",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/regions/1/private-service-connect/40/endpoints/52/deletionScripts",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscServiceEndpointScriptGetRequest",
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
				  "commandType": "activeActivePscServiceEndpointScriptGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2025-01-14T16:18:42.048633704Z",
				  "response": {
					"resourceId": 42,
					"resource": {
					  "script": {
						"bash": "bash script",
						"powershell": "powershell script"
					  }
					}
				  },
				  "links": [
					{
					  "rel": "self",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/502fc31f-fd44-4cb0-a429-07882309a971"
					}
				  ]
				}`,
				),
			},
			expectedResult: &psc.DeletionScript{
				Script: &psc.GCPDeletionScript{
					Bash:       redis.String("bash script"),
					Powershell: redis.String("powershell script"),
				},
			},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithStatus(
					t,
					"/subscriptions/114019/regions/1/private-service-connect/40/endpoints/52/deletionScripts",
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/regions/1/private-service-connect/40/endpoints/52/deletionScripts"
					}`,
				),
			},
			expectedError:   errors.New("resource not found - subscription 114019 and region 1"),
			expectedErrorAs: &psc.NotFoundActiveActive{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateServiceConnect.GetActiveActiveEndpointDeletionScripts(context.TODO(),
				114019, 1, 40, 52)
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

func TestCreatePrivateServiceConnectServiceActiveActive(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequestWithNoRequest(
				t,
				"/subscriptions/114019/regions/1/private-service-connect",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "activeActivePscServiceCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-01-13T11:58:49.673306547Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "activeActivePscServiceGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 40
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.PrivateServiceConnect.CreateActiveActiveService(context.TODO(), 114019, 1)
	assert.NoError(t, err)
	assert.Equal(t, 40, actual)
}

func TestUpdatePrivateServiceConnectEndpoint(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/subscriptions/114019/private-service-connect/40/endpoints/50",
				`{
				  "gcpProjectId": "my-gcp-project",
				  "gcpVpcName": "my-vpc",
				  "gcpVpcSubnetName": "my-vpc-subnet",
				  "endpointConnectionName": "my-endpoint-connection",
				  "action":"accept"
				}`,
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "pscEndpointUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-01-13T11:58:49.673306547Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "pscEndpointUpdateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 39
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.PrivateServiceConnect.UpdateEndpoint(context.TODO(), 114019, 40, 50, &psc.UpdatePrivateServiceConnectEndpoint{
		GCPProjectID:           redis.String("my-gcp-project"),
		GCPVPCName:             redis.String("my-vpc"),
		GCPVPCSubnetName:       redis.String("my-vpc-subnet"),
		EndpointConnectionName: redis.String("my-endpoint-connection"),
		Action:                 redis.String(psc.EndpointActionAccept),
	})
	assert.NoError(t, err)
}

func TestUpdatePrivateServiceConnectEndpointActiveActive(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/subscriptions/114019/regions/1/private-service-connect/40/endpoints/50",
				`{
				  "gcpProjectId": "my-gcp-project",
				  "gcpVpcName": "my-vpc",
				  "gcpVpcSubnetName": "my-vpc-subnet",
				  "endpointConnectionName": "my-endpoint-connection",
				  "action":"accept"
				}`,
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "activeActivePscEndpointUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-01-13T11:58:49.673306547Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "activeActivePscEndpointUpdateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 39
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.PrivateServiceConnect.UpdateActiveActiveEndpoint(context.TODO(), 114019, 1, 40, 50, &psc.UpdatePrivateServiceConnectEndpoint{
		GCPProjectID:           redis.String("my-gcp-project"),
		GCPVPCName:             redis.String("my-vpc"),
		GCPVPCSubnetName:       redis.String("my-vpc-subnet"),
		EndpointConnectionName: redis.String("my-endpoint-connection"),
		Action:                 redis.String(psc.EndpointActionAccept),
	})
	assert.NoError(t, err)
}

func TestUpdatePrivateServiceConnectEndpointAccept(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/subscriptions/114019/private-service-connect/40/endpoints/50",
				`{
				  "action":"accept"
				}`,
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "pscEndpointUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-01-13T11:58:49.673306547Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "pscEndpointUpdateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 39
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.PrivateServiceConnect.UpdateEndpoint(context.TODO(), 114019, 40, 50, &psc.UpdatePrivateServiceConnectEndpoint{
		Action: redis.String(psc.EndpointActionAccept),
	})
	assert.NoError(t, err)
}

func TestUpdatePrivateServiceConnectEndpointActiveActiveAccept(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/subscriptions/114019/regions/1/private-service-connect/40/endpoints/50",
				`{
				  "action":"accept"
				}`,
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "activeActivePscEndpointUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-01-13T11:58:49.673306547Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/41a37bac-91dc-4468-bb4b-45dfd61df515",
				`{
				  "taskId": "41a37bac-91dc-4468-bb4b-45dfd61df515",
				  "commandType": "activeActivePscEndpointUpdateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 39
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.PrivateServiceConnect.UpdateActiveActiveEndpoint(context.TODO(), 114019, 1, 40, 50, &psc.UpdatePrivateServiceConnectEndpoint{
		Action: redis.String(psc.EndpointActionAccept),
	})
	assert.NoError(t, err)
}

func TestDeletePrivateServiceConnectService(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully delete a private service connect service",
			mockedResponse: []endpointRequest{
				deleteRequest(
					t,
					"/subscriptions/114019/private-service-connect",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscServiceDeleteRequest",
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
				  "commandType": "pscServiceDeleteRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019
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
		},
		{
			description: "should fail when private service connect is not found",
			mockedResponse: []endpointRequest{
				deleteRequestWithStatus(
					t,
					"/subscriptions/114019/private-service-connect",
					404,
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscServiceDeleteRequest",
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
				  "commandType": "pscServiceDeleteRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2025-01-13T11:22:51.204189721Z",
				  "response": {
					"error": {
					  "type": "PSC_SERVICE_NOT_FOUND",
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
			expectedErrorAs: &psc.NotFound{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			err = subject.PrivateServiceConnect.DeleteService(context.TODO(), 114019)
			if testCase.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.IsType(t, err, testCase.expectedErrorAs)
				assert.EqualError(t, err, testCase.expectedError.Error())
			}
		})
	}
}

func TestDeletePrivateServiceConnectServiceActiveActive(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully delete a private service connect service",
			mockedResponse: []endpointRequest{
				deleteRequest(
					t,
					"/subscriptions/114019/regions/1/private-service-connect",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscServiceDeleteRequest",
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
				  "commandType": "activeActivePscServiceDeleteRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019
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
		},
		{
			description: "should fail when private service connect is not found",
			mockedResponse: []endpointRequest{
				deleteRequestWithStatus(
					t,
					"/subscriptions/114019/regions/1/private-service-connect",
					404,
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscServiceDeleteRequest",
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
				  "commandType": "activeActivePscServiceDeleteRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2025-01-13T11:22:51.204189721Z",
				  "response": {
					"error": {
					  "type": "PSC_SERVICE_NOT_FOUND",
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
			expectedError:   errors.New("resource not found - subscription 114019 and region 1"),
			expectedErrorAs: &psc.NotFoundActiveActive{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			err = subject.PrivateServiceConnect.DeleteActiveActiveService(context.TODO(), 114019, 1)
			if testCase.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.IsType(t, err, testCase.expectedErrorAs)
				assert.EqualError(t, err, testCase.expectedError.Error())
			}
		})
	}
}

func TestDeletePrivateServiceConnectEndpoint(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully delete a private service connect endpoint",
			mockedResponse: []endpointRequest{
				deleteRequest(
					t,
					"/subscriptions/114019/private-service-connect/40/endpoints/50",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscEndpointDeleteRequest",
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
				  "commandType": "pscEndpointDeleteRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019
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
		},
		{
			description: "should fail when private service connect endpoint is not found",
			mockedResponse: []endpointRequest{
				deleteRequestWithStatus(
					t,
					"/subscriptions/114019/private-service-connect/40/endpoints/50",
					404,
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "pscEndpointDeleteRequest",
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
				  "commandType": "pscEndpointDeleteRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2025-01-13T11:22:51.204189721Z",
				  "response": {
					"error": {
					  "type": "PSC_SERVICE_NOT_FOUND",
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
			expectedErrorAs: &psc.NotFound{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			err = subject.PrivateServiceConnect.DeleteEndpoint(context.TODO(), 114019, 40, 50)
			if testCase.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.IsType(t, err, testCase.expectedErrorAs)
				assert.EqualError(t, err, testCase.expectedError.Error())
			}
		})
	}
}

func TestDeletePrivateServiceConnectEndpointActiveActive(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully delete a private service connect endpoint",
			mockedResponse: []endpointRequest{
				deleteRequest(
					t,
					"/subscriptions/114019/regions/1/private-service-connect/40/endpoints/50",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscEndpointDeleteRequest",
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
				  "commandType": "activeActivePscEndpointDeleteRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019
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
		},
		{
			description: "should fail when private service connect endpoint is not found",
			mockedResponse: []endpointRequest{
				deleteRequestWithStatus(
					t,
					"/subscriptions/114019/regions/1/private-service-connect/40/endpoints/50",
					404,
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePscEndpointDeleteRequest",
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
				  "commandType": "activeActivePscEndpointDeleteRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2025-01-13T11:22:51.204189721Z",
				  "response": {
					"error": {
					  "type": "PSC_SERVICE_NOT_FOUND",
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
			expectedError:   errors.New("resource not found - subscription 114019 and region 1"),
			expectedErrorAs: &psc.NotFoundActiveActive{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			err = subject.PrivateServiceConnect.DeleteActiveActiveEndpoint(context.TODO(), 114019, 1, 40, 50)
			if testCase.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.IsType(t, err, testCase.expectedErrorAs)
				assert.EqualError(t, err, testCase.expectedError.Error())
			}
		})
	}
}
