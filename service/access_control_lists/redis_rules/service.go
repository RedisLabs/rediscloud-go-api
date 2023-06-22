package redis_rules

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

// List will list all of the current account's redisRules.
func (a API) List(ctx context.Context) ([]*GetRedisRuleResponse, error) {
	var response ListRedisRulesResponse
	err := a.client.Get(ctx, "list redisRules", "/acl/redisRules", &response)
	if err != nil {
		return nil, err
	}

	return response.RedisRules, nil
}

// No getById

// Create will create a new redisRule and return the identifier of the redisRule.
func (a *API) Create(ctx context.Context, redisRule CreateRedisRuleRequest) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, "create redisRule", "/acl/redisRules", redisRule, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the redisRule", task)

	id, err := a.task.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, fmt.Errorf("failed when creating redisRule %d: %w", id, err)
	}

	return id, nil
}

// Update will make changes to an existing redisRule.
func (a *API) Update(ctx context.Context, id int, redisRule CreateRedisRuleRequest) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update redisRule %d", id), fmt.Sprintf("/acl/redisRules/%d", id), redisRule, &task)
	if err != nil {
		return internal.Wrap404Error(id, "redisRule", err)
	}

	a.logger.Printf("Waiting for task %s to finish updating the redisRule", task)

	err = a.task.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when updating redisRule %d: %w", id, err)
	}

	return nil
}

// Delete will destroy an existing redisRule.
func (a *API) Delete(ctx context.Context, id int) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("delete redisRule %d", id), fmt.Sprintf("/acl/redisRules/%d", id), &task)
	if err != nil {
		return internal.Wrap404Error(id, "redisRule", err)
	}

	a.logger.Printf("Waiting for redisRule %d to finish being deleted", id)

	err = a.task.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when deleting redisRule %d: %w", id, err)
	}

	return nil
}
