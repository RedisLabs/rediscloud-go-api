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

type Api interface {
	// WaitForResourceId will poll the Task, waiting for the Task to finish processing, where it will then return.
	// An error will be returned if the Task couldn't be retrieved or the Task was not processed successfully.
	//
	// The Task will be continuously polled until the Task either fails or succeeds - cancellation can be achieved
	// by cancelling the context.
	WaitForResourceId(ctx context.Context, id string) (int, error)

	// Wait will poll the Task, waiting for the Task to finish processing, where it will then return.
	// An error will be returned if the Task couldn't be retrieved or the Task was not processed successfully.
	//
	// The Task will be continuously polled until the Task either fails or succeeds - cancellation can be achieved
	// by cancelling the context.
	Wait(ctx context.Context, id string) error

	// WaitForResource will poll the Task, waiting for the Task to finish processing, where it will then marshal the
	// returned resource into the value pointed to be `resource`.
	//
	// The Task will be continuously polled until the Task either fails or succeeds - cancellation can be achieved
	// by cancelling the context.
	WaitForResource(ctx context.Context, id string, resource interface{}) error

	// WaitForTask will poll the Task, waiting for it to enter a terminal state (i.e Done or Error). This Task
	// will then be returned, or an error in case it cannot be retrieved.
	WaitForTask(ctx context.Context, id string) (*Task, error)
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
	return err
}

func (a *api) WaitForResource(ctx context.Context, id string, resource interface{}) error {
	task, err := a.waitForTaskToComplete(ctx, id)
	if err != nil {
		return err
	}

	return json.Unmarshal(*task.Response.Resource, resource)
}

func (a *api) waitForTaskToComplete(ctx context.Context, id string) (*Task, error) {
	var task *Task
	notFoundCount := 0
	err := retry.Do(
		func() error {
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
		retry.Attempts(math.MaxUint16),
		retry.Delay(1*time.Second),
		retry.MaxDelay(30*time.Second),
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

func (a *api) WaitForTask(ctx context.Context, id string) (*Task, error) {
	var task *Task
	notFoundCount := 0
	err := retry.Do(
		func() error {
			var err error
			task, err = a.get(ctx, id)
			if err != nil {
				// An error is a terminal state (any repeated pre-task 404s will have been exhausted by this point)
				return nil
			}

			status := redis.StringValue(task.Status)

			if _, ok := processingStates[status]; !ok {
				// The task is no longer processing for whatever reason
				return nil
			}

			return fmt.Errorf("task %s not processed yet: %s", id, status)
		},
		retry.Attempts(math.MaxUint16),
		retry.Delay(1*time.Second),
		retry.MaxDelay(30*time.Second),
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

func (a *api) get(ctx context.Context, id string) (*Task, error) {
	var task Task
	if err := a.client.Get(ctx, fmt.Sprintf("retrieve Task %s", id), "/tasks/"+url.PathEscape(id), &task); err != nil {
		return nil, err
	}

	if task.Response != nil && task.Response.Error != nil {
		return &task, task.Response.Error
	}

	return &task, nil
}

// Number of 404 errors to swallow before returning an error while waiting for a Task to finish.
//
// There's a short window between the API returning a Task ID and the Task being known by the
// Task service, so by ignoring _a number_ of 404 errors we give the Task service enough time to
// learn about the Task but also handle the situation where there really is no Task.
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
