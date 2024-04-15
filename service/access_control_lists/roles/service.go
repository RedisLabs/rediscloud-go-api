package roles

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

// List will list all of the current account's roles.
func (a API) List(ctx context.Context) ([]*GetRoleResponse, error) {
	var response ListRolesResponse
	err := a.client.Get(ctx, "list roles", "/acl/roles", &response)
	if err != nil {
		return nil, err
	}

	return response.Roles, nil
}

// Get has to use the List behaviour to simulate getById
func (a API) Get(ctx context.Context, id int) (*GetRoleResponse, error) {
	roles, err := a.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		if id == *role.ID {
			return role, nil
		}
	}

	return nil, &NotFound{ID: id}
}

// Create will create a new role and return the identifier of the role.
func (a *API) Create(ctx context.Context, role CreateRoleRequest) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, "create role", "/acl/roles", role, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the role", task)

	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, fmt.Errorf("failed when creating role %d: %w", id, err)
	}

	return id, nil
}

// Update will make changes to an existing role.
func (a *API) Update(ctx context.Context, id int, role CreateRoleRequest) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update role %d", id), fmt.Sprintf("/acl/roles/%d", id), role, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish updating the role", task)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when updating role %d: %w", id, err)
	}

	return nil
}

// Delete will destroy an existing role.
func (a *API) Delete(ctx context.Context, id int) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("delete role %d", id), fmt.Sprintf("/acl/roles/%d", id), &task)
	if err != nil {
		return wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for role %d to finish being deleted", id)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when deleting role %d: %w", id, err)
	}

	return nil
}

type NotFound struct {
	ID int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("role %d not found", f.ID)
}

func wrap404Error(id int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{ID: id}
	}
	return err
}
