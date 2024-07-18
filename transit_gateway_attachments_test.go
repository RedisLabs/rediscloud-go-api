package rediscloud_api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/transit_gateway/attachments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAttachments(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/114019/transitGateways",
				`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "tgwGetRequest",
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
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "tgws": [
						{
						  "id": 36,
						  "awsTgwUid": "tgw-0b92afdae97faaef8",
						  "status": "available",
						  "awsAccountId": "620187402834",
						  "cidrs": []
						}
					  ]
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
			getRequest(
				t,
				"/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
				`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "tgws": [
						{
						  "id": 36,
						  "awsTgwUid": "tgw-0b92afdae97faaef8",
						  "status": "available",
						  "awsAccountId": "620187402834",
						  "cidrs": []
						}
					  ]
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.TransitGatewayAttachments.Get(context.TODO(), 114019)
	require.NoError(t, err)

	assert.Equal(t, &attachments.GetAttachmentsTask{
		CommandType: redis.String("tgwGetRequest"),
		Description: redis.String("Request processing completed successfully and its resources are now being provisioned / de-provisioned."),
		Status:      redis.String("processing-completed"),
		ID:          redis.String("502fc31f-fd44-4cb0-a429-07882309a971"),
		Response: &attachments.Response{
			ResourceId: redis.Int(114019),
			Resource: &attachments.Resource{
				TransitGatewayAttachment: []*attachments.TransitGatewayAttachment{
					{
						Id:               redis.Int(36),
						AwsTgwUid:        redis.String("tgw-0b92afdae97faaef8"),
						AttachmentUid:    nil,
						Status:           redis.String("available"),
						AttachmentStatus: nil,
						AwsAccountId:     redis.String("620187402834"),
						Cidrs:            []*attachments.Cidr{},
					},
				},
			},
		},
	}, actual)
}

