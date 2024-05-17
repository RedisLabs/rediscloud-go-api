package latest_imports

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
}

type TaskWaiter interface {
	WaitForTask(ctx context.Context, id string) (*internal.Task, error)
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

func (a *API) Get(ctx context.Context, subscription int, database int) (*LatestImportStatus, error) {
	message := fmt.Sprintf("get latest import information for database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("/subscriptions/%d/databases/%d/import", subscription, database)
	task, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}
	return NewLatestImportStatus(task), nil
}

func (a *API) GetFixed(ctx context.Context, subscription int, database int) (*LatestImportStatus, error) {
	message := fmt.Sprintf("get latest import information for database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("/fixed/subscriptions/%d/databases/%d/import", subscription, database)
	task, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}
	return NewLatestImportStatus(task), nil
}

func (a *API) get(ctx context.Context, message string, address string) (*internal.Task, error) {
	var taskResponse internal.TaskResponse
	err := a.client.Get(ctx, message, address, &taskResponse)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for backup status request %d to complete", taskResponse.ID)

	return a.taskWaiter.WaitForTask(ctx, *taskResponse.ID)
}

func wrap404Error(subId int, dbId int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{subId: subId, dbId: dbId}
	}
	return err
}
