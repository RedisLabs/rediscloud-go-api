package rediscloud_api

import (
	"context"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/access_control_lists/users"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateBadUser(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequest(
				t,
				"/acl/users",
				`{
				  "name": "ACL-user-example",
				  "role": "Redis-role-example-that-does-not-exist",
				  "password": "somerandompassword"
				}`,
				`{
				  "taskId": "09b60ba8-9e76-469c-8ea1-0b7596574afd",
				  "commandType": "aclUserCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T13:21:03.96463Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/09b60ba8-9e76-469c-8ea1-0b7596574afd",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/09b60ba8-9e76-469c-8ea1-0b7596574afd",
				`{
				  "taskId": "09b60ba8-9e76-469c-8ea1-0b7596574afd",
				  "commandType": "aclUserCreateRequest",
				  "status": "processing-error",
				  "description": "Task request failed during processing. See error information for failure details.",
				  "timestamp": "2023-06-21T13:21:08.203427Z",
				  "response": {
					"error": {
					  "type": "ACL_ROLE_DOES_NOT_EXISTS",
					  "status": "400 BAD_REQUEST",
					  "description": "ACL role associated with a user has to exist."
					}
				  },
				  "links": [
					{
					  "rel": "self",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/09b60ba8-9e76-469c-8ea1-0b7596574afd",
					  "type": "GET"
					}
				  ]
				}`,
			),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	_, err = subject.Users.Create(context.TODO(), users.CreateUserRequest{
		Name:     redis.String("ACL-user-example"),
		Role:     redis.String("Redis-role-example-that-does-not-exist"),
		Password: redis.String("somerandompassword"),
	})

	require.Error(t, err)
}

func TestCreateUser(t *testing.T) {
	expected := 212
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			postRequest(
				t,
				"/acl/users",
				`{
				  "name": "ACL-user-example",
				  "role": "ACL-role-example",
				  "password": "someRandom.pa55word"
				}`,
				`{
				  "taskId": "72f1e802-bd98-47b1-ab89-7c402cc09736",
				  "commandType": "aclUserCreateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T13:31:08.833253Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/72f1e802-bd98-47b1-ab89-7c402cc09736",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`,
			),
			getRequest(
				t,
				"/tasks/72f1e802-bd98-47b1-ab89-7c402cc09736",
				fmt.Sprintf(`{
				  "taskId": "72f1e802-bd98-47b1-ab89-7c402cc09736",
				  "commandType": "aclUserCreateRequest",
				  "status": "processing-completed",
				  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
				  "timestamp": "2023-06-21T13:31:14.296891Z",
				  "response": {
					"resourceId": %[1]d
				  },
				  "links": [
					{
					  "rel": "resource",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/%[1]d/acl/users",
					  "title": "getACLUserInformation",
					  "type": "GET"
					},
					{
					  "rel": "self",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/72f1e802-bd98-47b1-ab89-7c402cc09736",
					  "type": "GET"
					}
				  ]
				}`, expected)),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Users.Create(context.TODO(), users.CreateUserRequest{
		Name:     redis.String("ACL-user-example"),
		Role:     redis.String("ACL-role-example"),
		Password: redis.String("someRandom.pa55word"),
	})
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUpdateNonExistentUser(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/acl/users/40004",
				`{
				  "role": "ACL-role-example",
				  "password": "someRandom.pa55word"
				}`,
				`{
				  "taskId": "b549b6ef-14c2-4ff1-b5b9-9e84afc9a6aa",
				  "commandType": "aclUserUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T13:36:42.993222Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/b549b6ef-14c2-4ff1-b5b9-9e84afc9a6aa",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/b549b6ef-14c2-4ff1-b5b9-9e84afc9a6aa", `{
			  "taskId": "b549b6ef-14c2-4ff1-b5b9-9e84afc9a6aa",
			  "commandType": "aclUserUpdateRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T13:36:43.309522Z",
			  "response": {
				"error": {
				  "type": "ACL_USER_NOT_FOUND",
				  "status": "404 NOT_FOUND",
				  "description": "ACL user not found"
				}
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/b549b6ef-14c2-4ff1-b5b9-9e84afc9a6aa",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Users.Update(context.TODO(), 40004, users.UpdateUserRequest{
		Role:     redis.String("ACL-role-example"),
		Password: redis.String("someRandom.pa55word"),
	})

	assert.Equal(t, &internal.Error{
		Type:        redis.String("ACL_USER_NOT_FOUND"),
		Description: redis.String("ACL user not found"),
		Status:      redis.String("404 NOT_FOUND"),
	}, errors.Unwrap(err))
}

func TestUpdateBadUser(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/acl/users/40000",
				`{
					"role": "ACL-role-example",
				    "password": "I do not meet requirements"
				}`,
				`{
				  "taskId": "612d6943-22f3-4ba1-ae18-b906e3749b2d",
				  "commandType": "aclUserUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T13:40:12.738053Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/612d6943-22f3-4ba1-ae18-b906e3749b2d",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/612d6943-22f3-4ba1-ae18-b906e3749b2d", `{
			  "taskId": "612d6943-22f3-4ba1-ae18-b906e3749b2d",
			  "commandType": "aclUserUpdateRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T13:40:13.177847Z",
			  "response": {
				"error": {
				  "type": "ACL_USER_PASSWORD_NOT_VALID",
				  "status": "400 BAD_REQUEST",
				  "description": "ACL user password is not valid."
				},
				"additionalInfo": "cluster-user-password-wrong-pattern"
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/612d6943-22f3-4ba1-ae18-b906e3749b2d",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Users.Update(context.TODO(), 40000, users.UpdateUserRequest{
		Role:     redis.String("ACL-role-example"),
		Password: redis.String("I do not meet requirements"),
	})

	assert.Equal(t, &internal.Error{
		Type:        redis.String("ACL_USER_PASSWORD_NOT_VALID"),
		Description: redis.String("ACL user password is not valid."),
		Status:      redis.String("400 BAD_REQUEST"),
	}, errors.Unwrap(err))
}

func TestUpdateUser(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			putRequest(
				t,
				"/acl/users/20000",
				`{
					"role": "ACL-role-example",
				    "password": "someOther.pa55word"
				}`,
				`{
				  "taskId": "c54fbe6a-c67c-461f-b02f-7d4504fc222d",
				  "commandType": "aclUserUpdateRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T13:43:00.183997Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/c54fbe6a-c67c-461f-b02f-7d4504fc222d",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/c54fbe6a-c67c-461f-b02f-7d4504fc222d", `{
			  "taskId": "c54fbe6a-c67c-461f-b02f-7d4504fc222d",
			  "commandType": "aclUserUpdateRequest",
			  "status": "processing-completed",
			  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
			  "timestamp": "2023-06-21T13:43:05.626511Z",
			  "response": {
				"resourceId": 20000
			  },
			  "links": [
				{
				  "rel": "resource",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/1209/acl/users",
				  "title": "getACLUserInformation",
				  "type": "GET"
				},
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/c54fbe6a-c67c-461f-b02f-7d4504fc222d",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Users.Update(context.TODO(), 20000, users.UpdateUserRequest{
		Role:     redis.String("ACL-role-example"),
		Password: redis.String("someOther.pa55word"),
	})
	require.NoError(t, err)
}

func TestDeleteNonExistentUser(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			deleteRequest(
				t,
				"/acl/users/40004",
				`{
				  "taskId": "644744ad-4b94-4a3c-bf36-cca5891f3f35",
				  "commandType": "aclUserDeleteRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T13:46:36.013348Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/644744ad-4b94-4a3c-bf36-cca5891f3f35",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/644744ad-4b94-4a3c-bf36-cca5891f3f35", `{
			  "taskId": "644744ad-4b94-4a3c-bf36-cca5891f3f35",
			  "commandType": "aclUserDeleteRequest",
			  "status": "processing-error",
			  "description": "Task request failed during processing. See error information for failure details.",
			  "timestamp": "2023-06-21T13:46:36.397443Z",
			  "response": {
				"error": {
				  "type": "ACL_USER_NOT_FOUND",
				  "status": "404 NOT_FOUND",
				  "description": "ACL user not found"
				}
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/644744ad-4b94-4a3c-bf36-cca5891f3f35",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Users.Delete(context.TODO(), 40004)

	assert.Equal(t, &internal.Error{
		Type:        redis.String("ACL_USER_NOT_FOUND"),
		Description: redis.String("ACL user not found"),
		Status:      redis.String("404 NOT_FOUND"),
	}, errors.Unwrap(err))
}

func TestDeleteUser(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			deleteRequest(
				t,
				"/acl/users/20000",
				`{
				  "taskId": "7dd16a39-2c77-4659-8457-8b9cc8bbf5f4",
				  "commandType": "aclUserDeleteRequest",
				  "status": "received",
				  "description": "Task request received and is being queued for processing.",
				  "timestamp": "2023-06-21T13:47:37.819348Z",
				  "links": [
					{
					  "rel": "task",
					  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/7dd16a39-2c77-4659-8457-8b9cc8bbf5f4",
					  "title": "getTaskStatusUpdates",
					  "type": "GET"
					}
				  ]
				}`),
			getRequest(t, "/tasks/7dd16a39-2c77-4659-8457-8b9cc8bbf5f4", `{
			  "taskId": "7dd16a39-2c77-4659-8457-8b9cc8bbf5f4",
			  "commandType": "aclUserDeleteRequest",
			  "status": "processing-completed",
			  "description": "Request processing completed successfully and its resources are now being provisioned / de-provisioned.",
			  "timestamp": "2023-06-21T13:47:38.165523Z",
			  "response": {
				"resourceId": 20000
			  },
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/tasks/7dd16a39-2c77-4659-8457-8b9cc8bbf5f4",
				  "type": "GET"
				}
			  ]
			}`),
		))

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	err = subject.Users.Delete(context.TODO(), 20000)
	require.NoError(t, err)
}

func TestListUsers(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(t, "/acl/users", `{
			  "accountId": 53012,
			  "users": [
				{
				  "id": 24,
				  "name": "test-user",
				  "role": "test-role",
				  "links": []
				}
			  ],
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/users",
				  "type": "GET"
				}
			  ]
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Users.List(context.TODO())
	require.NoError(t, err)

	assert.ElementsMatch(t, []*users.GetUserResponse{
		{
			ID:   redis.Int(24),
			Name: redis.String("test-user"),
			Role: redis.String("test-role"),
		},
	}, actual)

}

func TestGetNonExistentUser(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequestWithStatus(t, "/acl/users/40004", 404, `{
			  "timestamp": "2023-06-21T13:51:50.571+0000",
			  "status": 404,
			  "error": "Not Found",
			  "message": "ACL user 40004 not found",
			  "path": "/v1/acl/users/40004"
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	_, err = subject.Users.Get(context.TODO(), 40004)
	require.Error(t, err)
}

func TestGetUser(t *testing.T) {
	server := httptest.NewServer(
		testServer(
			"key",
			"secret",
			getRequest(t, "/acl/users/20000", `{
			  "id": 20000,
			  "name": "test-user",
			  "role": "test-role",
			  "status": "pending",
			  "links": [
				{
				  "rel": "self",
				  "href": "https://api-cloudapi.qa.redislabs.com/v1/acl/users/20000",
				  "type": "GET"
				}
			  ]
			}`),
		),
	)

	subject, err := clientFromTestServer(server, "key", "secret")
	require.NoError(t, err)

	actual, err := subject.Users.Get(context.TODO(), 20000)
	require.NoError(t, err)

	assert.Equal(t, &users.GetUserResponse{
		ID:     redis.Int(20000),
		Name:   redis.String("test-user"),
		Role:   redis.String("test-role"),
		Status: redis.String(users.StatusPending),
	}, actual)
}
