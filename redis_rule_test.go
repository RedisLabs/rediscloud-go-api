package rediscloud_api

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/access_control_lists/redis_rules"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBadRedisRule(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequest(
				t,
				"/acl/redisRules",
				`{"name": "Test Create Bad Redis Rule", "redisRule": "let-me-create-resources"}`,
				`{
				  "taskId": "9462ac71-2d6a-403c-835e-bb999c5fad4a",
				  "commandType": "aclRedisRuleCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T09:21:18.521282Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/9462ac71-2d6a-403c-835e-bb999c5fad4a",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(t, "/tasks/9462ac71-2d6a-403c-835e-bb999c5fad4a", `{
			  "taskId": "9462ac71-2d6a-403c-835e-bb999c5fad4a",
			  "commandType": "aclRedisRuleCreateRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T10:21:37.889922Z",
			  "response": {
				"error": {
				  "type": "ACL_REDIS_RULE_PATTERN_NOT_VALID",
				  "status": "400 BAD_REQUEST",
				  "description": "Invalid ACL redis rule: commands must start with a + or - sign, categories must start with +@ or -@ characters and keys must start with the ~ symbol"
				},
				"additionalInfo": "redis-acl-wrong-acl-rule-pattern"
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/9462ac71-2d6a-403c-835e-bb999c5fad4a",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	_, err = subject.RedisRules.Create(context.TODO(), redis_rules.CreateRedisRuleRequest{
		Name:      redis.String("Test Create Bad Redis Rule"),
		RedisRule: redis.String("let-me-create-resources"),
	})

	assert.Equal(t, &internal.Error{
		Type:        redis.String("ACL_REDIS_RULE_PATTERN_NOT_VALID"),
		Description: redis.String("Invalid ACL redis rule: commands must start with a + or - sign, categories must start with +@ or -@ characters and keys must start with the ~ symbol"),
		Status:      redis.String("400 BAD_REQUEST"),
	}, errors.Unwrap(err))
}

