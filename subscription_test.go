package rediscloud_api

import (
	"context"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/subscriptions"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubscription_Create(t *testing.T) {
	expected := 1235
	// Also test that the task API will poll for the finished task
	s := httptest.NewServer(testServer("key", "secret", postRequest(t, "/subscriptions", `{
  "name": "Test subscription",
  "dryRun": false,
  "paymentMethodId": 2,
  "paymentMethod": "credit-card",
  "memoryStorage": "ram",
  "cloudProviders": [
    {
      "provider": "AWS",
      "cloudAccountId": 1,
      "regions": [
        {
          "region": "eu-west-1"
        }
      ]
    }
  ],
  "databases": [
    {
      "name": "example",
      "protocol": "redis",
      "datasetSizeInGb": 1,
      "supportOSSClusterApi": true,
      "dataPersistence": "none",
      "replication": false,
      "throughputMeasurement": {
        "by": "operations-per-second",
        "value": 10000
      },
      "quantity": 1
    }
  ],
  "redisVersion": "latest"
}`, `{
  "taskId": "task-id",
  "commandType": "subscriptionCreateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequestWithStatus(t, "/tasks/task-id", 404, ""), getRequest(t, "/tasks/task-id", `{
  "taskId": "task-id",
  "commandType": "subscriptionCreateRequest",
  "status": "initialized",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {},
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task-id", `{
  "taskId": "task-id",
  "commandType": "subscriptionCreateRequest",
  "status": "processing-in-progress",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {},
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task-id", fmt.Sprintf(`{
  "taskId": "task-id",
  "commandType": "subscriptionCreateRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
    "resourceId": %d
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`, expected))))

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Subscription.Create(context.TODO(), subscriptions.CreateSubscription{
		Name:            redis.String("Test subscription"),
		DryRun:          redis.Bool(false),
		PaymentMethodID: redis.Int(2),
		PaymentMethod:   redis.String("credit-card"),
		MemoryStorage:   redis.String("ram"),
		CloudProviders: []*subscriptions.CreateCloudProvider{
			{
				Provider:       redis.String("AWS"),
				CloudAccountID: redis.Int(1),
				Regions: []*subscriptions.CreateRegion{
					{
						Region: redis.String("eu-west-1"),
					},
				},
			},
		},
		Databases: []*subscriptions.CreateDatabase{
			{
				Name:     redis.String("example"),
				Protocol: redis.String("redis"),
				// MemoryLimitInGB:      redis.Float64(1),
				DatasetSizeInGB:      redis.Float64(1),
				SupportOSSClusterAPI: redis.Bool(true),
				DataPersistence:      redis.String("none"),
				Replication:          redis.Bool(false),
				ThroughputMeasurement: &subscriptions.CreateThroughput{
					By:    redis.String("operations-per-second"),
					Value: redis.Int(10000),
				},
				Quantity: redis.Int(1),
			},
		},
		RedisVersion: redis.String("latest"),
	})
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSubscription_List(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequest(t, "/subscriptions", `{
  "accountId": 53012,
  "subscriptions": [
    {
      "id": 1,
      "name": "sdk",
      "status": "active",
      "paymentMethodType": "credit-card",
      "paymentMethodId": 2,
      "memoryStorage": "ram",
      "storageEncryption": false,
      "numberOfDatabases": 1,
      "subscriptionPricing": [
        {
          "type": "Shards",
          "quantity": 1,
          "quantityMeasurement": "shards",
          "pricePerUnit": 1,
          "priceCurrency": "USD",
          "pricePeriod": "hour"
        },
        {
          "type": "EBS Volume",
          "quantity": 71,
          "quantityMeasurement": "GB"
        },
        {
          "type": "c5.xlarge",
          "quantity": 2,
          "quantityMeasurement": "instances"
        },
        {
          "type": "m5.large",
          "quantity": 1,
          "quantityMeasurement": "instances"
        }
      ],
      "cloudDetails": [
        {
          "provider": "AWS",
          "cloudAccountId": 2,
          "totalSizeInGb": 0.0062,
          "regions": [
            {
              "region": "eu-west-1",
              "networking": [
                {
                  "deploymentCIDR": "10.0.0.0/24",
                  "subnetId": "subnet-12345"
                }
              ],
              "preferredAvailabilityZones": [
                "eu-west-1a"
              ],
              "multipleAvailabilityZones": false
            }
          ]
        }
      ]
    },
    {
      "id": 2,
      "name": "TF Example Subscription demo",
      "status": "pending",
      "paymentMethodId": 3,
      "memoryStorage": "ram",
      "storageEncryption": false,
      "numberOfDatabases": 0,
      "subscriptionPricing": [
        {
          "type": "Shards",
          "quantity": 1,
          "quantityMeasurement": "shards",
          "pricePerUnit": 2,
          "priceCurrency": "USD",
          "pricePeriod": "hour"
        },
        {
          "type": "EBS Volume",
          "quantity": 71,
          "quantityMeasurement": "GB"
        },
        {
          "type": "c5.xlarge",
          "quantity": 2,
          "quantityMeasurement": "instances"
        },
        {
          "type": "m5.large",
          "quantity": 1,
          "quantityMeasurement": "instances"
        }
      ],
      "cloudDetails": []
    }
  ],
  "_links": {
    "self": {
      "href": "https://qa-api.redislabs.com/v1/subscriptions",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.Subscription.List(context.TODO())
	require.NoError(t, err)

	assert.ElementsMatch(t, []*subscriptions.Subscription{
		{
			ID:                redis.Int(1),
			Name:              redis.String("sdk"),
			Status:            redis.String("active"),
			PaymentMethod:     redis.String("credit-card"),
			PaymentMethodID:   redis.Int(2),
			MemoryStorage:     redis.String("ram"),
			StorageEncryption: redis.Bool(false),
			NumberOfDatabases: redis.Int(1),
			CloudDetails: []*subscriptions.CloudDetail{
				{
					Provider:       redis.String("AWS"),
					CloudAccountID: redis.Int(2),
					TotalSizeInGB:  redis.Float64(0.0062),
					Regions: []*subscriptions.Region{
						{
							Region: redis.String("eu-west-1"),
							Networking: []*subscriptions.Networking{
								{
									DeploymentCIDR: redis.String("10.0.0.0/24"),
									SubnetID:       redis.String("subnet-12345"),
								},
							},
							PreferredAvailabilityZones: redis.StringSlice("eu-west-1a"),
							MultipleAvailabilityZones:  redis.Bool(false),
						},
					},
				},
			},
		},
		{
			ID:                redis.Int(2),
			Name:              redis.String("TF Example Subscription demo"),
			Status:            redis.String("pending"),
			PaymentMethodID:   redis.Int(3),
			MemoryStorage:     redis.String("ram"),
			StorageEncryption: redis.Bool(false),
			NumberOfDatabases: redis.Int(0),
			CloudDetails:      []*subscriptions.CloudDetail{},
		},
	}, actual)
}

func TestSubscription_Get(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequest(t, "/subscriptions/98765", `{
  "id": 1,
  "name": "Get-test",
  "status": "active",
  "paymentMethodType": "credit-card",
  "paymentMethodId": 2,
  "memoryStorage": "ram",
  "storageEncryption": false,
  "numberOfDatabases": 1,
  "subscriptionPricing": [
    {
      "type": "Shards",
      "quantity": 1,
      "quantityMeasurement": "shards",
      "pricePerUnit": 0.1,
      "priceCurrency": "USD",
      "pricePeriod": "hour"
    },
    {
      "type": "EBS Volume",
      "quantity": 71,
      "quantityMeasurement": "GB"
    },
    {
      "type": "c5.xlarge",
      "quantity": 2,
      "quantityMeasurement": "instances"
    },
    {
      "type": "m5.large",
      "quantity": 1,
      "quantityMeasurement": "instances"
    }
  ],
  "cloudDetails": [
    {
      "provider": "AWS",
      "cloudAccountId": 3,
      "totalSizeInGb": 4,
      "regions": [
        {
          "region": "eu-west-1",
          "networking": [
            {
              "deploymentCIDR": "10.0.0.0/24",
              "subnetId": "subnet-123456"
            }
          ],
          "preferredAvailabilityZones": [
            "eu-west-1a"
          ],
          "multipleAvailabilityZones": false
        }
      ]
    }
  ],
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.Subscription.Get(context.TODO(), 98765)
	require.NoError(t, err)

	assert.Equal(t, &subscriptions.Subscription{
		ID:                redis.Int(1),
		Name:              redis.String("Get-test"),
		Status:            redis.String("active"),
		PaymentMethod:     redis.String("credit-card"),
		PaymentMethodID:   redis.Int(2),
		MemoryStorage:     redis.String("ram"),
		StorageEncryption: redis.Bool(false),
		NumberOfDatabases: redis.Int(1),
		CloudDetails: []*subscriptions.CloudDetail{
			{
				Provider:       redis.String("AWS"),
				CloudAccountID: redis.Int(3),
				TotalSizeInGB:  redis.Float64(4),
				Regions: []*subscriptions.Region{
					{
						Region: redis.String("eu-west-1"),
						Networking: []*subscriptions.Networking{
							{
								DeploymentCIDR: redis.String("10.0.0.0/24"),
								SubnetID:       redis.String("subnet-123456"),
							},
						},
						PreferredAvailabilityZones: redis.StringSlice("eu-west-1a"),
						MultipleAvailabilityZones:  redis.Bool(false),
					},
				},
			},
		},
	}, actual)
}

func TestSubscription_Get_wraps404Error(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequestWithStatus(t, "/subscriptions/123", 404, "")))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.Subscription.Get(context.TODO(), 123)

	assert.Nil(t, actual)
	assert.IsType(t, &subscriptions.NotFound{}, err)
}

func TestSubscription_Update(t *testing.T) {
	s := httptest.NewServer(testServer("key", "secret", putRequest(t, "/subscriptions/1234", `{
  "name": "test",
  "paymentMethodId": 7
}`, `{
  "taskId": "task",
  "commandType": "subscriptionUpdateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "subscriptionUpdateRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	err = subject.Subscription.Update(context.TODO(), 1234, subscriptions.UpdateSubscription{
		Name:            redis.String("test"),
		PaymentMethodID: redis.Int(7),
	})
	require.NoError(t, err)
}

func TestSubscription_Delete(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", deleteRequest(t, "/subscriptions/12356", `{
  "taskId": "task",
  "commandType": "subscriptionDeleteRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "subscriptionDeleteRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.Subscription.Delete(context.TODO(), 12356)
	require.NoError(t, err)
}

func TestSubscription_GetCIDRWhitelist(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequest(t, "/subscriptions/12356/cidr", `{
  "taskId": "task",
  "commandType": "peeringListRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "peeringListRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
    "resource": {
      "cidr_ips": ["1", "2", "3"],
      "security_group_ids": ["4", "5", "6"]
    }
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.Subscription.GetCIDRAllowlist(context.TODO(), 12356)
	require.NoError(t, err)

	assert.Equal(t, &subscriptions.CIDRAllowlist{
		CIDRIPs:          redis.StringSlice("1", "2", "3"),
		SecurityGroupIDs: redis.StringSlice("4", "5", "6"),
	}, actual)
}

func TestSubscription_UpdateCIDRWhitelist(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", putRequest(t, "/subscriptions/12356/cidr", `{
  "cidrIps": ["6", "5"],
  "securityGroupIds": ["a", "b"]
}`, `{
  "taskId": "task",
  "commandType": "cidrUpdateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "cidrUpdateRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
    "resource": {
      "cidr_ips": ["1", "2", "3"],
      "security_group_ids": ["4", "5", "6"]
    }
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.Subscription.UpdateCIDRAllowlist(context.TODO(), 12356, subscriptions.UpdateCIDRAllowlist{
		CIDRIPs:          redis.StringSlice("6", "5"),
		SecurityGroupIDs: redis.StringSlice("a", "b"),
	})
	require.NoError(t, err)
}

func TestSubscription_ListVPCPeering(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequest(t, "/subscriptions/12356/peerings", `{
  "taskId": "task",
  "commandType": "peeringListRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "peeringListRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
    "resourceId" : 12356,
    "resource" : {
      "peerings" : [
		{
          "vpcPeeringId": 10,
          "awsAccountId": "4291",
          "vpcUid": "vpc-deadbeef",
          "vpcCidr": "10.0.0.0/24",
          "awsPeeringUid": "pcx-0123456789",
          "status": "done",
          "regionName": "eu-west-2"
		}
	  ]
    }
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.Subscription.ListVPCPeering(context.TODO(), 12356)
	require.NoError(t, err)

	assert.ElementsMatch(t, []*subscriptions.VPCPeering{
		{
			ID:           redis.Int(10),
			AWSAccountID: redis.String("4291"),
			VPCId:        redis.String("vpc-deadbeef"),
			VPCCidr:      redis.String("10.0.0.0/24"),
			AWSPeeringID: redis.String("pcx-0123456789"),
			Status:       redis.String("done"),
			Region:       redis.String("eu-west-2"),
		},
	}, actual)
}

func TestSubscription_ListVPCPeering_gcp(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", getRequest(t, "/subscriptions/12356/peerings", `{
  "taskId": "task",
  "commandType": "peeringListRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", `{
  "taskId": "task",
  "commandType": "vpcPeeringGetRequest",
  "status": "processing-completed",
  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
  "timestamp": "2020-12-01T14:56:09.204Z",
  "response": {
    "resourceId": 12356,
    "resource": {
      "peerings": [
        {
          "vpcPeeringId": 11,
          "projectUid": "cloud-api-123456",
          "networkName": "cloud-api-vpc-peering-test",
          "redisProjectUid": "v00d1c1f22233333f-tp",
          "redisNetworkName": "c12345-us-east1-2-rlrcp",
          "cloudPeeringId": "redislabs-peering-f123abc4-d56",
          "status": "inactive"
        }
      ]
    }
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	actual, err := subject.Subscription.ListVPCPeering(context.TODO(), 12356)
	require.NoError(t, err)

	assert.ElementsMatch(t, []*subscriptions.VPCPeering{
		{
			ID:               redis.Int(11),
			Status:           redis.String("inactive"),
			GCPProjectUID:    redis.String("cloud-api-123456"),
			NetworkName:      redis.String("cloud-api-vpc-peering-test"),
			RedisProjectUID:  redis.String("v00d1c1f22233333f-tp"),
			RedisNetworkName: redis.String("c12345-us-east1-2-rlrcp"),
			CloudPeeringID:   redis.String("redislabs-peering-f123abc4-d56"),
		},
	}, actual)
}

func TestSubscription_CreateVPCPeering(t *testing.T) {
	expected := 1235
	s := httptest.NewServer(testServer("key", "secret", postRequest(t, "/subscriptions/42/peerings", `{
  "region": "us-east-1",
  "awsAccountId": "098765",
  "vpcId": "17",
  "vpcCidr": "192.168.0.0/24"
}`, `{
  "taskId": "task-id",
  "commandType": "peeringCreateRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task-id", fmt.Sprintf(`{
  "taskId": "task-id",
  "commandType": "peeringCreateRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
    "resourceId": %d
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`, expected))))

	subject, err := clientFromTestServer(s, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Subscription.CreateVPCPeering(context.TODO(), 42, subscriptions.CreateVPCPeering{
		Region:       redis.String("us-east-1"),
		AWSAccountID: redis.String("098765"),
		VPCId:        redis.String("17"),
		VPCCidr:      redis.String("192.168.0.0/24"),
	})
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestSubscription_DeleteVPCPeering(t *testing.T) {
	s := httptest.NewServer(testServer("apiKey", "secret", deleteRequest(t, "/subscriptions/2/peerings/20", `{
  "taskId": "task",
  "commandType": "peeringDeleteRequest",
  "status": "received",
  "description": "Task request received and is being queued for processing.",
  "timestamp": "2020-11-02T09:05:34.3Z",
  "_links": {
    "task": {
      "href": "https://example.org",
      "title": "getTaskStatusUpdates",
      "type": "GET"
    }
  }
}`), getRequest(t, "/tasks/task", `{
  "taskId": "e02b40d6-1395-4861-a3b9-ecf829d835fd",
  "commandType": "peeringDeleteRequest",
  "status": "processing-completed",
  "timestamp": "2020-10-28T09:58:16.798Z",
  "response": {
  },
  "_links": {
    "self": {
      "href": "https://example.com",
      "type": "GET"
    }
  }
}`)))

	subject, err := clientFromTestServer(s, "apiKey", "secret")
	require.NoError(t, err)

	err = subject.Subscription.DeleteVPCPeering(context.TODO(), 2, 20)
	require.NoError(t, err)
}
