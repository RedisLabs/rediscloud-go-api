package cloud_accounts

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
	Wait(ctx context.Context, id string) error
}

type Api struct {
	client HttpClient
	task   Task
	logger Log
}

func NewApi(client HttpClient, task Task, logger Log) *Api {
	return &Api{client: client, task: task, logger: logger}
}

// Create will create a new Cloud Account and return the identifier of the new account.
func (a *Api) Create(ctx context.Context, account CreateCloudAccount) (int, error) {
	var response taskResponse
	if err := a.client.Post(ctx, "cloud account", "/cloud-accounts", account, &response); err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for task %s to finish creating the cloud account", response.TaskId)

	id, err := a.task.WaitForResourceId(ctx, response.TaskId)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Get will retrieve an existing Cloud Account.
func (a *Api) Get(ctx context.Context, id int) (*CloudAccount, error) {
	var response CloudAccount
	if err := a.client.Get(ctx, fmt.Sprintf("retrieve cloud account %d", id), fmt.Sprintf("/cloud-accounts/%d", id), &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Update will update certain values of an existing Cloud Account.
func (a *Api) Update(ctx context.Context, id int, account UpdateCloudAccount) error {
	var response taskResponse
	if err := a.client.Put(ctx, fmt.Sprintf("update cloud account %d", id), fmt.Sprintf("/cloud-accounts/%d", id), account, &response); err != nil {
		return err
	}

	a.logger.Printf("Waiting for cloud account %d to finish being updated", id)

	err := a.task.Wait(ctx, response.TaskId)
	if err != nil {
		return fmt.Errorf("failed when updating account %d: %w", id, err)
	}

	return nil
}

// Delete will delete an existing Cloud Account.
func (a *Api) Delete(ctx context.Context, id int) error {
	var response taskResponse
	if err := a.client.Delete(ctx, fmt.Sprintf("delete cloud account %d", id), fmt.Sprintf("/cloud-accounts/%d", id), &response); err != nil {
		return err
	}

	a.logger.Printf("Waiting for cloud account %d to finish being deleted", id)

	if err := a.task.Wait(ctx, response.TaskId); err != nil {
		return fmt.Errorf("failed when deleting account %d: %w", id, err)
	}

	return nil
}
