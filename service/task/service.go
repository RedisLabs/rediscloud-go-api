package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Api struct {
	client  *http.Client
	baseUrl string
}

func NewApi(client *http.Client, baseUrl string) *Api {
	return &Api{client: client, baseUrl: baseUrl}
}

// WaitForTaskToComplete will poll the task, waiting for the task to finish processing, where it will then return.
// An error will be returned if the task couldn't be retrieved or the task was not processed successfully.
//
// The task will be continuously polled until the task either fails or succeeds - cancellation can be achieved
// by cancelling the context.
func (a *Api) WaitForTaskToComplete(ctx context.Context, id string) (*Task, error) {
	for true {
		task, err := a.Get(ctx, id)
		if err != nil {
			return nil, err
		}

		if task.Status == processedState {
			return task, nil
		}

		if _, ok := processingStates[task.Status]; !ok {
			return nil, fmt.Errorf("task %s failed %s - %s", id, task.Status, task.Description)
		}

		time.Sleep(1 * time.Second)
	}

	return nil, fmt.Errorf("an infinite loop managed to finish")
}

// Get will retrieve a task. An error will be returned if the task couldn't be retrieved or the task itself
// failed.
func (a *Api) Get(ctx context.Context, id string) (*Task, error) {
	parsed, err := url.Parse(a.baseUrl)
	if err != nil {
		return nil, err
	}

	parsed.Path += "/tasks/" + url.PathEscape(id)

	u := parsed.String()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for task %s: %w", id, err)
	}

	response, err := a.client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve task %s: %w", id, err)
	}

	if response.StatusCode > 299 {
		return nil, fmt.Errorf("failed to get task %s: %d", id, response.StatusCode)
	}

	defer response.Body.Close()

	var task Task
	if err := json.NewDecoder(response.Body).Decode(&task); err != nil {
		return nil, fmt.Errorf("failed to decode response for task %s: %w", id, err)
	}

	if task.Response.Error != nil {
		return nil, task.Response.Error
	}

	return &task, nil
}

var processingStates = map[string]bool{
	"initialized":            true,
	"received":               true,
	"processing-in-progress": true,
}

const processedState = "processing-completed"
