package latest_imports

import (
	"context"
	"encoding/json"
	"errors"
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
	return task, nil
}

func (a *API) GetFixed(ctx context.Context, subscription int, database int) (*LatestImportStatus, error) {
	message := fmt.Sprintf("get latest import information for database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("/fixed/subscriptions/%d/databases/%d/import", subscription, database)
	task, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}
	return task, nil
}

func (a *API) get(ctx context.Context, message string, address string) (*LatestImportStatus, error) {
	var task internal.TaskResponse
	err := a.client.Get(ctx, message, address, &task)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Waiting for import status request %d to complete", task.ID)

	taskResp, err := a.taskWaiter.WaitForTask(ctx, *task.ID)
	if err != nil {
		var iErr *internal.Error
		if errors.As(err, &iErr) && taskResp != nil {
			importStatusTask, err := createLatestImportStatusFromTask(taskResp)
			if err != nil {
				return nil, err
			}
			return importStatusTask, nil
		}
		return nil, fmt.Errorf("failed to retrieve completed backup status %d: %w", task.ID, err)
	}

	a.logger.Printf("Import status request %d completed, possibly with error: %v", task.ID, err)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve completed import status %d: %w", task.ID, err)
	}

	return createLatestImportStatusFromTask(taskResp)
}

func wrap404Error(subId int, dbId int, err error) error {
	var httpErr *internal.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
		return &NotFound{subId: subId, dbId: dbId}
	}
	return err
}

func createLatestImportStatusFromTask(task *internal.Task) (*LatestImportStatus, error) {
	latestImportStatus := &LatestImportStatus{}
	if task != nil {
		latestImportStatus.CommandType = task.CommandType
		latestImportStatus.Description = task.Description
		latestImportStatus.Status = task.Status
		latestImportStatus.ID = task.ID
		if task.Response != nil {
			latestImportStatus.Response = &Response{
				ID: task.Response.ID,
			}
			if task.Response.Error != nil {
				latestImportStatus.Response.Error = &Error{
					Type:        task.Response.Error.Type,
					Description: task.Response.Error.Description,
					Status:      task.Response.Error.Status,
				}
			}

			if task.Response.Resource != nil {
				latestImportStatus.Response.Resource = &Resource{}
				err := json.Unmarshal(*task.Response.Resource, latestImportStatus.Response.Resource)
				if err != nil {
					return nil, fmt.Errorf("failed to unmarshal task response: %w", err)
				}
			}
		}
	}
	return latestImportStatus, nil
}
