package rediscloud_api

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/access_control_lists/roles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateNoSubscriptionRole(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequest(
				t,
				"/acl/roles",
				`{
				  "name": "ACL-role-example",
				  "redisRules": [
					{
					  "ruleName": "Read-Only",
					  "databases": [
						{
						  "subscriptionId": 0,
						  "databaseId": 0
						}
					  ]
					}
				  ]
				}`,
				`{
				  "taskId": "4c3073ae-f085-4e9c-aa5d-ea89980ed1fe",
				  "commandType": "aclRoleCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T10:44:36.203408Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/4c3073ae-f085-4e9c-aa5d-ea89980ed1fe",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/4c3073ae-f085-4e9c-aa5d-ea89980ed1fe",
				`{
				  "taskId": "4c3073ae-f085-4e9c-aa5d-ea89980ed1fe",
				  "commandType": "aclRoleCreateRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2023-06-21T10:44:36.552105Z",
				  "response": {
					"error": {
					  "type": "SUBSCRIPTION_NOT_FOUND",
					  "status": "404 NOT_FOUND",
					  "description": "Subscription was not found"
					}
				  },
				  "links": [
					{
					  "rel": "self",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/4c3073ae-f085-4e9c-aa5d-ea89980ed1fe",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	_, err = subject.Roles.Create(context.TODO(), roles.CreateRoleRequest{
		Name: redis.String("ACL-role-example"),
		RedisRules: []*roles.CreateRuleInRoleRequest{
			{
				RuleName: redis.String("Read-Only"),
				Databases: []*roles.CreateDatabaseInRuleInRoleRequest{
					{
						SubscriptionId: redis.Int(0),
						DatabaseId:     redis.Int(0),
						Regions:        []*string{},
					},
				},
			},
		},
	})

	require.Error(t, err)
}

func TestCreateRole(t *testing.T) {
	expected := 998
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequest(
				t,
				"/acl/roles",
				`{
				  "name": "ACL-role-example",
				  "redisRules": [
					{
					  "ruleName": "Read-Only",
					  "databases": [
						{
						  "subscriptionId": 156982,
						  "databaseId": 51332744
						}
					  ]
					}
				  ]
				}`,
				`{
				  "taskId": "78321c3e-5f6b-4968-a36e-c82d5f9e8a46",
				  "commandType": "aclRoleCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T11:30:42.605256Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/78321c3e-5f6b-4968-a36e-c82d5f9e8a46",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`,
			),
			// Task doesn't exist just yet
			getRequestWithStatus(t, "/tasks/78321c3e-5f6b-4968-a36e-c82d5f9e8a46", 404, ""),
			// Task exists, has just started
			getRequest(t, "/tasks/78321c3e-5f6b-4968-a36e-c82d5f9e8a46", `{
			  "taskId": "78321c3e-5f6b-4968-a36e-c82d5f9e8a46",
			  "commandType": "aclRoleCreateRequest",
			  "status": "initialized",
			  "timestamp": "2023-06-21T09:22:18.521282Z",
			  "response": {},
			  "_links": {
				"self": {
				  "rel": "task",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/78321c3e-5f6b-4968-a36e-c82d5f9e8a46",
				  "title": "getTaskStatusUpdates",
				  "type": "GET"
				}
			  }
			}`),
			// Task exists, is in progress
			getRequest(t, "/tasks/78321c3e-5f6b-4968-a36e-c82d5f9e8a46", `{
			  "taskId": "78321c3e-5f6b-4968-a36e-c82d5f9e8a46",
			  "commandType": "aclRoleCreateRequest",
			  "status": "processing-in-progress",
			  "timestamp": "2023-06-21T09:23:18.521282Z",
			  "response": {},
			  "_links": {
				"self": {
				  "rel": "task",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/78321c3e-5f6b-4968-a36e-c82d5f9e8a46",
				  "title": "getTaskStatusUpdates",
				  "type": "GET"
				}
			  }
			}`),
			// Task complete
			getRequest(t, "/tasks/78321c3e-5f6b-4968-a36e-c82d5f9e8a46", fmt.Sprintf(`{
			  "taskId": "78321c3e-5f6b-4968-a36e-c82d5f9e8a46",
			  "commandType": "aclRoleCreateRequest",
			  "status": "processing-completed",
			  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
			  "timestamp": "2023-06-21T11:30:43.265042Z",
			  "response": {
				"resourceId": %[1]d
			  },
			  "links": [
				{
				  "rel": "resource",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/%[1]d/acl/roles",
				  "title": "getACLRolesInformation",
				  "type": "GET"
				},
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/78321c3e-5f6b-4968-a36e-c82d5f9e8a46",
				  "type": "GET"
				}
			  ]
			}`, expected)),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Roles.Create(context.TODO(), roles.CreateRoleRequest{
		Name: redis.String("ACL-role-example"),
		RedisRules: []*roles.CreateRuleInRoleRequest{
			{
				RuleName: redis.String("Read-Only"),
				Databases: []*roles.CreateDatabaseInRuleInRoleRequest{
					{
						SubscriptionId: redis.Int(156982),
						DatabaseId:     redis.Int(51332744),
					},
				},
			},
		},
	})
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUpdateNonExistentRole(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/acl/roles/40004",
				`{
				  "name": "ACL-role-example",
				    "redisRules": [
  					{
  					  "ruleName": "Read-Only",
  					  "databases": [
  						{
						  "subscriptionId": 156982,
						  "databaseId": 51332744
						}
					  ]
					}
				  ]
				}`,
				`{
				  "taskId": "603aff45-eaed-48da-b085-6bd8cace67a7",
				  "commandType": "aclRoleUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T11:43:45.286377Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/603aff45-eaed-48da-b085-6bd8cace67a7",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/603aff45-eaed-48da-b085-6bd8cace67a7", `{
			  "taskId": "603aff45-eaed-48da-b085-6bd8cace67a7",
			  "commandType": "aclRoleUpdateRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T11:43:45.587192Z",
			  "response": {
				"error": {
				  "type": "ACL_ROLE_NOT_FOUND",
				  "status": "404 NOT_FOUND",
				  "description": "ACL role not found"
				}
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/603aff45-eaed-48da-b085-6bd8cace67a7",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Roles.Update(context.TODO(), 40004, roles.CreateRoleRequest{
		Name: redis.String("ACL-role-example"),
		RedisRules: []*roles.CreateRuleInRoleRequest{
			{
				RuleName: redis.String("Read-Only"),
				Databases: []*roles.CreateDatabaseInRuleInRoleRequest{
					{
						SubscriptionId: redis.Int(156982),
						DatabaseId:     redis.Int(51332744),
					},
				},
			},
		},
	})

	assert.Equal(t, &internal.Error{
		Type:        redis.String("ACL_ROLE_NOT_FOUND"),
		Description: redis.String("ACL role not found"),
		Status:      redis.String("404 NOT_FOUND"),
	}, errors.Unwrap(err))
}

func TestUpdateBadRole(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/acl/roles/40000",
				`{
					"name": "ACL-role-example",
				    "redisRules": [
  					{
  					  "ruleName": "Read-Only",
  					  "databases": [
  						{
						  "subscriptionId": 156983,
						  "databaseId": 0
						}
					  ]
					}
				  ]
				}`,
				`{
				  "taskId": "aaf9c2bf-23dd-471e-ab16-84cce2c53b51",
				  "commandType": "aclRoleUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T12:04:39.156272Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/aaf9c2bf-23dd-471e-ab16-84cce2c53b51",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/aaf9c2bf-23dd-471e-ab16-84cce2c53b51", `{
			  "taskId": "aaf9c2bf-23dd-471e-ab16-84cce2c53b51",
			  "commandType": "aclRoleUpdateRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T12:04:39.544693Z",
			  "response": {
				"error": {
				  "type": "DATABASE_NOT_FOUND",
				  "status": "404 NOT_FOUND",
				  "description": "Database was not found"
				}
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/aaf9c2bf-23dd-471e-ab16-84cce2c53b51",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Roles.Update(context.TODO(), 40000, roles.CreateRoleRequest{
		Name: redis.String("ACL-role-example"),
		RedisRules: []*roles.CreateRuleInRoleRequest{
			{
				RuleName: redis.String("Read-Only"),
				Databases: []*roles.CreateDatabaseInRuleInRoleRequest{
					{
						SubscriptionId: redis.Int(156983),
						DatabaseId:     redis.Int(0),
					},
				},
			},
		},
	})

	assert.Equal(t, &internal.Error{
		Type:        redis.String("DATABASE_NOT_FOUND"),
		Description: redis.String("Database was not found"),
		Status:      redis.String("404 NOT_FOUND"),
	}, errors.Unwrap(err))
}

func TestUpdateRole(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/acl/roles/20000",
				`{
					"name": "ACL-role-example",
				    "redisRules": [
  					{
  					  "ruleName": "Read-Only",
  					  "databases": [
  						{
						  "subscriptionId": 156983,
						  "databaseId": 51332750
						}
					  ]
					}
				  ]
				}`,
				`{
				  "taskId": "aaf9c2bf-23dd-471e-ab16-84cce2c53b51",
				  "commandType": "aclRoleUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T12:04:39.156272Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/aaf9c2bf-23dd-471e-ab16-84cce2c53b51",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/aaf9c2bf-23dd-471e-ab16-84cce2c53b51", `{
			  "taskId": "aaf9c2bf-23dd-471e-ab16-84cce2c53b51",
			  "commandType": "aclRoleUpdateRequest",
			  "status": "processing-completed",
			  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
			  "timestamp": "2023-06-21T12:30:39.787664Z",
			  "response": {
				"resourceId": 20000
			  },
			  "links": [
				{
				  "rel": "resource",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/998/acl/roles",
				  "title": "getACLRolesInformation",
				  "type": "GET"
				},
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/aaf9c2bf-23dd-471e-ab16-84cce2c53b51",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Roles.Update(context.TODO(), 20000, roles.CreateRoleRequest{
		Name: redis.String("ACL-role-example"),
		RedisRules: []*roles.CreateRuleInRoleRequest{
			{
				RuleName: redis.String("Read-Only"),
				Databases: []*roles.CreateDatabaseInRuleInRoleRequest{
					{
						SubscriptionId: redis.Int(156983),
						DatabaseId:     redis.Int(51332750),
						Regions:        nil,
					},
				},
			},
		},
	})
	require.NoError(t, err)
}

func TestDeleteNonExistentRole(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			deleteRequest(
				t,
				"/acl/roles/40004",
				`{
				  "taskId": "8ea93075-b99e-484f-9386-e0e623ec3245",
				  "commandType": "aclRoleDeleteRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T12:44:50.52562Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/8ea93075-b99e-484f-9386-e0e623ec3245",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/8ea93075-b99e-484f-9386-e0e623ec3245", `{
			  "taskId": "8ea93075-b99e-484f-9386-e0e623ec3245",
			  "commandType": "aclRoleDeleteRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T12:44:50.918624Z",
			  "response": {
				"error": {
				  "type": "ACL_ROLE_NOT_FOUND",
				  "status": "404 NOT_FOUND",
				  "description": "ACL role not found"
				}
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/8ea93075-b99e-484f-9386-e0e623ec3245",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Roles.Delete(context.TODO(), 40004)

	assert.Equal(t, &internal.Error{
		Type:        redis.String("ACL_ROLE_NOT_FOUND"),
		Description: redis.String("ACL role not found"),
		Status:      redis.String("404 NOT_FOUND"),
	}, errors.Unwrap(err))
}

func TestDeleteRole(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			deleteRequest(
				t,
				"/acl/roles/20000",
				`{
				  "taskId": "23b0bf8b-f97f-4e69-adb3-3f115ee74b7a",
				  "commandType": "aclRoleDeleteRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T12:48:24.959182Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/23b0bf8b-f97f-4e69-adb3-3f115ee74b7a",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/23b0bf8b-f97f-4e69-adb3-3f115ee74b7a", `{
			  "taskId": "23b0bf8b-f97f-4e69-adb3-3f115ee74b7a",
			  "commandType": "aclRoleDeleteRequest",
			  "status": "processing-completed",
			  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
			  "timestamp": "2023-06-21T12:48:25.353751Z",
			  "response": {
				"resourceId": 999
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/23b0bf8b-f97f-4e69-adb3-3f115ee74b7a",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Roles.Delete(context.TODO(), 20000)
	require.NoError(t, err)
}

func TestListRoles(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(t, "/acl/roles", `{
			  "accountId": 53012,
			  "roles": [
				{
				  "id": 998,
				  "name": "ACL-role-example",
				  "redisRules": [
					{
					  "ruleId": 78,
					  "ruleName": "Read-Only",
					  "databases": [
						{
						  "subscriptionId": 156983,
						  "databaseId": 51332750,
						  "databaseName": "john-test-database-keen-lab",
						  "regions": []
						}
					  ]
					}
				  ],
				  "users": [],
				  "status": "active"
				},
				{
				  "id": 999,
				  "name": "ACL-role-another-example",
				  "redisRules": [
					{
					  "ruleId": 78,
					  "ruleName": "Read-Only",
					  "databases": [
						{
						  "subscriptionId": 156983,
						  "databaseId": 51332750,
						  "databaseName": "john-test-database-keen-lab",
						  "regions": []
						}
					  ]
					}
				  ],
				  "users": [],
				  "status": "active"
				},
				{
				  "id": 27,
				  "name": "test-role",
				  "redisRules": [],
				  "users": [
					{
					  "id": 24,
					  "name": "test-user"
					}
				  ],
				  "status": "active"
				}
			  ],
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/roles",
				  "type": "GET"
				}
			  ]
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Roles.List(context.TODO())
	require.NoError(t, err)

	assert.ElementsMatch(t, []*roles.GetRoleResponse{
		{
			ID:   redis.Int(998),
			Name: redis.String("ACL-role-example"),
			RedisRules: []*roles.GetRuleInRoleResponse{
				{
					RuleId:   redis.Int(78),
					RuleName: redis.String("Read-Only"),
					Databases: []*roles.GetDatabaseInRuleInRoleResponse{
						{
							SubscriptionId: redis.Int(156983),
							DatabaseId:     redis.Int(51332750),
							DatabaseName:   redis.String("john-test-database-keen-lab"),
							Regions:        []*string{},
						},
					},
				},
			},
			Users:  []*roles.GetUserInRoleResponse{},
			Status: redis.String("active"),
		},
		{
			ID:   redis.Int(999),
			Name: redis.String("ACL-role-another-example"),
			RedisRules: []*roles.GetRuleInRoleResponse{
				{
					RuleId:   redis.Int(78),
					RuleName: redis.String("Read-Only"),
					Databases: []*roles.GetDatabaseInRuleInRoleResponse{
						{
							SubscriptionId: redis.Int(156983),
							DatabaseId:     redis.Int(51332750),
							DatabaseName:   redis.String("john-test-database-keen-lab"),
							Regions:        []*string{},
						},
					},
				},
			},
			Users:  []*roles.GetUserInRoleResponse{},
			Status: redis.String("active"),
		},
		{
			ID:         redis.Int(27),
			Name:       redis.String("test-role"),
			RedisRules: []*roles.GetRuleInRoleResponse{},
			Users: []*roles.GetUserInRoleResponse{
				{
					ID:   redis.Int(24),
					Name: redis.String("test-user"),
				},
			},
			Status: redis.String("active"),
		},
	}, actual)

}

func TestGetNonExistentRole(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(t, "/acl/roles", `{
			  "accountId": 53012,
			  "roles": [
				{
				  "id": 998,
				  "name": "ACL-role-example",
				  "redisRules": [
					{
					  "ruleId": 78,
					  "ruleName": "Read-Only",
					  "databases": [
						{
						  "subscriptionId": 156983,
						  "databaseId": 51332750,
						  "databaseName": "john-test-database-keen-lab",
						  "regions": []
						}
					  ]
					}
				  ],
				  "users": [],
				  "status": "active"
				},
				{
				  "id": 999,
				  "name": "ACL-role-another-example",
				  "redisRules": [
					{
					  "ruleId": 78,
					  "ruleName": "Read-Only",
					  "databases": [
						{
						  "subscriptionId": 156983,
						  "databaseId": 51332750,
						  "databaseName": "john-test-database-keen-lab",
						  "regions": []
						}
					  ]
					}
				  ],
				  "users": [],
				  "status": "active"
				},
				{
				  "id": 27,
				  "name": "test-role",
				  "redisRules": [],
				  "users": [
					{
					  "id": 24,
					  "name": "test-user"
					}
				  ],
				  "status": "active"
				}
			  ],
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/roles",
				  "type": "GET"
				}
			  ]
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Roles.Get(context.TODO(), 40004)

	assert.Nil(t, actual)
	assert.IsType(t, &roles.NotFound{}, err)

}

func TestGetRole(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(t, "/acl/roles", `{
			  "accountId": 53012,
			  "roles": [
				{
				  "id": 998,
				  "name": "ACL-role-example",
				  "redisRules": [
					{
					  "ruleId": 78,
					  "ruleName": "Read-Only",
					  "databases": [
						{
						  "subscriptionId": 156983,
						  "databaseId": 51332750,
						  "databaseName": "john-test-database-keen-lab",
						  "regions": []
						}
					  ]
					}
				  ],
				  "users": [],
				  "status": "active"
				},
				{
				  "id": 999,
				  "name": "ACL-role-another-example",
				  "redisRules": [
					{
					  "ruleId": 78,
					  "ruleName": "Read-Only",
					  "databases": [
						{
						  "subscriptionId": 156983,
						  "databaseId": 51332750,
						  "databaseName": "john-test-database-keen-lab",
						  "regions": []
						}
					  ]
					}
				  ],
				  "users": [
				    {
					  "id": 24,
					  "name": "test-user"
					}
				  ],
				  "status": "active"
				},
				{
				  "id": 27,
				  "name": "test-role",
				  "redisRules": [],
				  "users": [
					{
					  "id": 24,
					  "name": "test-user"
					}
				  ],
				  "status": "active"
				}
			  ],
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/roles",
				  "type": "GET"
				}
			  ]
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Roles.Get(context.TODO(), 999)
	require.NoError(t, err)

	assert.Equal(t, &roles.GetRoleResponse{
		ID:   redis.Int(999),
		Name: redis.String("ACL-role-another-example"),
		RedisRules: []*roles.GetRuleInRoleResponse{
			{
				RuleId:   redis.Int(78),
				RuleName: redis.String("Read-Only"),
				Databases: []*roles.GetDatabaseInRuleInRoleResponse{
					{
						SubscriptionId: redis.Int(156983),
						DatabaseId:     redis.Int(51332750),
						DatabaseName:   redis.String("john-test-database-keen-lab"),
						Regions:        []*string{},
					},
				},
			},
		},
		Users: []*roles.GetUserInRoleResponse{
			{
				ID:   redis.Int(24),
				Name: redis.String("test-user"),
			},
		},
		Status: redis.String("active"),
	}, actual)

}
