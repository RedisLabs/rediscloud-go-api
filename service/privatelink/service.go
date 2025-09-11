package privatelink

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	GetWithQuery(ctx context.Context, name, path string, query url.Values, responseBody interface{}) error
	Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Delete(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
}

type TaskWaiter interface {
	WaitForResourceId(ctx context.Context, id string) (int, error)
	Wait(ctx context.Context, id string) error
	WaitForResource(ctx context.Context, id string, resource interface{}) error
}

type Log interface {
	Printf(format string, args ...interface{})
}

type API struct {
	client     HttpClient
	taskWaiter TaskWaiter
	logger     Log
}

func NewAPI(client HttpClient, taskWaiter TaskWaiter, logger Log) *API {
	return &API{client: client, taskWaiter: taskWaiter, logger: logger}
}

// // CreatePrivateLink will create a new PrivateLink.
func (a *API) CreatePrivateLink(ctx context.Context, subscriptionId int, privateLink CreatePrivateLink) error {
	message := fmt.Sprintf("create privatelink for subscription %d", subscriptionId)
	path := fmt.Sprintf("/subscriptions/%d/private-link", subscriptionId)
	err := a.create(ctx, message, path, privateLink)
	if err != nil {
		return wrap404Error(subscriptionId, err)
	}
	return nil
}

func (a *API) create(ctx context.Context, message string, path string, link CreatePrivateLink) error {
	var task internal.TaskResponse
	err := a.client.Post(ctx, message, path, nil, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish creating the PrivateLink", task)

	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when creating PrivateLink %d: %w", id, err)
	}

	return nil
}

// GetPrivateLink will get a new PrivateLink.
func (a *API) GetPrivateLink(ctx context.Context, subscription int) (*PrivateLink, error) {
	message := fmt.Sprintf("get private link for subscription %d", subscription)
	path := fmt.Sprintf("/subscriptions/%d/private-link", subscription)
	task, err := a.get(ctx, message, path)
	if err != nil {
		return nil, wrap404Error(subscription, err)
	}
	return task, nil
}

func (a *API) get(ctx context.Context, message string, path string) (*PrivateLink, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, message, path, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for privatelink request %d to complete", task.ID)

	var response PrivateLink
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeletePrincipal will remove a principal from a PrivateLink.
func (a *API) DeletePrincipal(ctx context.Context, subscriptionId int, principal string) error {
	message := fmt.Sprintf("delete principal %s for subscription %d", principal, subscriptionId)
	path := fmt.Sprintf("/subscriptions/%d/private-link/principals", subscriptionId)

	requestBody := map[string]interface{}{
		"principal": principal,
	}

	err := a.delete(ctx, message, path, requestBody, nil)
	if err != nil {
		return wrap404Error(subscriptionId, err)
	}
	return nil
}

func (a *API) delete(ctx context.Context, message string, path string, requestBody interface{}, responseBody interface{}) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, message, path, requestBody, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish deleting", task)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when deleting PrivateLink %w", err)
	}

	return nil
}

func wrap404Error(subId int, err error) error {
	var e *internal.HTTPError
	if errors.As(err, &e) && e.StatusCode == http.StatusNotFound {
		return &NotFound{subscriptionID: subId}
	}
	var v *internal.Error
	if errors.As(err, &v) && v.StatusCode() == strconv.Itoa(http.StatusNotFound) {
		return &NotFound{subscriptionID: subId}
	}
	return err
}
