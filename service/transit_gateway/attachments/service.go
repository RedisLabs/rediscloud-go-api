package attachments

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	Post(ctx context.Context, name, path string, requsetBody interface{}, responseBody interface{}) error
	Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Delete(ctx context.Context, name, path string, responseBody interface{}) error
}

type TaskWaiter interface {
	WaitForResourceId(ctx context.Context, id string) (int, error)
	Wait(ctx context.Context, id string) error
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

func (a *API) Get(ctx context.Context, subscription int) (*GetAttachmentsTask, error) {
	message := fmt.Sprintf("get TGw attachments for subscription %d", subscription)
	address := fmt.Sprintf("/subscriptions/%d/transitGateways", subscription)
	task, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404Error(subscription, err)
	}
	return task, nil
}

func (a *API) GetActiveActive(ctx context.Context, subscription int, regionId int) (*GetAttachmentsTask, error) {
	message := fmt.Sprintf("get TGw attachments for subscription %d in region %d", subscription, regionId)
	address := fmt.Sprintf("/subscriptions/%d/regions/%d/transitGateways", subscription, regionId)
	task, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return task, nil
}

func (a *API) Create(ctx context.Context, subscription int, tgwId int) (int, error) {
	message := fmt.Sprintf("create TGw attachment for subscription %d", subscription)
	address := fmt.Sprintf("/subscriptions/%d/transitGateways/%d/attachment", subscription, tgwId)
	resourceId, err := a.create(ctx, message, address)
	if err != nil {
		return 0, wrap404Error(subscription, err)
	}
	return resourceId, nil
}

func (a *API) CreateActiveActive(ctx context.Context, subscription int, regionId int, tgwId int) (int, error) {
	message := fmt.Sprintf("create TGw attachment for subscription %d in region %d", subscription, regionId)
	address := fmt.Sprintf("/subscriptions/%d/regions/%d/transitGateways/%d/attachment", subscription, regionId, tgwId)
	resourceId, err := a.create(ctx, message, address)
	if err != nil {
		return 0, wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return resourceId, nil
}

func (a *API) Update(ctx context.Context, subscription int, tgwId int, cidrs []*string) error {
	message := fmt.Sprintf("update TGw attachment %d for subscription %d", tgwId, subscription)
	address := fmt.Sprintf("/subscriptions/%d/transitGateways/%d/attachment", subscription, tgwId)
	err := a.update(ctx, message, address, cidrs)
	if err != nil {
		return wrap404Error(subscription, err)
	}
	return nil
}

func (a *API) UpdateActiveActive(ctx context.Context, subscription int, regionId int, tgwId int, cidrs []*string) error {
	message := fmt.Sprintf("update TGw attachment %d for subscription %d in region %d", tgwId, subscription, regionId)
	address := fmt.Sprintf("/subscriptions/%d/regions/%d/transitGateways/%d/attachment", subscription, regionId, tgwId)
	err := a.update(ctx, message, address, cidrs)
	if err != nil {
		return wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return nil
}

func (a *API) Delete(ctx context.Context, subscription int, tgwId int) error {
	message := fmt.Sprintf("delete TGw attachment %d for subscription %d", tgwId, subscription)
	address := fmt.Sprintf("/subscriptions/%d/transitGateways/%d/attachment", subscription, tgwId)
	err := a.delete(ctx, message, address)
	if err != nil {
		return wrap404Error(subscription, err)
	}
	return nil
}

func (a *API) DeleteActiveActive(ctx context.Context, subscription int, regionId int, tgwId int) error {
	message := fmt.Sprintf("delete TGw attachment %d for subscription %d in region %d", tgwId, subscription, regionId)
	address := fmt.Sprintf("/subscriptions/%d/regions/%d/transitGateways/%d/attachment", subscription, regionId, tgwId)
	err := a.delete(ctx, message, address)
	if err != nil {
		return wrap404ErrorActiveActive(subscription, regionId, err)
	}
	return nil
}

func (a *API) get(ctx context.Context, message string, address string) (*GetAttachmentsTask, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, message, address, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for tgwGetRequest %d to complete", task.ID)

	err = a.taskWaiter.Wait(ctx, *task.ID)

	a.logger.Printf("tgwGetRequest %d completed, possibly with error", task.ID, err)

	var getAttachmentsTask *GetAttachmentsTask
	err = a.client.Get(ctx,
		fmt.Sprintf("retrieve completed tgwGetRequest task %d", task.ID),
		"/tasks/"+*task.ID,
		&getAttachmentsTask,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve completed tgwGetRequest %d: %w", task.ID, err)
	}

	return getAttachmentsTask, nil
}

func (a *API) create(ctx context.Context, message string, address string) (int, error) {
	var task internal.TaskResponse
	// TODO Assuming nil is an allowed body
	err := a.client.Post(ctx, message, address, nil, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the TGw attachment", task)

	// TODO Assuming the ID at task.Response.ID is what we want?
	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, fmt.Errorf("failed when creating TGw attachment %d: %w", id, err)
	}

	return id, nil
}

func (a *API) update(ctx context.Context, message string, address string, cidrs []*string) error {
	var task internal.TaskResponse
	// TODO Assuming this request body ([]*string) is acceptable and parsed correctly
	err := a.client.Put(ctx, message, address, cidrs, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish updating the TGw attachment", task)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when updating TGw attachment %w", err)
	}

	return nil
}

func (a *API) delete(ctx context.Context, message string, address string) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, message, address, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for task %s to finish deleting the TGw attachment", task)

	err = a.taskWaiter.Wait(ctx, *task.ID)
	if err != nil {
		return fmt.Errorf("failed when deleting TGw attachment %w", err)
	}

	return nil
}

func wrap404Error(subId int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{subId: subId}
	}
	return err
}

func wrap404ErrorActiveActive(subId int, regionId int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFoundActiveActive{subId: subId, regionId: regionId}
	}
	return err
}
