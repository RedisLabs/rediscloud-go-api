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
	Delete(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
}

type TaskWaiter interface {
	WaitForResourceId(ctx context.Context, id string) (int, error)
	WaitForResource(ctx context.Context, id string, resource interface{}) error
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
func (a *API) Create(ctx context.Context, subscription CreateSubscription) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, "create subscription", "/subscriptions", subscription, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the subscription", task)

	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// List will list all of the current account's subscriptions.
func (a *API) List(ctx context.Context) ([]*Subscription, error) {
	var response listSubscriptionResponse
	err := a.client.Get(ctx, "list subscriptions", "/subscriptions", &response)
	if err != nil {
		return nil, err
	}

	return response.Subscriptions, nil
}

// Get will retrieve an existing subscription.
func (a *API) Get(ctx context.Context, id int) (*Subscription, error) {
	var response Subscription
	err := a.client.Get(ctx, fmt.Sprintf("retrieve subscription %d", id), fmt.Sprintf("/subscriptions/%d", id), &response)
	if err != nil {
		return nil, wrap404Error(id, err)
	}

	return &response, nil
}

// Update will make changes to an existing subscription.
func (a *API) Update(ctx context.Context, id int, subscription UpdateSubscription) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update subscription %d", id), fmt.Sprintf("/subscriptions/%d", id), subscription, &task)
	if err != nil {
		return wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for task %s to finish updating the subscription", task)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when updating subscription %d: %w", id, err)
	}

	return nil
}

// Update will make changes to an existing subscription's CMKs.
func (a *API) UpdateCMKs(ctx context.Context, id int, subscriptionCMKs UpdateSubscriptionCMKs) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update subscription %d", id), fmt.Sprintf("/subscriptions/%d", id), subscriptionCMKs, &task)
	if err != nil {
		return wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for task %s to finish updating subscription %d", task, id)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when updating subscription %d: %w", id, err)
	}

	return nil
}

// Delete will destroy an existing subscription. All existing databases within the subscription should already be
// deleted, otherwise this function will fail.
func (a *API) Delete(ctx context.Context, id int) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("delete subscription %d", id), fmt.Sprintf("/subscriptions/%d", id), nil, &task)
	if err != nil {
		return wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for subscription %d to finish being deleted", id)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

// GetCIDRAllowlist retrieves the CIDR addresses that are allowed to access an endpoint for a database associated with
// a the subscription.
func (a *API) GetCIDRAllowlist(ctx context.Context, id int) (*CIDRAllowlist, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, fmt.Sprintf("get cidr for subscription %d", id), fmt.Sprintf("/subscriptions/%d/cidr", id), &task)
	if err != nil {
		return nil, wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for subscription %d CIDR allowlist to be retrieved", id)

	var response CIDRAllowlist
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateCIDRAllowlist modifies the CIDR addresses that are allowed to access an endpoint for a database associated with
// a the subscription.
func (a *API) UpdateCIDRAllowlist(ctx context.Context, id int, cidr UpdateCIDRAllowlist) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update cidr for subscription %d", id), fmt.Sprintf("/subscriptions/%d/cidr", id), cidr, &task)
	if err != nil {
		return wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for subscription %d CIDR allowlist to finish being updated", id)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

// ListVPCPeering retrieves the VPCs that have been peered to the subscription VPC.
func (a *API) ListVPCPeering(ctx context.Context, id int) ([]*VPCPeering, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, fmt.Sprintf("get peerings for subscription %d", id), fmt.Sprintf("/subscriptions/%d/peerings", id), &task)
	if err != nil {
		return nil, wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for subscription %d peering details to be retrieved", id)

	var peering listVpcPeering
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &peering)
	if err != nil {
		return nil, err
	}

	return peering.Peerings, nil
}

func (a *API) ListActiveActiveVPCPeering(ctx context.Context, id int) ([]*ActiveActiveVpcRegion, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, fmt.Sprintf("get peerings for subscription %d", id), fmt.Sprintf("/subscriptions/%d/regions/peerings/", id), &task)
	if err != nil {
		return nil, wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for subscription %d peering details to be retrieved", id)

	var peering listActiveActiveVpcPeering
	err = a.taskWaiter.WaitForResource(ctx, *task.ID, &peering)
	if err != nil {
		return nil, err
	}

	return peering.Regions, nil
}

// CreateVPCPeering creates a new VPC peering from the subscription VPC and returns the identifier of the VPC peering.
func (a *API) CreateVPCPeering(ctx context.Context, id int, create CreateVPCPeering) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, fmt.Sprintf("create peering for subscription %d", id), fmt.Sprintf("/subscriptions/%d/peerings", id), create, &task)
	if err != nil {
		return 0, wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for subscription %d peering details to be retrieved", id)

	id, err = a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *API) CreateActiveActiveVPCPeering(ctx context.Context, id int, create CreateActiveActiveVPCPeering) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, fmt.Sprintf("create peering for subscription %d", id), fmt.Sprintf("/subscriptions/%d/regions/peerings/", id), create, &task)
	if err != nil {
		return 0, wrap404Error(id, err)
	}

	a.logger.Printf("Waiting for subscription %d peering details to be retrieved", id)

	id, err = a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// DeleteVPCPeering destroys an existing VPC peering connection.
func (a *API) DeleteVPCPeering(ctx context.Context, subscription int, peering int) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("deleting peering %d for subscription %d", peering, subscription), fmt.Sprintf("/subscriptions/%d/peerings/%d", subscription, peering), nil, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for peering %d for subscription %d to be deleted", peering, subscription)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

func (a *API) DeleteActiveActiveVPCPeering(ctx context.Context, subscription int, peering int) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("deleting peering %d for subscription %d", peering, subscription), fmt.Sprintf("/subscriptions/%d/regions/peerings/%d", subscription, peering), nil, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for peering %d for subscription %d to be deleted", peering, subscription)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

func (a *API) ListActiveActiveRegions(ctx context.Context, subscription int) ([]*ActiveActiveRegion, error) {
	var response ListAASubscriptionRegionsResponse
	err := a.client.Get(ctx, "list regions", fmt.Sprintf("/subscriptions/%d/regions", subscription), &response)

	if err != nil {
		return nil, err
	}

	return response.Regions, nil
}

// GetRedisVersions retrieves the Redis database versions available for this subscription.
func (a *API) GetRedisVersions(ctx context.Context, subscription int) (*RedisVersions, error) {
	var redisVersions RedisVersions
	getRedisVersionsUrl := "/subscriptions/redis-versions?subscriptionId=%d"

	path := fmt.Sprintf(getRedisVersionsUrl, subscription)
	err := a.client.Get(
		ctx,
		fmt.Sprintf("get versions for subscription %d", subscription),
		path,
		&redisVersions,
	)

	if err != nil {
		return nil, wrap404Error(subscription, err)
	}
	return &redisVersions, nil
}

func wrap404Error(id int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{ID: id}
	}
	return err
}
