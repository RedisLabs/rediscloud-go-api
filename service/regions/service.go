package regions

import (
	"context"
	"fmt"
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
	WaitForResource(ctx context.Context, id string, resource interface{}) error
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

// List will list all of a given subscription's active-active regions.
func (a API) List(ctx context.Context, subId int) (*Regions, error) {
	var response Regions
	err := a.client.Get(ctx, "list regions", fmt.Sprintf("/subscriptions/%d/regions", subId), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
