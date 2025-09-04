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
	Delete(ctx context.Context, name, path string, responseBody interface{}) error
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

func (a *API) GetPrivateLink(ctx context.Context, subscription int) (*PrivateLink, error) {
	message := fmt.Sprintf("get private link for subscription %d", subscription)
	path := fmt.Sprintf("/subscriptions/%d/private-link", subscription)
	task, err := a.getLink(ctx, message, path)
	if err != nil {
		return nil, wrap404Error(subscription, err)
	}
	return task, nil
}

func (a *API) getLink(ctx context.Context, message string, path string) (*PrivateLink, error) {
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

func wrap404ErrorActiveActive(subId int, regionId int, err error) error {
	var e *internal.HTTPError
	if errors.As(err, &e) && e.StatusCode == http.StatusNotFound {
		return &NotFoundActiveActive{subscriptionID: subId, regionID: regionId}
	}
	var v *internal.Error
	if errors.As(err, &v) && v.StatusCode() == strconv.Itoa(http.StatusNotFound) {
		return &NotFoundActiveActive{subscriptionID: subId, regionID: regionId}
	}
	return err
}
