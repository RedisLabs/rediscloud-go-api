package latest_backups

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	GetWithQuery(ctx context.Context, name, path string, query url.Values, responseBody interface{}) error
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
	return task, nil
}

func (a *API) GetFixed(ctx context.Context, subscription int, database int) (*LatestBackupStatus, error) {
	message := fmt.Sprintf("get latest backup information for database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("/fixed/subscriptions/%d/databases/%d/backup", subscription, database)
	task, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}
	return task, nil
}

func (a *API) GetActiveActive(ctx context.Context, subscription int, database int, region string) (*LatestBackupStatus, error) {
	message := fmt.Sprintf("get latest backup information for database %d in subscription %d and region %s", subscription, database, region)
	address := fmt.Sprintf("/subscriptions/%d/databases/%d/backup", subscription, database)

	q := map[string][]string{
		"regionName": {region},
	}

	var task internal.TaskResponse
	err := a.client.GetWithQuery(ctx, message, address, q, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for backup status request %d to complete", task.ID)

	taskResp, err := a.taskWaiter.WaitForTask(ctx, *task.ID)
	if err != nil {
		var iErr *internal.Error
		if errors.As(err, &iErr) && taskResp != nil {
			backupStatusTask, err := createLatestBackupStatusFromTask(taskResp)
			if err != nil {
				return nil, err
			}
			return backupStatusTask, nil
		}
		return nil, wrap404ErrorActiveActive(subscription, database, region,
			fmt.Errorf("failed to retrieve completed backup status %d: %w", task.ID, err))
	}

	return createLatestBackupStatusFromTask(taskResp)
}

func (a *API) get(ctx context.Context, message string, address string) (*LatestBackupStatus, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, message, address, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for backup status request %d to complete", task.ID)

	taskResp, err := a.taskWaiter.WaitForTask(ctx, *task.ID)
	if err != nil {
		var iErr *internal.Error
		if errors.As(err, &iErr) && taskResp != nil {
			backupStatusTask, err := createLatestBackupStatusFromTask(taskResp)
			if err != nil {
				return nil, err
			}
			return backupStatusTask, nil
		}
		return nil, fmt.Errorf("failed to retrieve completed backup status %d: %w", task.ID, err)
	}

	return createLatestBackupStatusFromTask(taskResp)
}

func wrap404Error(subId int, dbId int, err error) error {
	var httpErr *internal.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
		return &NotFound{subId: subId, dbId: dbId}
	}
	return err
}

func wrap404ErrorActiveActive(subId int, dbId int, region string, err error) error {
	var httpErr *internal.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
		return &NotFoundActiveActive{subId: subId, dbId: dbId, region: region}
	}
	return err
}

func createLatestBackupStatusFromTask(task *internal.Task) (*LatestBackupStatus, error) {
	latestBackupStatus := &LatestBackupStatus{}
	if task != nil {
		latestBackupStatus.CommandType = task.CommandType
		latestBackupStatus.Description = task.Description
		latestBackupStatus.Status = task.Status
		latestBackupStatus.ID = task.ID
		if task.Response != nil {
			latestBackupStatus.Response = &Response{
				ID: task.Response.ID,
			}
			if task.Response.Error != nil {
				latestBackupStatus.Response.Error = &Error{
					Type:        task.Response.Error.Type,
					Description: task.Response.Error.Description,
					Status:      task.Response.Error.Status,
				}
			}

			if task.Response.Resource != nil {
				latestBackupStatus.Response.Resource = &Resource{}
				err := json.Unmarshal(*task.Response.Resource, latestBackupStatus.Response.Resource)
				if err != nil {
					return nil, fmt.Errorf("failed to unmarshal task response: %w", err)
				}
			}
		}
	}
	return latestBackupStatus, nil
}
