package roles

import (
	"context"
	"fmt"
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

type Task interface {
	WaitForResourceId(ctx context.Context, id string) (int, error)
	Wait(ctx context.Context, id string) error
}

type API struct {
	client HttpClient
	task   Task
	logger Log
}

func NewAPI(client HttpClient, task Task, logger Log) *API {
	return &API{client: client, task: task, logger: logger}
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

// No getById

// Create will create a new role and return the identifier of the role.
func (a *API) Create(ctx context.Context, role CreateRoleRequest) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, "create role", "/acl/roles", role, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the role", task)

	id, err := a.task.WaitForResourceId(ctx, *task.ID)
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

	err = a.task.Wait(ctx, *task.ID)
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
		return internal.Wrap404Error(id, "role", err)
	}

	a.logger.Printf("Waiting for role %d to finish being deleted", id)

	err = a.task.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when deleting role %d: %w", id, err)
	}

	return nil
}
