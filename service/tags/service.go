package tags

import (
	"context"
	"fmt"
	"net/http"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	Put(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
}

type API struct {
	client HttpClient
}

func NewAPI(client HttpClient) *API {
	return &API{client: client}
}

func (a *API) Get(ctx context.Context, subscription int, database int) (*AllTags, error) {
	message := fmt.Sprintf("get tags for database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("/subscriptions/%d/databases/%d/tags", subscription, database)
	tags, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}
	return tags, nil
}

func (a *API) GetFixed(ctx context.Context, subscription int, database int) (*AllTags, error) {
	message := fmt.Sprintf("get tags for fixed database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("fixed/subscriptions/%d/databases/%d/tags", subscription, database)
	tags, err := a.get(ctx, message, address)
	if err != nil {
		return nil, wrap404Error(subscription, database, err)
	}
	return tags, nil
}

func (a *API) Put(ctx context.Context, subscription int, database int, tags AllTags) error {
	message := fmt.Sprintf("update tags for database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("/subscriptions/%d/databases/%d/tags", subscription, database)
	err := a.put(ctx, message, address, tags)
	if err != nil {
		return wrap404Error(subscription, database, err)
	}
	return nil
}

func (a *API) PutFixed(ctx context.Context, subscription int, database int, tags AllTags) error {
	message := fmt.Sprintf("update tags for fixed database %d in subscription %d", subscription, database)
	address := fmt.Sprintf("fixed/subscriptions/%d/databases/%d/tags", subscription, database)
	err := a.put(ctx, message, address, tags)
	if err != nil {
		return wrap404Error(subscription, database, err)
	}
	return nil
}

func (a *API) get(ctx context.Context, message string, address string) (*AllTags, error) {
	var tags AllTags
	err := a.client.Get(ctx, message, address, &tags)
	if err != nil {
		return nil, err
	}
	return &tags, nil
}

func (a *API) put(ctx context.Context, message string, address string, tags AllTags) error {
	var tagsResponse AllTags
	return a.client.Put(ctx, message, address, tags, &tagsResponse)
}

func wrap404Error(subId int, dbId int, err error) error {
	if v, ok := err.(*internal.HTTPError); ok && v.StatusCode == http.StatusNotFound {
		return &NotFound{subId: subId, dbId: dbId}
	}
	return err
}
