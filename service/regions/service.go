package regions

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
	DeleteWithQuery(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
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

// Create will create a new subscription.
func (a *API) Create(ctx context.Context, subId int, region CreateRegion) (int, error) {
	var task taskResponse
	err := a.client.Post(ctx, "create subscription", fmt.Sprintf("/subscriptions/%d/regions", subId), region, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the subscription", task)

	id, err := a.task.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
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

func (a *API) DeleteWithQuery(ctx context.Context, id int, regions DeleteRegions) error {
	var task taskResponse
	err := a.client.DeleteWithQuery(ctx, fmt.Sprintf("delete region %d", id), fmt.Sprintf("/subscriptions/%d/regions/", id), regions, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for region %d to finish being deleted", id)

	err = a.task.Wait(ctx, *task.ID)
	if err != nil {
		return err
	}

	return nil
}

type NotFound struct {
	id int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("subscription %d not found", f.id)
}

func wrap404Error(id int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{id: id}
	}
	return err
}
