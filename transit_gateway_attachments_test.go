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
				"/subscriptions/113779/transitGateways",
				`{
				  "taskId": "268e9fa1-81de-4946-88a0-cd6c5c1c8dd2",
				  "commandType": "tgwGetRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2024-07-11T10:06:30.413894868Z",
				  "links": [
					{
					  "rel": "task",
					  "type": "GET",
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919"
					}
				  ]
				}
				`,
			),
			getRequest(
				t,
				"/tasks/268e9fa1-81de-4946-88a0-cd6c5c1c8dd2",
				`{
				  "taskId": "268e9fa1-81de-4946-88a0-cd6c5c1c8dd2",
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T14:42:07.772817377Z",
				  "response": {
					"resourceId": 113977,
					"resource": {
					  "tgws": []
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/268e9fa1-81de-4946-88a0-cd6c5c1c8dd2",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/268e9fa1-81de-4946-88a0-cd6c5c1c8dd2",
				`{
				  "taskId": "268e9fa1-81de-4946-88a0-cd6c5c1c8dd2",
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2024-07-15T14:42:07.772817377Z",
				  "response": {
					"resourceId": 113977,
					"resource": {
					  "tgws": []
					}
				  },
				  "links": [
					{
					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/268e9fa1-81de-4946-88a0-cd6c5c1c8dd2",
					  "rel": "self",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.TransitGatewayAttachments.Get(context.TODO(), 113779)
	require.NoError(t, err)

	assert.Equal(t, &attachments.GetAttachmentsTask{
		CommandType: redis.String("tgwGetRequest"),
		Description: redis.String("Request processing completed successfully and its resources are now being provisioned / de-provisioned."),
		Status:      redis.String("processing-completed"),
		ID:          redis.String("ce2cbfea-9b15-4250-a516-f014161a8dd3"),
		Response: &attachments.Response{
			ResourceId: redis.Int(1), // TODO What will this identify?
			Resource: &attachments.Resource{
				TransitGatewayAttachment: []*attachments.TransitGatewayAttachment{
					{
						Id:               redis.String("1"), // TODO What will this identify?
						AwsTgwUid:        nil,               // TODO Whose identifier is this?
						AttachmentUid:    nil,               // TODO Use this as the resource id in terraform!
						Status:           redis.String("ready"),
						AttachmentStatus: nil,
						AwsAccountId:     nil,
						Cidrs: []*attachments.Cidr{
							{
								CidrAddress: redis.String("10.0.0.0/24"),
								Status:      redis.String("ready"),
							},
						},
					},
				},
			},
		},
	}, actual)
}

// TODO
//func TestGetActiveActiveAttachments(t *testing.T) {
//	server := httptest.NewServer(
//		testServer(
//			"key",
//			"secret",
//			getRequest(
//				t,
//				"/subscriptions/113779/regions/1/transitGateways",
//				`{
//				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
//				  "commandType": "tgwGetRequest",
//				  "status": "received",
//				  "description": "Task request received and is being queued for processing.",
//				  "timestamp": "2024-07-11T10:06:30.413894868Z",
//				  "links": [
//					{
//					  "rel": "task",
//					  "type": "GET",
//					  "href": "https://api-staging.qa.redislabs.com/v1/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919"
//					}
//				  ]
//				}
//				`,
//			),
//			getRequest(
//				t,
//				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
//				`{
//				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
//				  "commandType": "tgwGetRequest",
//				  "status": "processing-completed",
//				  "moreDetail": "goes here"
//				}`,
//			),
//			getRequest(
//				t,
//				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
//				`{
//				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
//				  "commandType": "tgwGetRequest",
//				  "status": "processing-completed",
//				  "moreDetail": "goes here"
//				}`,
//			),
//		))
//
//	subject, err := clientFromTestServer(server, "key", "secret")
//	require.NoError(t, err)
//
//	actual, err := subject.TransitGatewayAttachments.GetActiveActive(context.TODO(), 113779, 1)
//	require.NoError(t, err)
//
//	assert.Equal(t, &attachments.GetAttachmentsTask{
//		CommandType: redis.String("tgwGetRequest"),
//		Description: redis.String("Request processing completed successfully and its resources are now being provisioned / de-provisioned."),
//		Status:      redis.String("processing-completed"),
//		ID:          redis.String("ce2cbfea-9b15-4250-a516-f014161a8dd3"),
//		Response: &attachments.Response{
//			ResourceId: redis.Int(1),
//			Resource: &attachments.Resource{
//				TransitGatewayAttachment: []*attachments.TransitGatewayAttachment{
//					{
//						Id:               redis.String("1"),
//						AwsTgwUid:        nil,
//						AttachmentUid:    nil, // TODO Use this as the resource id in terraform! Hopefully the POST endpoint returns it...
//						Status:           redis.String("ready"),
//						AttachmentStatus: nil,
//						AwsAccountId:     nil,
//						Cidrs: []*attachments.Cidr{
//							{
//								CidrAddress: redis.String("10.0.0.0/24"),
//								Status:      redis.String("ready"),
//							},
//						},
//					},
//				},
//			},
//		},
//	}, actual)
//}

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

// TODO CreateAA
// TODO Update
// TODO UpdateAA

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

// TODO DeleteAA