func TestCreateRedisRule(t *testing.T) {
	expected := 7457
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequest(
				t,
				"/acl/redisRules",
				`{"name": "Test Create Redis Rule", "redisRule": "+@let-me-create-resources"}`,
				`{
				  "taskId": "e3946019-994e-49f6-83bb-26694b3c241f",
				  "commandType": "aclRedisRuleCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T09:21:18.521282Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/e3946019-994e-49f6-83bb-26694b3c241f",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`,
			),
			// Task doesn't exist just yet
			getRequestWithStatus(t, "/tasks/e3946019-994e-49f6-83bb-26694b3c241f", 404, ""),
			// Task exists, has just started
			getRequest(t, "/tasks/e3946019-994e-49f6-83bb-26694b3c241f", `{
			  "taskId": "e3946019-994e-49f6-83bb-26694b3c241f",
			  "commandType": "aclRedisRuleCreateRequest",
			  "status": "initialized",
			  "timestamp": "2023-06-21T09:22:18.521282Z",
			  "response": {},
			  "_links": {
				"self": {
				  "rel": "task",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/e3946019-994e-49f6-83bb-26694b3c241f",
				  "title": "getTaskStatusUpdates",
				  "type": "GET"
				}
			  }
			}`),
			// Task exists, is in progress
			getRequest(t, "/tasks/e3946019-994e-49f6-83bb-26694b3c241f", `{
			  "taskId": "e3946019-994e-49f6-83bb-26694b3c241f",
			  "commandType": "aclRedisRuleCreateRequest",
			  "status": "processing-in-progress",
			  "timestamp": "2023-06-21T09:23:18.521282Z",
			  "response": {},
			  "_links": {
				"self": {
				  "rel": "task",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/e3946019-994e-49f6-83bb-26694b3c241f",
				  "title": "getTaskStatusUpdates",
				  "type": "GET"
				}
			  }
			}`),
			// Task complete
			getRequest(t, "/tasks/e3946019-994e-49f6-83bb-26694b3c241f", fmt.Sprintf(`{
			  "taskId": "e3946019-994e-49f6-83bb-26694b3c241f",
			  "commandType": "aclRedisRuleCreateRequest",
			  "status": "processing-completed",
			  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
			  "timestamp": "2023-06-21T09:24:18.521282Z",
			  "response": {
				"resourceId": %[1]d
			  },
			  "links": [
			  	{
				  "rel": "resource",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/%[1]d/acl/redisRules",
				  "title": "getACLRedisRulesInformation",
				  "type": "GET"
				},
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/e3946019-994e-49f6-83bb-26694b3c241f",
				  "type": "GET"
				}
			  ]
			}`, expected)),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.RedisRules.Create(context.TODO(), redis_rules.CreateRedisRuleRequest{
		Name:      redis.String("Test Create Redis Rule"),
		RedisRule: redis.String("+@let-me-create-resources"),
	})
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUpdateNonExistentRedisRule(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/acl/redisRules/40004",
				`{"name": "Test Update Redis Rule", "redisRule": "+@let-me-update-resources"}`,
				`{
				  "taskId": "dc134208-bb52-458a-a2c9-ce572b28e5be",
				  "commandType": "aclRedisRuleUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T09:37:08.99213Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/dc134208-bb52-458a-a2c9-ce572b28e5be",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/dc134208-bb52-458a-a2c9-ce572b28e5be", `{
			  "taskId": "dc134208-bb52-458a-a2c9-ce572b28e5be",
			  "commandType": "aclRedisRuleUpdateRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T09:37:09.221196Z",
			  "response": {
				"error": {
				  "type": "ACL_REDIS_RULE_NOT_FOUND",
				  "status": "404 NOT_FOUND",
				  "description": "ACL redis rule not found"
				}
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/dc134208-bb52-458a-a2c9-ce572b28e5be",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.RedisRules.Update(context.TODO(), 40004, redis_rules.CreateRedisRuleRequest{
		Name:      redis.String("Test Update Redis Rule"),
		RedisRule: redis.String("+@let-me-update-resources"),
	})

	assert.Equal(t, &internal.Error{
		Type:        redis.String("ACL_REDIS_RULE_NOT_FOUND"),
		Description: redis.String("ACL redis rule not found"),
		Status:      redis.String("404 NOT_FOUND"),
	}, errors.Unwrap(err))
}

func TestUpdateBadRedisRule(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/acl/redisRules/40000",
				`{"name": "Test Update Bad Redis Rule", "redisRule": "let-me-update-resources"}`,
				`{
				  "taskId": "51c67510-976d-495d-8917-50734b44ccfd",
				  "commandType": "aclRedisRuleUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T09:37:08.99213Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/51c67510-976d-495d-8917-50734b44ccfd",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/51c67510-976d-495d-8917-50734b44ccfd", `{
			  "taskId": "51c67510-976d-495d-8917-50734b44ccfd",
			  "commandType": "aclRedisRuleUpdateRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T10:34:06.295657Z",
			  "response": {
				"error": {
				  "type": "ACL_REDIS_RULE_PATTERN_NOT_VALID",
				  "status": "400 BAD_REQUEST",
				  "description": "Invalid ACL redis rule: commands must start with a + or - sign, categories must start with +@ or -@ characters and keys must start with the ~ symbol"
				},
				"additionalInfo": "redis-acl-wrong-acl-rule-pattern"
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/51c67510-976d-495d-8917-50734b44ccfd",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.RedisRules.Update(context.TODO(), 40000, redis_rules.CreateRedisRuleRequest{
		Name:      redis.String("Test Update Bad Redis Rule"),
		RedisRule: redis.String("let-me-update-resources"),
	})

	assert.Equal(t, &internal.Error{
		Type:        redis.String("ACL_REDIS_RULE_PATTERN_NOT_VALID"),
		Description: redis.String("Invalid ACL redis rule: commands must start with a + or - sign, categories must start with +@ or -@ characters and keys must start with the ~ symbol"),
		Status:      redis.String("400 BAD_REQUEST"),
	}, errors.Unwrap(err))
}

func TestUpdateRedisRule(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/acl/redisRules/20000",
				`{"name": "Test Update Redis Rule", "redisRule": "+@let-me-update-resources"}`,
				`{
				  "taskId": "9a4178d9-0c84-4a7c-b02a-b0a559757df9",
				  "commandType": "aclRedisRuleUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T09:21:18.521282Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/e3946019-994e-49f6-83bb-26694b3c241f",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`,
			),
			// Task doesn't exist just yet
			getRequestWithStatus(t, "/tasks/9a4178d9-0c84-4a7c-b02a-b0a559757df9", 404, ""),
			// Task exists, has just started
			getRequest(t, "/tasks/9a4178d9-0c84-4a7c-b02a-b0a559757df9", `{
			  "taskId": "9a4178d9-0c84-4a7c-b02a-b0a559757df9",
			  "commandType": "aclRedisRuleUpdateRequest",
			  "status": "initialized",
			  "timestamp": "2023-06-21T09:22:18.521282Z",
			  "response": {},
			  "_links": {
				"self": {
				  "rel": "task",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/9a4178d9-0c84-4a7c-b02a-b0a559757df9",
				  "title": "getTaskStatusUpdates",
				  "type": "GET"
				}
			  }
			}`),
			// Task exists, is in progress
			getRequest(t, "/tasks/9a4178d9-0c84-4a7c-b02a-b0a559757df9", `{
			  "taskId": "9a4178d9-0c84-4a7c-b02a-b0a559757df9",
			  "commandType": "aclRedisRuleUpdateRequest",
			  "status": "processing-in-progress",
			  "timestamp": "2023-06-21T09:23:18.521282Z",
			  "response": {},
			  "_links": {
				"self": {
				  "rel": "task",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/9a4178d9-0c84-4a7c-b02a-b0a559757df9",
				  "title": "getTaskStatusUpdates",
				  "type": "GET"
				}
			  }
			}`),
			// Task complete
			getRequest(t, "/tasks/9a4178d9-0c84-4a7c-b02a-b0a559757df9", `{
			  "taskId": "9a4178d9-0c84-4a7c-b02a-b0a559757df9",
			  "commandType": "aclRedisRuleUpdateRequest",
			  "status": "processing-completed",
			  "timestamp": "2023-06-21T09:24:18.521282Z",
			  "response": {
				"resourceId": 20000
			  },
			  "_links": {
				"self": {
				  "rel": "task",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/9a4178d9-0c84-4a7c-b02a-b0a559757df9",
				  "title": "getTaskStatusUpdates",
				  "type": "GET"
				}
			  }
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.RedisRules.Update(context.TODO(), 20000, redis_rules.CreateRedisRuleRequest{
		Name:      redis.String("Test Update Redis Rule"),
		RedisRule: redis.String("+@let-me-update-resources"),
	})
	require.NoError(t, err)
}

func TestDeleteNonExistentRedisRule(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			deleteRequest(
				t,
				"/acl/redisRules/40004",
				`{
				  "taskId": "2577eb1d-d5a8-468a-b504-947c7f0916a7",
				  "commandType": "aclRedisRuleDeleteRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T09:37:08.99213Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/2577eb1d-d5a8-468a-b504-947c7f0916a7",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/2577eb1d-d5a8-468a-b504-947c7f0916a7", `{
			  "taskId": "2577eb1d-d5a8-468a-b504-947c7f0916a7",
			  "commandType": "aclRedisRuleDeleteRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T09:37:09.221196Z",
			  "response": {
				"error": {
				  "type": "ACL_REDIS_RULE_NOT_FOUND",
				  "status": "404 NOT_FOUND",
				  "description": "ACL redis rule not found"
				}
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/2577eb1d-d5a8-468a-b504-947c7f0916a7",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.RedisRules.Delete(context.TODO(), 40004)

	assert.Equal(t, &internal.Error{
		Type:        redis.String("ACL_REDIS_RULE_NOT_FOUND"),
		Description: redis.String("ACL redis rule not found"),
		Status:      redis.String("404 NOT_FOUND"),
	}, errors.Unwrap(err))
}

func TestDeleteRedisRule(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			deleteRequest(
				t,
				"/acl/redisRules/20000",
				`{
				  "taskId": "7b54806b-cba7-4a14-8234-15c1538ccb7d",
				  "commandType": "aclRedisRuleDeleteRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T10:05:13.98099Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/7b54806b-cba7-4a14-8234-15c1538ccb7d",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`,
			),
			// Task doesn't exist just yet
			getRequestWithStatus(t, "/tasks/7b54806b-cba7-4a14-8234-15c1538ccb7d", 404, ""),
			// Task exists, has just started
			getRequest(t, "/tasks/7b54806b-cba7-4a14-8234-15c1538ccb7d", `{
			  "taskId": "7b54806b-cba7-4a14-8234-15c1538ccb7d",
			  "commandType": "aclRedisRuleDeleteRequest",
			  "status": "initialized",
			  "timestamp": "2023-06-21T09:22:18.521282Z",
			  "response": {},
			  "_links": {
				"self": {
				  "rel": "task",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/7b54806b-cba7-4a14-8234-15c1538ccb7d",
				  "title": "getTaskStatusUpdates",
				  "type": "GET"
				}
			  }
			}`),
			// Task exists, is in progress
			getRequest(t, "/tasks/7b54806b-cba7-4a14-8234-15c1538ccb7d", `{
			  "taskId": "7b54806b-cba7-4a14-8234-15c1538ccb7d",
			  "commandType": "aclRedisRuleDeleteRequest",
			  "status": "processing-in-progress",
			  "timestamp": "2023-06-21T09:23:18.521282Z",
			  "response": {},
			  "_links": {
				"self": {
				  "rel": "task",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/7b54806b-cba7-4a14-8234-15c1538ccb7d",
				  "title": "getTaskStatusUpdates",
				  "type": "GET"
				}
			  }
			}`),
			// Task complete
			getRequest(t, "/tasks/7b54806b-cba7-4a14-8234-15c1538ccb7d", `{
			  "taskId": "7b54806b-cba7-4a14-8234-15c1538ccb7d",
			  "commandType": "aclRedisRuleDeleteRequest",
			  "status": "processing-completed",
			  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
			  "timestamp": "2023-06-21T09:24:18.521282Z",
			  "response": {
				"resourceId": 20000
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/7b54806b-cba7-4a14-8234-15c1538ccb7d",
				  "type": "GET"
				}
			  ]
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.RedisRules.Delete(context.TODO(), 20000)
	require.NoError(t, err)
}

func TestListRedisRules(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(t, "/acl/redisRules", `{
			  "accountId": 53012,
			  "redisRules": [
				{
				  "id": 3923,
				  "name": "ACL-rule-example",
				  "acl": "+@all",
				  "isDefault": false,
				  "status": "active"
				},
				{
				  "id": 76,
				  "name": "Full-Access",
				  "acl": "+@all  ~*",
				  "isDefault": true,
				  "status": "active"
				},
				{
				  "id": 77,
				  "name": "Read-Write",
				  "acl": "+@all -@dangerous ~*",
				  "isDefault": true,
				  "status": "active"
				},
				{
				  "id": 78,
				  "name": "Read-Only",
				  "acl": "+@read ~*",
				  "isDefault": true,
				  "status": "active"
				}
			  ],
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/redisRules",
				  "type": "GET"
				}
			  ]
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.RedisRules.List(context.TODO())
	require.NoError(t, err)

	assert.ElementsMatch(t, []*redis_rules.GetRedisRuleResponse{
		{
			ID:        redis.Int(3923),
			Name:      redis.String("ACL-rule-example"),
			ACL:       redis.String("+@all"),
			IsDefault: redis.Bool(false),
			Status:    redis.String("active"),
		},
		{
			ID:        redis.Int(76),
			Name:      redis.String("Full-Access"),
			ACL:       redis.String("+@all  ~*"),
			IsDefault: redis.Bool(true),
			Status:    redis.String("active"),
		},
		{
			ID:        redis.Int(77),
			Name:      redis.String("Read-Write"),
			ACL:       redis.String("+@all -@dangerous ~*"),
			IsDefault: redis.Bool(true),
			Status:    redis.String("active"),
		},
		{
			ID:        redis.Int(78),
			Name:      redis.String("Read-Only"),
			ACL:       redis.String("+@read ~*"),
			IsDefault: redis.Bool(true),
			Status:    redis.String("active"),
		},
	}, actual)

}

func TestGetNonExistentRedisRule(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(t, "/acl/redisRules", `{
			  "accountId": 53012,
			  "redisRules": [
				{
				  "id": 3923,
				  "name": "ACL-rule-example",
				  "acl": "+@all",
				  "isDefault": false,
				  "status": "active"
				},
				{
				  "id": 76,
				  "name": "Full-Access",
				  "acl": "+@all  ~*",
				  "isDefault": true,
				  "status": "active"
				},
				{
				  "id": 77,
				  "name": "Read-Write",
				  "acl": "+@all -@dangerous ~*",
				  "isDefault": true,
				  "status": "active"
				},
				{
				  "id": 78,
				  "name": "Read-Only",
				  "acl": "+@read ~*",
				  "isDefault": true,
				  "status": "active"
				}
			  ],
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/redisRules",
				  "type": "GET"
				}
			  ]
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.RedisRules.Get(context.TODO(), 40004)

	assert.Nil(t, actual)
	assert.IsType(t, &redis_rules.NotFound{}, err)

}

func TestGetRedisRule(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(t, "/acl/redisRules", `{
			  "accountId": 53012,
			  "redisRules": [
				{
				  "id": 3923,
				  "name": "ACL-rule-example",
				  "acl": "+@all",
				  "isDefault": false,
				  "status": "active"
				},
				{
				  "id": 76,
				  "name": "Full-Access",
				  "acl": "+@all  ~*",
				  "isDefault": true,
				  "status": "active"
				},
				{
				  "id": 77,
				  "name": "Read-Write",
				  "acl": "+@all -@dangerous ~*",
				  "isDefault": true,
				  "status": "active"
				},
				{
				  "id": 78,
				  "name": "Read-Only",
				  "acl": "+@read ~*",
				  "isDefault": true,
				  "status": "active"
				}
			  ],
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/redisRules",
				  "type": "GET"
				}
			  ]
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.RedisRules.Get(context.TODO(), 3923)
	require.NoError(t, err)

	assert.Equal(t, &redis_rules.GetRedisRuleResponse{
		ID:        redis.Int(3923),
		Name:      redis.String("ACL-rule-example"),
		ACL:       redis.String("+@all"),
		IsDefault: redis.Bool(false),
		Status:    redis.String("active"),
	}, actual)

}
