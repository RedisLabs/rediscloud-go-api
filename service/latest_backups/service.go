package latest_backups

import (
	"context"
	"fmt"
	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/service"
	"net/http"
)

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
}

type TaskWaiter interface {
	Wait(ctx context.Context, id string) (*service.Task, error)
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

func (a *API) Get(ctx context.Context, subscription int, database int) (*service.Task, error) {
	var taskResponse service.TaskResponse
	message := fmt.Sprintf("get latest backup information for database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("/subscriptions/%d/databases/%d/backup", subscription, database)
	err := a.client.Get(ctx, message, address, &taskResponse)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for backup status request %d to complete", taskResponse.ID)

	completedTask, err := a.taskWaiter.Wait(ctx, *taskResponse.ID)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}
	return completedTask, nil
}

func wrap404Error(subId int, dbId int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{subId: subId, dbId: dbId}
	}
	return err
}
