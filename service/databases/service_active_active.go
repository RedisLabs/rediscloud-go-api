package databases

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
)

// ActiveActiveCreate will create a new database for the subscription and return the identifier of the database.
func (a *API) ActiveActiveCreate(ctx context.Context, subscription int, db CreateActiveActiveDatabase) (int, error) {
	var task internal.TaskResponse
	err := a.client.Post(ctx, fmt.Sprintf("create database for subscription %d", subscription), fmt.Sprintf("/subscriptions/%d/databases", subscription), db, &task)
	if err != nil {
		return 0, err
	}

	a.logger.Printf("Waiting for new database for subscription %d to finish being created", subscription)

	id, err := a.taskWaiter.WaitForResourceId(ctx, *task.ID)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// ActiveActiveUpdate will update certain values of an existing database.
func (a *API) ActiveActiveUpdate(ctx context.Context, subscription int, database int, update UpdateActiveActiveDatabase) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update database %d for subscription %d", database, subscription), fmt.Sprintf("/subscriptions/%d/databases/%d/regions", subscription, database), update, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for database %d for subscription %d to finish being updated", database, subscription)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

// ListActiveActive will return a ListDatabase that is capable of paging through all of the databases associated with a
// subscription.
func (a *API) ListActiveActive(ctx context.Context, subscription int) *ListActiveActiveDatabase {
	return newListActiveActiveDatabase(ctx, a.client, subscription, 100)
}

// GetActiveActive will retrieve an existing database.
func (a *API) GetActiveActive(ctx context.Context, subscription int, database int) (*ActiveActiveDatabase, error) {
	var db ActiveActiveDatabase
	err := a.client.Get(ctx, fmt.Sprintf("get database %d for subscription %d", subscription, database), fmt.Sprintf("/subscriptions/%d/databases/%d", subscription, database), &db)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}

	return &db, nil
}

type ListActiveActiveDatabase struct {
	client       HttpClient
	subscription int
	ctx          context.Context
	pageSize     int

	offset int
	page   []*ActiveActiveDatabase
	err    error
	fin    bool
	value  *ActiveActiveDatabase
}

func newListActiveActiveDatabase(ctx context.Context, client HttpClient, subscription int, pageSize int) *ListActiveActiveDatabase {
	return &ListActiveActiveDatabase{client: client, subscription: subscription, ctx: ctx, pageSize: pageSize}
}

// Next attempts to retrieve the next page of databases and will return false if no more databases were found.
// Any error that occurs within this function can be retrieved from the `Err()` function.
func (d *ListActiveActiveDatabase) Next() bool {
	if d.err != nil {
		return false
	}

	if d.fin {
		return false
	}

	if len(d.page) == 0 {
		if err := d.nextPage(); err != nil {
			d.setError(err)
			return false
		}
		// If the page is still empty after fetching, we're done
		if len(d.page) == 0 {
			return false
		}
	}

	d.updateValue()

	return true
}

// Value returns the current page of databases.
func (d *ListActiveActiveDatabase) Value() *ActiveActiveDatabase {
	return d.value
}

// Err returns any error that occurred while trying to retrieve the next page of databases.
func (d *ListActiveActiveDatabase) Err() error {
	return d.err
}

func (d *ListActiveActiveDatabase) nextPage() error {
	u := fmt.Sprintf("/subscriptions/%d/databases", d.subscription)
	q := map[string][]string{
		"limit":  {strconv.Itoa(d.pageSize)},
		"offset": {strconv.Itoa(d.offset)},
	}

	var list listActiveActiveDatabaseResponse
	err := d.client.GetWithQuery(d.ctx, fmt.Sprintf("list databases for %d", d.subscription), u, q, &list)
	if err != nil {
		return err
	}

	if len(list.Subscription) != 1 || redis.IntValue(list.Subscription[0].ID) != d.subscription {
		return fmt.Errorf("server didn't respond with just a single subscription")
	}

	d.page = list.Subscription[0].Databases
	d.offset += d.pageSize

	return nil
}

func (d *ListActiveActiveDatabase) updateValue() {
	d.value = d.page[0]
	d.page = d.page[1:]
}

func (d *ListActiveActiveDatabase) setError(err error) {
	if httpErr, ok := err.(*internal.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
		d.fin = true
	} else {
		d.err = err
	}

	d.page = nil
	d.value = nil
}
