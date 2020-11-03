package databases

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

// TODO pagination
// TODO api returns 404 when no databases
func (a Api) List(ctx context.Context, subscription int) ([]*Database, error) {
	var list listDatabaseResponse
	err := a.client.Get(ctx, fmt.Sprintf("list databases for %d", subscription), fmt.Sprintf("/subscriptions/%d/databases", subscription), &list)
	if err != nil {
		return nil, err
	}

	if len(list.Subscription) != 1 || list.Subscription[0].ID != subscription {
		return nil, fmt.Errorf("server didn't respond with just a single subscription")
	}

	return list.Subscription[0].Databases, nil
}

func (a Api) Get(ctx context.Context, subscription int, database int) (*Database, error) {
	var db Database
	err := a.client.Get(ctx, fmt.Sprintf("get database %d for subscription %d", subscription, database), fmt.Sprintf("/subscriptions/%d/databases/%d", subscription, database), &db)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (a Api) Delete(ctx context.Context, subscription int, database int) error {
	var task taskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("delete database %d/%d", subscription, database), fmt.Sprintf("/subscriptions/%d/databases/%d", subscription, database), &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for database %d for subscription %d to finish being deleted", subscription, database)

	err = a.task.Wait(ctx, task.TaskId)
	if err != nil {
		return err
	}

	return nil
}
