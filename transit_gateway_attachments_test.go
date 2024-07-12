package rediscloud_api

import (
	"context"
	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/transit_gateway/attachments"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
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
				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
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
				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				`{
				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "moreDetail": "goes here"
				}`,
			),
			getRequest(
				t,
				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				`{
				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "moreDetail": "goes here"
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

func TestGetActiveActiveAttachments(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(
				t,
				"/subscriptions/113779/regions/1/transitGateways",
				`{
				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
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
				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				`{
				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "moreDetail": "goes here"
				}`,
			),
			getRequest(
				t,
				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				`{
				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				  "commandType": "tgwGetRequest",
				  "status": "processing-completed",
				  "moreDetail": "goes here"
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.TransitGatewayAttachments.GetActiveActive(context.TODO(), 113779, 1)
	require.NoError(t, err)

	assert.Equal(t, &attachments.GetAttachmentsTask{
		CommandType: redis.String("tgwGetRequest"),
		Description: redis.String("Request processing completed successfully and its resources are now being provisioned / de-provisioned."),
		Status:      redis.String("processing-completed"),
		ID:          redis.String("ce2cbfea-9b15-4250-a516-f014161a8dd3"),
		Response: &attachments.Response{
			ResourceId: redis.Int(1),
			Resource: &attachments.Resource{
				TransitGatewayAttachment: []*attachments.TransitGatewayAttachment{
					{
						Id:               redis.String("1"),
						AwsTgwUid:        nil,
						AttachmentUid:    nil, // TODO Use this as the resource id in terraform! Hopefully the POST endpoint returns it...
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

func TestCreateAttachment(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequest(
				t,
				"/subscriptions/113779/transitGateways/123",
				"",
				// TODO tgwCreateRequest is a guess
				`{
				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				  "commandType": "tgwCreateRequest",
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
			// TODO What do these look like?
			getRequest(
				t,
				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				`{
				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				  "commandType": "tgwCreateRequest",
				  "status": "processing-completed",
				  "moreDetail": "goes here"
				}`,
			),
			getRequest(
				t,
				"/tasks/0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				`{
				  "taskId": "0477a7dc-008a-4b6e-a1dc-fb47722e2919",
				  "commandType": "tgwCreateRequest",
				  "status": "processing-completed",
				  "moreDetail": "goes here"
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.TransitGatewayAttachments.Create(context.TODO(), 113779, 123)
	require.NoError(t, err)

	assert.Equal(t, 456, actual)
}

// TODO CreateAA
// TODO Update
// TODO UpdateAA
// TODO Delete
// TODO DeleteAA
