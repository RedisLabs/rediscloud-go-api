package databases

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
)

type Log interface {
	Printf(format string, args ...interface{})
}

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	GetWithQuery(ctx context.Context, name, path string, query url.Values, responseBody interface{}) error
	Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Delete(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
}

type TaskWaiter interface {
	WaitForResourceId(ctx context.Context, id string) (int, error)
	Wait(ctx context.Context, id string) error
}

type API struct {
	client     HttpClient
	taskWaiter TaskWaiter
	logger     Log
}

func NewAPI(client HttpClient, taskWaiter TaskWaiter, logger Log) *API {
	return &API{client: client, taskWaiter: taskWaiter, logger: logger}
}

// Create will create a new database for the subscription and return the identifier of the database.
func (a *API) Create(ctx context.Context, subscription int, db CreateDatabase) (int, error) {
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

// List will return a ListDatabase that is capable of paging through all the databases associated with a
// subscription.
func (a *API) List(ctx context.Context, subscription int) *ListDatabase {
	return newListDatabase(ctx, a.client, subscription, 100)
}

// Get will retrieve an existing database.
func (a *API) Get(ctx context.Context, subscription int, database int) (*Database, error) {
	var db Database
	err := a.client.Get(ctx, fmt.Sprintf("get database %d for subscription %d", subscription, database), fmt.Sprintf("/subscriptions/%d/databases/%d", subscription, database), &db)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}

	return &db, nil
}

// Update will update certain values of an existing database.
func (a *API) Update(ctx context.Context, subscription int, database int, update UpdateDatabase) error {
	var task internal.TaskResponse
	err := a.client.Put(ctx, fmt.Sprintf("update database %d for subscription %d", database, subscription), fmt.Sprintf("/subscriptions/%d/databases/%d", subscription, database), update, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for database %d for subscription %d to finish being updated", database, subscription)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

// UpgradeRedisVersion will upgrade the Redis version of an existing database.
func (a *API) UpgradeRedisVersion(ctx context.Context, subscription int, database int, upgradeVersion UpgradeRedisVersion) error {
	var task internal.TaskResponse
	err := a.client.Post(ctx, fmt.Sprintf("upgrade database %d version for subscription %d", database, subscription), fmt.Sprintf("/subscriptions/%d/databases/%d/upgrade", subscription, database), upgradeVersion, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for database %d for subscription %d to finish being upgraded", database, subscription)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

// Delete will destroy an existing database.
func (a *API) Delete(ctx context.Context, subscription int, database int) error {
	var task internal.TaskResponse
	err := a.client.Delete(ctx, fmt.Sprintf("delete database %d/%d", subscription, database), fmt.Sprintf("/subscriptions/%d/databases/%d", subscription, database), nil, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for database %d for subscription %d to finish being deleted", database, subscription)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

// Backup will create a manual backup of the database to the destination the database has been configured to backup to.
func (a *API) Backup(ctx context.Context, subscription int, database int) error {
	var task internal.TaskResponse
	err := a.client.Post(ctx, fmt.Sprintf("backup database %d for subscription %d", database, subscription), fmt.Sprintf("/subscriptions/%d/databases/%d/backup", subscription, database), nil, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for backup of database %d for subscription %d to finish", database, subscription)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

// Import will import data from an RDB file or another Redis database into an existing database.
func (a *API) Import(ctx context.Context, subscription int, database int, request Import) error {
	var task internal.TaskResponse
	err := a.client.Post(ctx, fmt.Sprintf("import database %d for subscription %d", database, subscription), fmt.Sprintf("/subscriptions/%d/databases/%d/import", subscription, database), request, &task)
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for import into database %d for subscription %d to finish", database, subscription)

	return a.taskWaiter.Wait(ctx, *task.ID)
}

// GetCertificate retrieves the TLS certificate for the specified database within a subscription.
func (a *API) GetCertificate(ctx context.Context, subscription int, database int) (*DatabaseCertificate, error) {
	var certificate DatabaseCertificate
	getCertificateUrl := "/subscriptions/%d/databases/%d/certificate"

	path := fmt.Sprintf(getCertificateUrl, subscription, database)
	err := a.client.Get(
		ctx,
		fmt.Sprintf("get TLS certificate for database %d in subscription %d", database, subscription),
		path,
		&certificate,
	)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}
	return &certificate, nil
}

type ListDatabase struct {
	client       HttpClient
	subscription int
	ctx          context.Context
	pageSize     int

	offset int
	page   []*Database
	err    error
	fin    bool
	value  *Database
}

func newListDatabase(ctx context.Context, client HttpClient, subscription int, pageSize int) *ListDatabase {
	return &ListDatabase{client: client, subscription: subscription, ctx: ctx, pageSize: pageSize}
}

// Next attempts to retrieve the next page of databases and will return false if no more databases were found.
// Any error that occurs within this function can be retrieved from the `Err()` function.
func (d *ListDatabase) Next() bool {
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
	}

	d.updateValue()

	return true
}

// Value returns the current page of databases.
func (d *ListDatabase) Value() *Database {
	return d.value
}

// Err returns any error that occurred while trying to retrieve the next page of databases.
func (d *ListDatabase) Err() error {
	return d.err
}

func (d *ListDatabase) nextPage() error {
	u := fmt.Sprintf("/subscriptions/%d/databases", d.subscription)
	q := map[string][]string{
		"limit":  {strconv.Itoa(d.pageSize)},
		"offset": {strconv.Itoa(d.offset)},
	}

	var list listDatabaseResponse
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

func (d *ListDatabase) updateValue() {
	d.value = d.page[0]
	d.page = d.page[1:]
}

func (d *ListDatabase) setError(err error) {
	if httpErr, ok := err.(*internal.HTTPError); ok && httpErr.StatusCode == http.StatusNotFound {
		d.fin = true
	} else {
		d.err = err
	}

	d.page = nil
	d.value = nil
}

func wrap404Error(subId int, dbId int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{subId: subId, dbId: dbId}
	}
	return err
}
