package databases

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
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

type API struct {
	client HttpClient
	task   Task
	logger Log
}

func NewAPI(client HttpClient, task Task, logger Log) *API {
	return &API{client: client, task: task, logger: logger}
}

func (a *API) List(ctx context.Context, subscription int) *ListDatabase {
	return newListDatabase(ctx, a.client, subscription, 100)
}

func (a *API) Get(ctx context.Context, subscription int, database int) (*Database, error) {
	var db Database
	err := a.client.Get(ctx, fmt.Sprintf("get database %d for subscription %d", subscription, database), fmt.Sprintf("/subscriptions/%d/databases/%d", subscription, database), &db)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

func (a *API) Delete(ctx context.Context, subscription int, database int) error {
	var task taskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("delete database %d/%d", subscription, database), fmt.Sprintf("/subscriptions/%d/databases/%d", subscription, database), &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for database %d for subscription %d to finish being deleted", subscription, database)

	err = a.task.Wait(ctx, redis.StringValue(task.ID))
	if err != nil {
		return err
	}

	return nil
}

type ListDatabase struct {
	client       HttpClient
	subscription int
	ctx          context.Context
	pageSize     int

	offset int
	value  []*Database
	err    error
}

func newListDatabase(ctx context.Context, client HttpClient, subscription int, pageSize int) *ListDatabase {
	return &ListDatabase{client: client, subscription: subscription, ctx: ctx, pageSize: pageSize}
}

func (d *ListDatabase) Next() bool {
	if d.err != nil {
		return false
	}

	url := fmt.Sprintf("/subscriptions/%d/databases?limit=%d&offset=%d", d.subscription, d.pageSize, d.offset)

	var list listDatabaseResponse
	err := d.client.Get(d.ctx, fmt.Sprintf("list databases for %d", d.subscription), url, &list)
	if err != nil {
		d.setError(err)
		return false
	}

	if len(list.Subscription) != 1 || redis.IntValue(list.Subscription[0].ID) != d.subscription {
		d.setError(fmt.Errorf("server didn't respond with just a single subscription"))
		return false
	}

	d.value = list.Subscription[0].Databases
	d.offset += d.pageSize

	return true
}

func (d *ListDatabase) Value() []*Database {
	return d.value
}

func (d *ListDatabase) Err() error {
	return d.err
}

func (d *ListDatabase) setError(err error) {
	if httpErr, ok := err.(*internal.HTTPError); !ok || httpErr.StatusCode != http.StatusNotFound {
		d.err = err
	}
	d.value = nil
}
