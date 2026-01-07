package maintenance

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type Log interface {
	Printf(format string, args ...interface{})
}

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
}

type TaskWaiter interface {
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

// Get will retrieve a subscription's maintenance detail
func (a *API) Get(ctx context.Context, subscription int) (*Maintenance, error) {
	var m Maintenance
	err := a.client.Get(ctx, fmt.Sprintf("get maintenance for subscription %d", subscription), fmt.Sprintf("/subscriptions/%d/maintenance-windows", subscription), &m)
	if err != nil {
		return nil, wrap404Error(subscription, err)
	}

	return &m, nil
}

// Update will update a subscription's maintenance detail
func (a *API) Update(ctx context.Context, subscription int, m Maintenance) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update maintenance for subscription %d", subscription), fmt.Sprintf("/subscriptions/%d/maintenance-windows", subscription), m, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for fixed database %d for subscription %d to finish being updated", subscription)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

func wrap404Error(subId int, err error) error {
	var httpErr *internal.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
		return &NotFound{subId: subId}
	}
	return err
}
