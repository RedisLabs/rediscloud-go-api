package subscriptions

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

// Create will create a new subscription.
func (a *API) Create(ctx context.Context, subscription FixedSubscription) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, "create fixed subscription", "/fixed/subscriptions", subscription, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the fixed subscription", task)

	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// List will list all of the current account's fixed subscriptions.
func (a *API) List(ctx context.Context) ([]*FixedSubscription, error) {
	var response listFixedSubscriptionResponse
	err := a.client.Get(ctx, "list fixed subscriptions", "/fixed/subscriptions", &response)
	if err != nil {
		return nil, err
	}

	return response.FixedSubscriptions, nil
}

// Get will retrieve an existing fixed subscription.
func (a *API) Get(ctx context.Context, id int) (*FixedSubscription, error) {
	var response FixedSubscription
	err := a.client.Get(ctx, fmt.Sprintf("retrieve fixed subscription %d", id), fmt.Sprintf("/fixed/subscriptions/%d", id), &response)
	if err != nil {
		return nil, wrap404Error(id, err)
	}

	return &response, nil
}

// Update will make changes to an existing fixed subscription.
func (a *API) Update(ctx context.Context, id int, subscription FixedSubscription) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update fixed subscription %d", id), fmt.Sprintf("/fixed/subscriptions/%d", id), subscription, &task)
	if err != nil {
		return wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for task %s to finish updating the fixed subscription", task)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when updating fixed subscription %d: %w", id, err)
	}

	return nil
}

// Delete will destroy an existing subscription. All existing databases within the subscription should already be
// deleted, otherwise this function will fail.
func (a *API) Delete(ctx context.Context, id int) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("delete fixed subscription %d", id), fmt.Sprintf("/fixed/subscriptions/%d", id), &task)
	if err != nil {
		return wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for fixed subscription %d to finish being deleted", id)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

func wrap404Error(id int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{ID: id}
	}
	return err
}