// TODO We haven't actually seen this working, at the moment we're assuming this is how it would happen
func TestGetActiveActiveAttachments(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/114019/regions/1/transitGateways",
				`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "tgwGetRequest",
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
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "tgws": [
						{
						  "id": 36,
						  "awsTgwUid": "tgw-0b92afdae97faaef8",
						  "status": "available",
						  "awsAccountId": "620187402834",
						  "cidrs": []
						}
					  ]
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
			getRequest(
				t,
				"/tasks/502fc31f-fd44-4cb0-a429-07882309a971",
				`{
				  "taskId": "502fc31f-fd44-4cb0-a429-07882309a971",
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "tgws": [
						{
						  "id": 36,
						  "awsTgwUid": "tgw-0b92afdae97faaef8",
						  "status": "available",
						  "awsAccountId": "620187402834",
						  "cidrs": []
						}
					  ]
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
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.TransitGatewayAttachments.GetActiveActive(context.TODO(), 114019, 1)
	require.NoError(t, err)

	assert.Equal(t, &attachments.GetAttachmentsTask{
		CommandType: redis.String("tgwGetRequest"),
		Description: redis.String("Request processing completed successfully and its resources are now being provisioned / de-provisioned."),
		Status:      redis.String("processing-completed"),
		ID:          redis.String("502fc31f-fd44-4cb0-a429-07882309a971"),
		Response: &attachments.Response{
			ResourceId: redis.Int(114019),
			Resource: &attachments.Resource{
				TransitGatewayAttachment: []*attachments.TransitGatewayAttachment{
					{
						Id:               redis.Int(36),
						AwsTgwUid:        redis.String("tgw-0b92afdae97faaef8"),
						AttachmentUid:    nil,
						Status:           redis.String("available"),
						AttachmentStatus: nil,
						AwsAccountId:     redis.String("620187402834"),
						Cidrs:            []*attachments.Cidr{},
					},
				},
			},
		},
	}, actual)
}

func TestCreateAttachment(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequestWithNoRequest(
				t,
				"/subscriptions/113991/transitGateways/35/attachment",
				`{
				  "taskId": "0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
				  "commandType": "tgwAttachmentCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-15T16:19:03.189459819Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}
				`,
			),
			getRequest(
				t,
				"/tasks/0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
				`{
				  "taskId": "0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
				  "commandType": "tgwAttachmentCreateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T16:19:12.197029024Z",
				  "response": {
					"resourceId": 35
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				`{
				  "taskId": "0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
				  "commandType": "tgwAttachmentCreateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T16:19:12.197029024Z",
				  "response": {
					"resourceId": 35
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.TransitGatewayAttachments.Create(context.TODO(), 113991, 35)
	require.NoError(t, err)

	assert.Equal(t, 35, actual)
}

// TODO Not seen this work in the real world
func TestCreateActiveActiveAttachment(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequestWithNoRequest(
				t,
				"/subscriptions/113991/regions/1/transitGateways/35/attachment",
				`{
				  "taskId": "0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
				  "commandType": "tgwAttachmentCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-15T16:19:03.189459819Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}
				`,
			),
			getRequest(
				t,
				"/tasks/0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
				`{
				  "taskId": "0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
				  "commandType": "tgwAttachmentCreateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T16:19:12.197029024Z",
				  "response": {
					"resourceId": 35
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				`{
				  "taskId": "0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
				  "commandType": "tgwAttachmentCreateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T16:19:12.197029024Z",
				  "response": {
					"resourceId": 35
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/0f9c5f49-3f71-428c-a50b-31cda9a35ed6",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.TransitGatewayAttachments.CreateActiveActive(context.TODO(), 113991, 1, 35)
	require.NoError(t, err)

	assert.Equal(t, 35, actual)
}

func TestUpdateAttachment(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/subscriptions/114133/transitGateways/47/attachment",
				`{
				  "cidrs": [
					"10.0.24.0/24",
					"10.0.25.0/24"
				  ]
				}`,
				`{
				  "taskId" : "175b7844-2b30-4934-bfd9-0bdcffc355cd",
				  "commandType" : "tgwUpdateCidrsRequest",
				  "status" : "received",
				  "description" : "Task request received and is being queued for processing.",
				  "timestamp" : "2024-07-18T10:15:19.575825248Z",
				  "links" : [ {
					"href" : "https://api-staging.qa.redislabs.com/v1/tasks/175b7844-2b30-4934-bfd9-0bdcffc355cd",
					"rel" : "task",
					"type" : "GET"
				  } ]
				}`,
			),
			getRequest(
				t,
				"/tasks/175b7844-2b30-4934-bfd9-0bdcffc355cd",
				`{
				  "taskId" : "175b7844-2b30-4934-bfd9-0bdcffc355cd",
				  "commandType" : "tgwUpdateCidrsRequest",
				  "status" : "processing-completed",
				  "description" : "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp" : "2024-07-18T10:15:30.526858477Z",
				  "response" : {
					"resourceId" : 47
				  },
				  "links" : [ {
					"href" : "https://api-staging.qa.redislabs.com/v1/tasks/175b7844-2b30-4934-bfd9-0bdcffc355cd",
					"rel" : "self",
					"type" : "GET"
				  } ]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	request := redis.StringSlice("10.0.24.0/24", "10.0.25.0/24")
	err = subject.TransitGatewayAttachments.Update(context.TODO(), 114133, 47, request)
	require.NoError(t, err)
}

// TODO Same issue as AA requests above
func TestUpdateActiveActiveAttachment(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/subscriptions/114133/regions/1/transitGateways/47/attachment",
				`{
				  "cidrs": [
					"10.0.24.0/24",
					"10.0.25.0/24"
				  ]
				}`,
				`{
				  "taskId" : "175b7844-2b30-4934-bfd9-0bdcffc355cd",
				  "commandType" : "tgwUpdateCidrsRequest",
				  "status" : "received",
				  "description" : "Task request received and is being queued for processing.",
				  "timestamp" : "2024-07-18T10:15:19.575825248Z",
				  "links" : [ {
					"href" : "https://api-staging.qa.redislabs.com/v1/tasks/175b7844-2b30-4934-bfd9-0bdcffc355cd",
					"rel" : "task",
					"type" : "GET"
				  } ]
				}`,
			),
			getRequest(
				t,
				"/tasks/175b7844-2b30-4934-bfd9-0bdcffc355cd",
				`{
				  "taskId" : "175b7844-2b30-4934-bfd9-0bdcffc355cd",
				  "commandType" : "tgwUpdateCidrsRequest",
				  "status" : "processing-completed",
				  "description" : "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp" : "2024-07-18T10:15:30.526858477Z",
				  "response" : {
					"resourceId" : 47
				  },
				  "links" : [ {
					"href" : "https://api-staging.qa.redislabs.com/v1/tasks/175b7844-2b30-4934-bfd9-0bdcffc355cd",
					"rel" : "self",
					"type" : "GET"
				  } ]
				}`,
			),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	request := redis.StringSlice("10.0.24.0/24", "10.0.25.0/24")
	err = subject.TransitGatewayAttachments.UpdateActiveActive(context.TODO(), 114133, 1, 47, request)
	require.NoError(t, err)
}

func TestDeleteAttachment(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			deleteRequest(
				t,
				"/subscriptions/113991/transitGateways/35/attachment",
				`{
				  "taskId": "df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				  "commandType": "tgwAttachmentDeleteRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-15T16:29:00.318877687Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				`{
				  "taskId": "df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				  "commandType": "tgwAttachmentDeleteRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T16:29:09.349089948Z",
				  "response": {
					"resourceId": 35
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				`{
				  "taskId": "df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				  "commandType": "tgwAttachmentDeleteRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T16:29:09.349089948Z",
				  "response": {
					"resourceId": 35
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.TransitGatewayAttachments.Delete(context.TODO(), 113991, 35)
	require.NoError(t, err)
}

// TODO Same as AA requests above
func TestDeleteActiveActiveAttachment(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			deleteRequest(
				t,
				"/subscriptions/113991/regions/1/transitGateways/35/attachment",
				`{
				  "taskId": "df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				  "commandType": "tgwAttachmentDeleteRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-15T16:29:00.318877687Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				`{
				  "taskId": "df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				  "commandType": "tgwAttachmentDeleteRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T16:29:09.349089948Z",
				  "response": {
					"resourceId": 35
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				`{
				  "taskId": "df1599d1-8ed2-46d0-819a-a92e5838d2cd",
				  "commandType": "tgwAttachmentDeleteRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T16:29:09.349089948Z",
				  "response": {
					"resourceId": 35
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/df1599d1-8ed2-46d0-819a-a92e5838d2cd",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.TransitGatewayAttachments.DeleteActiveActive(context.TODO(), 113991, 1, 35)
	require.NoError(t, err)
}
