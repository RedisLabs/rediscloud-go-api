package subscriptions

import (
	"context"
	"fmt"

	"github.com/RedisLabs/rediscloud-go-api/redis"
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

type Api struct {
	client HttpClient
	task   Task
	logger Log
}

func NewApi(client HttpClient, task Task, logger Log) *Api {
	return &Api{client: client, task: task, logger: logger}
}

func (a *Api) Create(ctx context.Context, subscription CreateSubscription) (int, error) {
	var task taskResponse
	err := a.client.Post(ctx, "create subscription", "/subscriptions", subscription, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the subscription", task)

	id, err := a.task.WaitForResourceId(ctx, redis.StringValue(task.ID))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *Api) Delete(ctx context.Context, id int) error {
	var task taskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("delete subscription %d", id), fmt.Sprintf("/subscriptions/%d", id), &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for subscription %d to finish being deleted", id)

	err = a.task.Wait(ctx, redis.StringValue(task.ID))
	if err != nil {
		return err
	}

	return nil
}
