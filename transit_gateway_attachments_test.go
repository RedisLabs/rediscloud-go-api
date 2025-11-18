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

func TestListInvitations(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/114019/transitGateways/invitations",
				`{
				  "taskId": "a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
				  "commandType": "tgwListInvitationsRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-16T09:26:40.929904847Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
				`{
				  "taskId": "a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
				  "commandType": "tgwListInvitationsRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "invitations": [
						{
						  "id": 1,
						  "tgwId": 36,
						  "awsTgwUid": "tgw-0b92afdae97faaef8",
						  "status": "pending",
						  "awsAccountId": "620187402834"
						},
						{
						  "id": 2,
						  "tgwId": 37,
						  "awsTgwUid": "tgw-0c93bfeaf98gbbfg9",
						  "status": "pending",
						  "awsAccountId": "620187402834"
						}
					  ]
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
				`{
				  "taskId": "a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
				  "commandType": "tgwListInvitationsRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "invitations": [
						{
						  "id": 1,
						  "tgwId": 36,
						  "awsTgwUid": "tgw-0b92afdae97faaef8",
						  "status": "pending",
						  "awsAccountId": "620187402834"
						},
						{
						  "id": 2,
						  "tgwId": 37,
						  "awsTgwUid": "tgw-0c93bfeaf98gbbfg9",
						  "status": "pending",
						  "awsAccountId": "620187402834"
						}
					  ]
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.TransitGatewayAttachments.ListInvitations(context.TODO(), 114019)
	require.NoError(t, err)

	assert.Equal(t, &attachments.InvitationsResponse{
		CommandType: redis.String("tgwListInvitationsRequest"),
		Description: redis.String("Request processing completed successfully and its resources are now being provisioned / de-provisioned."),
		Status:      redis.String("processing-completed"),
		ID:          redis.String("a1b2c3d4-5e6f-7a8b-9c0d-1e2f3a4b5c6d"),
		Response: &attachments.InvitationResponseData{
			ResourceId: redis.Int(114019),
			Resource: &attachments.InvitationsResource{
				Invitations: []*attachments.TransitGatewayInvitation{
					{
						Id:           redis.Int(1),
						TgwId:        redis.Int(36),
						AwsTgwUid:    redis.String("tgw-0b92afdae97faaef8"),
						Status:       redis.String("pending"),
						AwsAccountId: redis.String("620187402834"),
					},
					{
						Id:           redis.Int(2),
						TgwId:        redis.Int(37),
						AwsTgwUid:    redis.String("tgw-0c93bfeaf98gbbfg9"),
						Status:       redis.String("pending"),
						AwsAccountId: redis.String("620187402834"),
					},
				},
			},
		},
	}, actual)
}

func TestListInvitationsActiveActive(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/114019/regions/1/transitGateways/invitations",
				`{
				  "taskId": "b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
				  "commandType": "tgwListInvitationsRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-16T09:26:40.929904847Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
				`{
				  "taskId": "b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
				  "commandType": "tgwListInvitationsRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "invitations": [
						{
						  "id": 3,
						  "tgwId": 38,
						  "awsTgwUid": "tgw-0d94cgfbg09hcchh0",
						  "status": "pending",
						  "awsAccountId": "620187402834"
						}
					  ]
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
				`{
				  "taskId": "b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
				  "commandType": "tgwListInvitationsRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T09:26:49.847808891Z",
				  "response": {
					"resourceId": 114019,
					"resource": {
					  "invitations": [
						{
						  "id": 3,
						  "tgwId": 38,
						  "awsTgwUid": "tgw-0d94cgfbg09hcchh0",
						  "status": "pending",
						  "awsAccountId": "620187402834"
						}
					  ]
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.TransitGatewayAttachments.ListInvitationsActiveActive(context.TODO(), 114019, 1)
	require.NoError(t, err)

	assert.Equal(t, &attachments.InvitationsResponse{
		CommandType: redis.String("tgwListInvitationsRequest"),
		Description: redis.String("Request processing completed successfully and its resources are now being provisioned / de-provisioned."),
		Status:      redis.String("processing-completed"),
		ID:          redis.String("b2c3d4e5-6f7a-8b9c-0d1e-2f3a4b5c6d7e"),
		Response: &attachments.InvitationResponseData{
			ResourceId: redis.Int(114019),
			Resource: &attachments.InvitationsResource{
				Invitations: []*attachments.TransitGatewayInvitation{
					{
						Id:           redis.Int(3),
						TgwId:        redis.Int(38),
						AwsTgwUid:    redis.String("tgw-0d94cgfbg09hcchh0"),
						Status:       redis.String("pending"),
						AwsAccountId: redis.String("620187402834"),
					},
				},
			},
		},
	}, actual)
}

func TestAcceptInvitation(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequestWithNoRequest(
				t,
				"/subscriptions/114019/transitGateways/invitations/1/accept",
				`{
				  "taskId": "c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f",
				  "commandType": "tgwAcceptInvitationRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-16T10:00:00.000000000Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f",
				`{
				  "taskId": "c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f",
				  "commandType": "tgwAcceptInvitationRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T10:00:10.000000000Z",
				  "response": {
					"resourceId": 1
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/c3d4e5f6-7a8b-9c0d-1e2f-3a4b5c6d7e8f",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.TransitGatewayAttachments.AcceptInvitation(context.TODO(), 114019, 1)
	require.NoError(t, err)
}

func TestAcceptInvitationActiveActive(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequestWithNoRequest(
				t,
				"/subscriptions/114019/regions/1/transitGateways/invitations/3/accept",
				`{
				  "taskId": "d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
				  "commandType": "tgwAcceptInvitationRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-16T10:00:00.000000000Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
				`{
				  "taskId": "d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
				  "commandType": "tgwAcceptInvitationRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T10:00:10.000000000Z",
				  "response": {
					"resourceId": 3
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/d4e5f6a7-8b9c-0d1e-2f3a-4b5c6d7e8f9a",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.TransitGatewayAttachments.AcceptInvitationActiveActive(context.TODO(), 114019, 1, 3)
	require.NoError(t, err)
}

func TestRejectInvitation(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequestWithNoRequest(
				t,
				"/subscriptions/114019/transitGateways/invitations/2/reject",
				`{
				  "taskId": "e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b",
				  "commandType": "tgwRejectInvitationRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-16T10:05:00.000000000Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b",
				`{
				  "taskId": "e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b",
				  "commandType": "tgwRejectInvitationRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T10:05:10.000000000Z",
				  "response": {
					"resourceId": 2
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/e5f6a7b8-9c0d-1e2f-3a4b-5c6d7e8f9a0b",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.TransitGatewayAttachments.RejectInvitation(context.TODO(), 114019, 2)
	require.NoError(t, err)
}

func TestRejectInvitationActiveActive(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequestWithNoRequest(
				t,
				"/subscriptions/114019/regions/1/transitGateways/invitations/3/reject",
				`{
				  "taskId": "f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c",
				  "commandType": "tgwRejectInvitationRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-16T10:05:00.000000000Z",
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c",
					  "rel": "task",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c",
				`{
				  "taskId": "f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c",
				  "commandType": "tgwRejectInvitationRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-16T10:05:10.000000000Z",
				  "response": {
					"resourceId": 3
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/f6a7b8c9-0d1e-2f3a-4b5c-6d7e8f9a0b1c",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.TransitGatewayAttachments.RejectInvitationActiveActive(context.TODO(), 114019, 1, 3)
	require.NoError(t, err)
}
