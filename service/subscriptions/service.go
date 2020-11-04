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

func (a *API) Create(ctx context.Context, subscription CreateSubscription) (int, error) {
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

func (a *API) Get(ctx context.Context, id int) (*Subscription, error) {
	var response Subscription
	err := a.client.Get(ctx, fmt.Sprintf("retrieve subscription %d", id), fmt.Sprintf("/subscriptions/%d", id), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *API) Update(ctx context.Context, id int, subscription UpdateSubscription) error {
	var task taskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update subscription %d", id), fmt.Sprintf("/subscriptions/%d", id), subscription, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish creating the subscription", task)

	err = a.task.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when updating subscription %d: %w", id, err)
	}

	return nil
}

func (a *API) Delete(ctx context.Context, id int) error {
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

func (a *API) GetCIDRWhitelist(ctx context.Context, id int) (*CIDRWhitelist, error) {
	var task taskResponse
	err := a.client.Get(ctx, fmt.Sprintf("get cidr for subscription %d", id), fmt.Sprintf("/subscriptions/%d/cidr", id), &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for subscription %d CIDR whitelist to be retrieved", id)

	var response CIDRWhitelist
	err = a.task.WaitForResource(ctx, *task.ID, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (a *API) UpdateCIDRWhitelist(ctx context.Context, id int, cidr UpdateCIDRWhitelist) error {
	var task taskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update cidr for subscription %d", id), fmt.Sprintf("/subscriptions/%d/cidr", id), cidr, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for subscription %d CIDRs to finish being updated", id)

	err = a.task.Wait(ctx, *task.ID)
	if err != nil {
		return err
	}

	return nil
}

func (a *API) ListVPCPeering(ctx context.Context, id int) ([]*VPCPeering, error) {
	var task taskResponse
	err := a.client.Get(ctx, fmt.Sprintf("get peerings for subscription %d", id), fmt.Sprintf("/subscriptions/%d/peerings", id), &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for subscription %d peering details to be retrieved", id)

	var peering []*VPCPeering
	err = a.task.WaitForResource(ctx, *task.ID, &peering)
	if err != nil {
		return nil, err
	}

	return peering, nil
}

func (a *API) CreateVPCPeering(ctx context.Context, id int, create CreateVPCPeering) (int, error) {
	var task taskResponse
	err := a.client.Post(ctx, fmt.Sprintf("create peering for subscription %d", id), fmt.Sprintf("/subscriptions/%d/peerings", id), create, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for subscription %d peering details to be retrieved", id)

	id, err = a.task.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (a *API) DeleteVPCPeering(ctx context.Context, subscription int, peering int) error {
	var task taskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("deleting peering %d for subscription %d", peering, subscription), fmt.Sprintf("/subscriptions/%d/peerings/%d", subscription, peering), &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for peering %d for subscription %d to be deleted", peering, subscription)

	err = a.task.Wait(ctx, *task.ID)
	if err != nil {
		return err
	}

	return nil
}
