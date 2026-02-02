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

// CreatePrivateLink will create a new PrivateLink.
func (a *API) CreatePrivateLink(ctx context.Context, subscriptionId int, privateLink CreatePrivateLink) error {
	message := fmt.Sprintf("create privatelink for subscription %d", subscriptionId)
	path := fmt.Sprintf("/subscriptions/%d/private-link", subscriptionId)
	err := a.create(ctx, message, path, privateLink)
	if err != nil {
		return wrap404Error(subscriptionId, err)
	}
	return nil
}

// GetPrivateLink will get a new PrivateLink.
func (a *API) GetPrivateLink(ctx context.Context, subscription int) (*PrivateLink, error) {
	message := fmt.Sprintf("get private link for subscription %d", subscription)
	path := fmt.Sprintf("/subscriptions/%d/private-link", subscription)
	task, err := a.get(ctx, message, path)
	if err != nil {
		if errors.Is(err, errEmptyResponse) {
			return nil, &NotFound{subscriptionID: subscription}
		}
		return nil, wrap404Error(subscription, err)
	}
	return task, nil
}

// GetPrivateLinkEndpointScript will get the script for an endpoint.
func (a *API) GetPrivateLinkEndpointScript(ctx context.Context, subscriptionId int) (*PrivateLinkEndpointScript, error) {
	message := fmt.Sprintf("get private link for subscription %d", subscriptionId)
	path := fmt.Sprintf("/subscriptions/%d/private-link/endpoint-script?includeTerraformAwsScript=true", subscriptionId)
	task, err := a.getScript(ctx, message, path)
	if err != nil {
		return nil, wrap404Error(subscriptionId, err)
	}
	return task, nil
}

// CreatePrincipal will add a principal to a PrivateLink.
func (a *API) CreatePrincipal(ctx context.Context, subscriptionId int, principal CreatePrivateLinkPrincipal) error {
	message := fmt.Sprintf("create principal %s for subscription %d", *principal.Principal, subscriptionId)
	path := fmt.Sprintf("/subscriptions/%d/private-link/principals", subscriptionId)

	err := a.create(ctx, message, path, principal)
	if err != nil {
		return wrap404Error(subscriptionId, err)
	}
	return nil
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

// DeletePrivateLink will delete a PrivateLink for a subscription.
// This marks the PrivateLink record as deleted but does not remove the actual AWS RL resources.
func (a *API) DeletePrivateLink(ctx context.Context, subscriptionId int) error {
	message := fmt.Sprintf("delete privatelink for subscription %d", subscriptionId)
	path := fmt.Sprintf("/subscriptions/%d/private-link", subscriptionId)

	err := a.delete(ctx, message, path, nil, nil)
	if err != nil {
		return wrap404Error(subscriptionId, err)
	}
	return nil
}

// CreateActiveActivePrivateLink will create a new active active PrivateLink.
func (a *API) CreateActiveActivePrivateLink(ctx context.Context, subscriptionId int, regionId int, privateLink CreatePrivateLink) error {
	message := fmt.Sprintf("create active active PrivateLink for subscription %d", subscriptionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-link", subscriptionId, regionId)
	err := a.create(ctx, message, path, privateLink)
	if err != nil {
		return wrap404Error(subscriptionId, err)
	}
	return nil
}

// GetActiveActivePrivateLink will get a new active active PrivateLink.
func (a *API) GetActiveActivePrivateLink(ctx context.Context, subscription int, regionId int) (*PrivateLink, error) {
	message := fmt.Sprintf("get active active PrivateLink for subscription %d", subscription)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-link", subscription, regionId)
	task, err := a.get(ctx, message, path)
	if err != nil {
		if errors.Is(err, errEmptyResponse) {
			return nil, &NotFoundActiveActive{subscriptionID: subscription, regionID: regionId}
		}
		return nil, wrap404Error(subscription, err)
	}
	return task, nil
}

// GetPrivateLinkEndpointScript will get the script for an endpoint.
func (a *API) GetActiveActivePrivateLinkEndpointScript(ctx context.Context, subscription int, regionId int) (*PrivateLinkEndpointScript, error) {
	message := fmt.Sprintf("get private link for subscription %d", subscription)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-link/endpoint-script?includeTerraformAwsScript=true", subscription, regionId)
	task, err := a.getScript(ctx, message, path)
	if err != nil {
		return nil, wrap404Error(subscription, err)
	}
	return task, nil
}

// CreateActiveActivePrincipal will add a principal to an active active PrivateLink.
func (a *API) CreateActiveActivePrincipal(ctx context.Context, subscriptionId int, regionId int, principal CreatePrivateLinkPrincipal) error {
	message := fmt.Sprintf("create principal %s for subscription %d", *principal.Principal, subscriptionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-link/principals", subscriptionId, regionId)

	err := a.create(ctx, message, path, principal)
	if err != nil {
		return wrap404Error(subscriptionId, err)
	}
	return nil
}

// DeleteActiveActivePrincipal will remove a principal from an active active PrivateLink.
func (a *API) DeleteActiveActivePrincipal(ctx context.Context, subscriptionId int, regionId int, principal string) error {
	message := fmt.Sprintf("delete principal %s for subscription %d", principal, subscriptionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-link/principals", subscriptionId, regionId)

	requestBody := map[string]interface{}{
		"principal": principal,
	}

	err := a.delete(ctx, message, path, requestBody, nil)
	if err != nil {
		return wrap404Error(subscriptionId, err)
	}
	return nil
}

// DeleteActiveActivePrivateLink will delete an Active-Active PrivateLink for a subscription region.
// This marks the PrivateLink record as deleted but does not remove the actual AWS RL resources.
func (a *API) DeleteActiveActivePrivateLink(ctx context.Context, subscriptionId int, regionId int) error {
	message := fmt.Sprintf("delete active active privatelink for subscription %d region %d", subscriptionId, regionId)
	path := fmt.Sprintf("/subscriptions/%d/regions/%d/private-link", subscriptionId, regionId)

	err := a.delete(ctx, message, path, nil, nil)
	if err != nil {
		return wrap404Error(subscriptionId, err)
	}
	return nil
}

func (a *API) create(ctx context.Context, message string, path string, requestBody interface{}) error {
	var task internal.TaskResponse
	err := a.client.Post(ctx, message, path, requestBody, &task)
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

func (a *API) get(ctx context.Context, message string, path string) (*PrivateLink, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, message, path, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for PrivateLink get request %d to complete", task.ID)

	var response PrivateLink
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &response)
	if err != nil {
		return nil, err
	}

	// API returns empty resource (e.g., {"links": []}) when privatelink doesn't exist
	// instead of a proper 404. Detect this and return sentinel error for callers to handle.
	if response.Status == nil && response.ShareName == nil {
		return nil, errEmptyResponse
	}

	return &response, nil
}

func (a *API) getScript(ctx context.Context, message string, path string) (*PrivateLinkEndpointScript, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, message, path, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for PrivateLink script get request %s to complete", *task.ID)

	var response PrivateLinkEndpointScript
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *API) delete(ctx context.Context, message string, path string, requestBody interface{}, responseBody interface{}) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, message, path, requestBody, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish deleting the PrivateLink", task)

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
