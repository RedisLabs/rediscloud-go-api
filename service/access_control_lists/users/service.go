package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type Log interface {
	Printf(format string, args ...interface{})
}

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Delete(ctx context.Context, name, path string, responseBody interface{}) error
}

type TaskWaiter interface {
	WaitForResourceId(ctx context.Context, id string) (int, error)
	Wait(ctx context.Context, id string) error
}

type API struct {
	client     HttpClient
	taskWaiter TaskWaiter
	logger     Log
}

func NewAPI(client HttpClient, taskWaiter TaskWaiter, logger Log) *API {
	return &API{client: client, taskWaiter: taskWaiter, logger: logger}
}

// List will list all of the current account's users.
func (a *API) List(ctx context.Context) ([]*GetUserResponse, error) {
	var response ListUsersResponse
	err := a.client.Get(ctx, "list users", "/acl/users", &response)
	if err != nil {
		return nil, err
	}

	return response.Users, nil
}

// Get will retrieve an existing user.
func (a *API) Get(ctx context.Context, id int) (*GetUserResponse, error) {
	var response GetUserResponse
	err := a.client.Get(ctx, fmt.Sprintf("get user %d", id), fmt.Sprintf("/acl/users/%d", id), &response)
	if err != nil {
		return nil, wrap404Error(id, err)
	}

	return &response, nil
}

// Create will create a new user and return the identifier of the user.
func (a *API) Create(ctx context.Context, user CreateUserRequest) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, "create user", "/acl/users", user, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the user", task)

	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Update will make changes to an existing user.
func (a *API) Update(ctx context.Context, id int, user UpdateUserRequest) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update user %d", id), fmt.Sprintf("/acl/users/%d", id), user, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish updating the user", task)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when updating user %d: %w", id, err)
	}

	return nil
}

// Delete will destroy an existing user.
func (a *API) Delete(ctx context.Context, id int) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("delete user %d", id), fmt.Sprintf("/acl/users/%d", id), &task)
	if err != nil {
		return wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for user %d to finish being deleted", id)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when deleting user %d: %w", id, err)
	}

	return nil
}

type NotFound struct {
	ID int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("user %d not found", f.ID)
}

func wrap404Error(id int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{ID: id}
	}
	return err
}
