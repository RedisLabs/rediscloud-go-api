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

func TestGetPrivateLink(t *testing.T) {
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
								  "status": "received",
								  "principals": [
									{
									  "principal": "arn:aws:iam::123456789012:root",
									  "status": "ready",
									  "alias": "some alias",
									  "type": "aws_account"
									}
								  ],
								  "resourceConfigurationId": "123456789012",
								  "resourceConfigurationArn": "arn:aws:iam::123456789012:root",
								  "shareArn": "arn:aws:iam::123456789012:root",
								  "shareName": "share name",
								  "connections": [
									{
									  "associationId": "received",
									  "connectionId": "vpce-con-12345678",
									  "type": "connection type",
									  "ownerId": "123456789012",
									  "associationDate": "2024-07-16T09:26:40.929904847Z"
									}
								  ],
								  "databases": [
									{
									  "databaseId": 0,
									  "port": 6379,
									  "rlEndpoint": ""
									}
								  ],
								  "subscriptionId": 114019,
								  "errorMessage": "no error"
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
				Status: redis.String("received"),
				Principals: []*pl.PrivateLinkPrincipal{
					{
						Principal: redis.String("arn:aws:iam::123456789012:root"),
						Status:    redis.String("ready"),
						Alias:     redis.String("some alias"),
						Type:      redis.String("aws_account"),
					},
				},
				ResourceConfigurationId:  redis.String("123456789012"),
				ResourceConfigurationArn: redis.String("arn:aws:iam::123456789012:root"),
				ShareArn:                 redis.String("arn:aws:iam::123456789012:root"),
				ShareName:                redis.String("share name"),
				Connections: []*pl.PrivateLinkConnection{{
					AssociationId:   redis.String("received"),
					ConnectionId:    redis.String("vpce-con-12345678"),
					Type:            redis.String("connection type"),
					OwnerId:         redis.String("123456789012"),
					AssociationDate: redis.String("2024-07-16T09:26:40.929904847Z"),
				}},
				Databases: []*pl.PrivateLinkDatabase{{
					DatabaseId:           redis.Int(0),
					Port:                 redis.Int(6379),
					ResourceLinkEndpoint: redis.String(""),
				}},
				SubscriptionId: redis.Int(114019),
				ErrorMessage:   redis.String("no error"),
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
			expectedError:   errors.New("privatelink resource not found - subscription 114019"),
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
			expectedError:   errors.New("privatelink resource not found - subscription 114019"),
			expectedErrorAs: &pl.NotFound{},
		},
		{
			description: "should fail when privatelink returns empty response",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/private-link",
					`{
				  "taskId": "a357a65b-801c-43ec-9434-a52a4b9102c9",
				  "commandType": "privateLinkGetRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2026-01-16T13:39:40.929904847Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/a357a65b-801c-43ec-9434-a52a4b9102c9",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
				),
				getRequest(
					t,
					"/tasks/a357a65b-801c-43ec-9434-a52a4b9102c9",
					`{
				  "taskId": "a357a65b-801c-43ec-9434-a52a4b9102c9",
				  "commandType": "privateLinkGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2026-01-16T13:39:50.259963336Z",
				  "response": {
					"resourceId": 127352,
					"resource": {
					  "links": []
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/a357a65b-801c-43ec-9434-a52a4b9102c9",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
				),
			},
			expectedError:   errors.New("privatelink resource not found - subscription 114019"),
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

func TestGetActiveActivePrivateLink(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *pl.PrivateLink
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return an active active privatelink config",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/regions/1/private-link",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePrivateLinkGetRequest",
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
							  "commandType": "activeActivePrivateLinkGetRequest",
							  "status": "processing-completed",
							  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
							  "timestamp": "2024-07-16T09:26:49.847808891Z",
							  "response": {
								"resourceId": 114019,
								"resource": {
								  "status": "received",
								  "principals": [
									{
									  "principal": "arn:aws:iam::123456789012:root",
									  "status": "ready",
									  "alias": "some alias",
									  "type": "aws_account"
									}
								  ],
								  "resourceConfigurationId": "123456789012",
								  "resourceConfigurationArn": "arn:aws:iam::123456789012:root",
								  "shareArn": "arn:aws:iam::123456789012:root",
								  "shareName": "share name",
								  "connections": [
									{
									  "associationId": "received",
									  "connectionId": "vpce-con-12345678",
									  "type": "connection type",
									  "ownerId": "123456789012",
									  "associationDate": "2024-07-16T09:26:40.929904847Z"
									}
								  ],
								  "databases": [
									{
									  "databaseId": 0,
									  "port": 6379,
									  "rlEndpoint": ""
									}
								  ],
								  "subscriptionId": 114019,
								  "regionId": 1,
								  "errorMessage": "no error"
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
				Status: redis.String("received"),
				Principals: []*pl.PrivateLinkPrincipal{
					{
						Principal: redis.String("arn:aws:iam::123456789012:root"),
						Status:    redis.String("ready"),
						Alias:     redis.String("some alias"),
						Type:      redis.String("aws_account"),
					},
				},
				ResourceConfigurationId:  redis.String("123456789012"),
				ResourceConfigurationArn: redis.String("arn:aws:iam::123456789012:root"),
				ShareArn:                 redis.String("arn:aws:iam::123456789012:root"),
				ShareName:                redis.String("share name"),
				Connections: []*pl.PrivateLinkConnection{{
					AssociationId:   redis.String("received"),
					ConnectionId:    redis.String("vpce-con-12345678"),
					Type:            redis.String("connection type"),
					OwnerId:         redis.String("123456789012"),
					AssociationDate: redis.String("2024-07-16T09:26:40.929904847Z"),
				}},
				Databases: []*pl.PrivateLinkDatabase{{
					DatabaseId:           redis.Int(0),
					Port:                 redis.Int(6379),
					ResourceLinkEndpoint: redis.String(""),
				}},
				SubscriptionId: redis.Int(114019),
				RegionId:       redis.Int(1),
				ErrorMessage:   redis.String("no error"),
			},
		},
		{
			description: "should fail when private link is not found",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/regions/1/private-link",
					`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "activeActivePrivateLinkGetRequest",
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
				  "commandType": "activeActivePrivateLinkGetRequest",
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
			expectedError:   errors.New("privatelink resource not found - subscription 114019"),
			expectedErrorAs: &pl.NotFound{},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				getRequestWithStatus(
					t,
					"/subscriptions/114019/regions/1/private-link",
					404,
					`{
					  "timestamp" : "2025-01-17T09:34:25.803+00:00",
					  "status" : 404,
					  "error" : "Not Found",
					  "path" : "/v1/subscriptions/114019/regions/1/private-link"
					}`),
			},
			expectedError:   errors.New("privatelink resource not found - subscription 114019"),
			expectedErrorAs: &pl.NotFound{},
		},
		{
			description: "should fail when active active privatelink returns empty response",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/regions/1/private-link",
					`{
				  "taskId": "b468b76c-912d-54fd-a045-b63b5c213db0",
				  "commandType": "activeActivePrivateLinkGetRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2026-01-16T13:39:40.929904847Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/b468b76c-912d-54fd-a045-b63b5c213db0",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
				),
				getRequest(
					t,
					"/tasks/b468b76c-912d-54fd-a045-b63b5c213db0",
					`{
				  "taskId": "b468b76c-912d-54fd-a045-b63b5c213db0",
				  "commandType": "activeActivePrivateLinkGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2026-01-16T13:39:50.259963336Z",
				  "response": {
					"resourceId": 127352,
					"resource": {
					  "links": []
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/b468b76c-912d-54fd-a045-b63b5c213db0",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
				),
			},
			expectedError:   errors.New("privatelink resource not found - subscription 114019, region 1"),
			expectedErrorAs: &pl.NotFoundActiveActive{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateLink.GetActiveActivePrivateLink(context.TODO(), 114019, 1)
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

func TestGetPrivateLinkScript(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *pl.PrivateLinkEndpointScript
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return a privatelink script",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/private-link/endpoint-script?includeTerraformAwsScript=true",
					`{
				  "taskId": "4e62bd68-06ce-4207-91cd-62b9d8b60dff",
				  "commandType": "privateLinkEndpointScriptGetRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-10-30T14:13:46.539111211Z",
				  "links": [
					{
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/4e62bd68-06ce-4207-91cd-62b9d8b60dff",
					  "rel": "task"
					}
				  ]
				}`,
				),
				getRequest(
					t,
					"/tasks/4e62bd68-06ce-4207-91cd-62b9d8b60dff",
					`{
							  "taskId": "4e62bd68-06ce-4207-91cd-62b9d8b60dff",
							  "commandType": "privateLinkEndpointScriptGetRequest",
							  "status": "processing-completed",
							  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
							  "timestamp": "2025-10-30T14:13:51.991415541Z",
							  "response": {
								"resourceId": 114019,
								"resource": {
								  "resourceEndpointScript": "python script content here",
								  "terraformAwsScript": "terraform script content here"
								}
							  },
							  "links": [
								{
								  "type": "GET",
								  "href": "https://api-staging.qa.redislabs.com/v1/tasks/4e62bd68-06ce-4207-91cd-62b9d8b60dff",
								  "rel": "self"
								}
							  ]
							}`,
				),
			},
			expectedResult: &pl.PrivateLinkEndpointScript{
				ResourceEndpointScript: redis.String("python script content here"),
				TerraformAwsScript:     redis.String("terraform script content here"),
			},
		},
	}
	for _, testCase := range tc {

		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateLink.GetPrivateLinkEndpointScript(context.TODO(), 114019)
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

func TestDeletePrivateLink(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully delete a privatelink",
			mockedResponse: []endpointRequest{
				deleteRequest(
					t,
					"/subscriptions/114019/private-link",
					`{
					  "taskId": "612fd31f-fd44-4cb0-a429-07882309a972",
					  "commandType": "privateLinkDeleteRequest",
					  "status": "received",
					  "description": "Task request received and is being queued for processing.",
					  "timestamp": "2024-07-16T09:26:40.929904847Z",
					  "links": [
						{
						  "href": "https://api-staging.qa.redislabs.com/v1/tasks/612fd31f-fd44-4cb0-a429-07882309a972",
						  "rel": "task",
						  "type": "GET"
						}
					  ]
					}`,
				),
				getRequest(
					t,
					"/tasks/612fd31f-fd44-4cb0-a429-07882309a972",
					`{
					  "taskId": "612fd31f-fd44-4cb0-a429-07882309a972",
					  "commandType": "privateLinkDeleteRequest",
					  "status": "processing-completed",
					  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					  "timestamp": "2024-07-16T09:26:49.847808891Z",
					  "links": [
						{
						  "href": "https://api-staging.qa.redislabs.com/v1/tasks/612fd31f-fd44-4cb0-a429-07882309a972",
						  "rel": "self",
						  "type": "GET"
						}
					  ]
					}`,
				),
			},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				deleteRequestWithStatus(
					t,
					"/subscriptions/114019/private-link",
					404,
					`{
					  "timestamp": "2025-01-17T09:34:25.803+00:00",
					  "status": 404,
					  "error": "Not Found",
					  "path": "/v1/subscriptions/114019/private-link"
					}`,
				),
			},
			expectedError:   errors.New("privatelink resource not found - subscription 114019"),
			expectedErrorAs: &pl.NotFound{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			err = subject.PrivateLink.DeletePrivateLink(context.TODO(), 114019)
			if testCase.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.IsType(t, err, testCase.expectedErrorAs)
				assert.EqualError(t, err, testCase.expectedError.Error())
			}
		})
	}
}

func TestDeleteActiveActivePrivateLink(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully delete an active active privatelink",
			mockedResponse: []endpointRequest{
				deleteRequest(
					t,
					"/subscriptions/114019/regions/1/private-link",
					`{
					  "taskId": "723ge42g-ge55-5dc1-b530-18993fb82083",
					  "commandType": "activeActivePrivateLinkDeleteRequest",
					  "status": "received",
					  "description": "Task request received and is being queued for processing.",
					  "timestamp": "2024-07-16T09:26:40.929904847Z",
					  "links": [
						{
						  "href": "https://api-staging.qa.redislabs.com/v1/tasks/723ge42g-ge55-5dc1-b530-18993fb82083",
						  "rel": "task",
						  "type": "GET"
						}
					  ]
					}`,
				),
				getRequest(
					t,
					"/tasks/723ge42g-ge55-5dc1-b530-18993fb82083",
					`{
					  "taskId": "723ge42g-ge55-5dc1-b530-18993fb82083",
					  "commandType": "activeActivePrivateLinkDeleteRequest",
					  "status": "processing-completed",
					  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
					  "timestamp": "2024-07-16T09:26:49.847808891Z",
					  "links": [
						{
						  "href": "https://api-staging.qa.redislabs.com/v1/tasks/723ge42g-ge55-5dc1-b530-18993fb82083",
						  "rel": "self",
						  "type": "GET"
						}
					  ]
					}`,
				),
			},
		},
		{
			description: "should fail when subscription is not found",
			mockedResponse: []endpointRequest{
				deleteRequestWithStatus(
					t,
					"/subscriptions/114019/regions/1/private-link",
					404,
					`{
					  "timestamp": "2025-01-17T09:34:25.803+00:00",
					  "status": 404,
					  "error": "Not Found",
					  "path": "/v1/subscriptions/114019/regions/1/private-link"
					}`,
				),
			},
			expectedError:   errors.New("privatelink resource not found - subscription 114019"),
			expectedErrorAs: &pl.NotFound{},
		},
	}

	for _, testCase := range tc {
		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			err = subject.PrivateLink.DeleteActiveActivePrivateLink(context.TODO(), 114019, 1)
			if testCase.expectedError == nil {
				assert.NoError(t, err)
			} else {
				assert.IsType(t, err, testCase.expectedErrorAs)
				assert.EqualError(t, err, testCase.expectedError.Error())
			}
		})
	}
}

func TestGetActiveActivePrivateLinkScript(t *testing.T) {
	tc := []struct {
		description     string
		mockedResponse  []endpointRequest
		expectedResult  *pl.PrivateLinkEndpointScript
		expectedError   error
		expectedErrorAs error
	}{
		{
			description: "should successfully return an active active privatelink script",
			mockedResponse: []endpointRequest{
				getRequest(
					t,
					"/subscriptions/114019/regions/1/private-link/endpoint-script?includeTerraformAwsScript=true",
					`{
				  "taskId": "5f73ce79-17df-5db1-a2de-18993ea71f00",
				  "commandType": "activeActivePrivateLinkEndpointScriptGetRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2025-10-30T14:13:46.539111211Z",
				  "links": [
					{
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/5f73ce79-17df-5db1-a2de-18993ea71f00",
					  "rel": "task"
					}
				  ]
				}`,
				),
				getRequest(
					t,
					"/tasks/5f73ce79-17df-5db1-a2de-18993ea71f00",
					`{
							  "taskId": "5f73ce79-17df-5db1-a2de-18993ea71f00",
							  "commandType": "activeActivePrivateLinkEndpointScriptGetRequest",
							  "status": "processing-completed",
							  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
							  "timestamp": "2025-10-30T14:13:51.991415541Z",
							  "response": {
								"resourceId": 114019,
								"resource": {
								  "resourceEndpointScript": "active active python script content here",
								  "terraformAwsScript": "active active terraform script content here"
								}
							  },
							  "links": [
								{
								  "type": "GET",
								  "href": "https://api-staging.qa.redislabs.com/v1/tasks/5f73ce79-17df-5db1-a2de-18993ea71f00",
								  "rel": "self"
								}
							  ]
							}`,
				),
			},
			expectedResult: &pl.PrivateLinkEndpointScript{
				ResourceEndpointScript: redis.String("active active python script content here"),
				TerraformAwsScript:     redis.String("active active terraform script content here"),
			},
		},
	}
	for _, testCase := range tc {

		t.Run(testCase.description, func(t *testing.T) {
			server := httptest.NewServer(
				testServer("key", "secret", testCase.mockedResponse...))

			subject, err := clientFromTestServer(server, "key", "secret")
			require.NoError(t, err)

			actual, err := subject.PrivateLink.GetActiveActivePrivateLinkEndpointScript(context.TODO(), 114019, 1)
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
