package latest_backups

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

func (a *API) Get(ctx context.Context, subscription int, database int) (*LatestBackupStatus, error) {
	message := fmt.Sprintf("get latest backup information for database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("/subscriptions/%d/databases/%d/backup", subscription, database)
	task, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}
	return NewLatestBackupStatus(task), nil
}

func (a *API) GetActiveActive(ctx context.Context, subscription int, database int, region string) (*LatestBackupStatus, error) {
	message := fmt.Sprintf("get latest backup information for database %d in subscription %d and region %s", subscription, database, region)
	address := fmt.Sprintf("/subscriptions/%d/databases/%d/backup?regionName=%s", subscription, database, region)
	task, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404ErrorActiveActive(subscription, database, region, err)
	}
	return NewLatestBackupStatus(task), nil
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

func wrap404ErrorActiveActive(subId int, dbId int, region string, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFoundActiveActive{subId: subId, dbId: dbId, region: region}
	}
	return err
}
