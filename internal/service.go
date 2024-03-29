package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/avast/retry-go/v4"
)

type Log interface {
	Println(v ...interface{})
}

type Api interface {
	// WaitForResourceId will poll the task, waiting for the task to finish processing, where it will then return.
	// An error will be returned if the task couldn't be retrieved or the task was not processed successfully.
	//
	// The task will be continuously polled until the task either fails or succeeds - cancellation can be achieved
	// by cancelling the context.
	WaitForResourceId(ctx context.Context, id string) (int, error)

	// Wait will poll the task, waiting for the task to finish processing, where it will then return.
	// An error will be returned if the task couldn't be retrieved or the task was not processed successfully.
	//
	// The task will be continuously polled until the task either fails or succeeds - cancellation can be achieved
	// by cancelling the context.
	Wait(ctx context.Context, id string) error

	// WaitForResource will poll the task, waiting for the task to finish processing, where it will then marshal the
	// returned resource into the value pointed to be `resource`.
	//
	// The task will be continuously polled until the task either fails or succeeds - cancellation can be achieved
	// by cancelling the context.
	WaitForResource(ctx context.Context, id string, resource interface{}) error
}

type api struct {
	client *HttpClient
	logger Log
}

func NewAPI(client *HttpClient, logger Log) Api {
	return &api{client: client, logger: logger}
}

func (a *api) WaitForResourceId(ctx context.Context, id string) (int, error) {
	task, err := a.waitForTaskToComplete(ctx, id)
	if err != nil {
		return 0, err
	}

	return redis.IntValue(task.Response.ID), nil
}

func (a *api) Wait(ctx context.Context, id string) error {
	_, err := a.waitForTaskToComplete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (a *api) WaitForResource(ctx context.Context, id string, resource interface{}) error {
	task, err := a.waitForTaskToComplete(ctx, id)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*task.Response.Resource, resource)
	if err != nil {
		return err
	}

	return nil
}

func (a *api) waitForTaskToComplete(ctx context.Context, id string) (*task, error) {
	var task *task
	notFoundCount := 0
	err := retry.Do(func() error {
		var err error
		task, err = a.get(ctx, id)
		if err != nil {
			if status, ok := err.(*HTTPError); ok && status.StatusCode == 404 {
				return &taskNotFoundError{err}
			}
			return retry.Unrecoverable(err)
		}

		status := redis.StringValue(task.Status)
		if status == processedState {
			return nil
		}

		if _, ok := processingStates[status]; !ok {
			return retry.Unrecoverable(fmt.Errorf("task %s failed %s - %s", id, status, redis.StringValue(task.Description)))
		}

		return fmt.Errorf("task %s not processed yet: %s", id, status)
	},
		retry.Attempts(math.MaxUint16), retry.Delay(1*time.Second), retry.MaxDelay(30*time.Second),
		retry.RetryIf(func(err error) bool {
			if !retry.IsRecoverable(err) {
				return false
			}
			if _, ok := err.(*taskNotFoundError); ok {
				notFoundCount++
				if notFoundCount > max404Errors {
					return false
				}
			}
			return true
		}),
		retry.LastErrorOnly(true), retry.Context(ctx), retry.OnRetry(func(_ uint, err error) {
			a.logger.Println(err)
		}))
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (a *api) get(ctx context.Context, id string) (*task, error) {
	var task task
	if err := a.client.Get(ctx, fmt.Sprintf("retrieve task %s", id), "/tasks/"+url.PathEscape(id), &task); err != nil {
		return nil, err
	}

	if task.Response != nil && task.Response.Error != nil {
		return nil, task.Response.Error
	}

	return &task, nil
}

// Number of 404 errors to swallow before returning an error while waiting for a task to finish.
//
// There's a short window between the API returning a task ID and the task being known by the
// Task service, so by ignoring _a number_ of 404 errors we give the task service enough time to
// learn about the task but also handle the situation where there really is no task.
const max404Errors = 5

var processingStates = map[string]bool{
	"initialized":            true,
	"received":               true,
	"processing-in-progress": true,
}

const processedState = "processing-completed"

type taskNotFoundError struct {
	wrapped error
}

func (e *taskNotFoundError) Error() string {
	return e.wrapped.Error()
}

func (e taskNotFoundError) Unwrap() error {
	return e.wrapped
}
