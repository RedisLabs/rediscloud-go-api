package agentmemory

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

const memoryStoresPath = "/memory-stores"

type Log interface {
	Printf(format string, args ...interface{})
}

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Patch(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Delete(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
}

type TaskWaiter interface {
	WaitForResource(ctx context.Context, id string, resource interface{}) error
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

func (a *API) Create(ctx context.Context, store CreateStore) (string, error) {
	var task internal.TaskResponse
	if err := a.client.Post(ctx, "create Agent Memory store", memoryStoresPath, store, &task); err != nil {
		return "", err
	}

	taskID, err := requireTaskID(task, "create Agent Memory store")
	if err != nil {
		return "", err
	}

	a.logger.Printf("Waiting for Agent Memory store to finish being created")

	var response CreateStoreResponse
	if err := a.taskWaiter.WaitForResource(ctx, taskID, &response); err != nil {
		return "", err
	}
	if response.StoreID == "" {
		return "", errors.New("create Agent Memory store task response did not contain a store ID")
	}

	return response.StoreID, nil
}

func (a *API) List(ctx context.Context) ([]StoreSummary, error) {
	var response ListStoresResponse
	if err := a.client.Get(ctx, "list Agent Memory stores", memoryStoresPath, &response); err != nil {
		return nil, err
	}

	return response.Data, nil
}

func (a *API) Get(ctx context.Context, storeID string) (*Store, error) {
	var store Store
	path := storePath(storeID)
	if err := a.client.Get(ctx, "get Agent Memory store", path, &store); err != nil {
		return nil, wrapNotFound(storeID, err)
	}

	return &store, nil
}

func (a *API) Update(ctx context.Context, storeID string, store UpdateStore) error {
	var task internal.TaskResponse
	path := storePath(storeID)
	if err := a.client.Patch(ctx, "update Agent Memory store", path, store, &task); err != nil {
		return wrapNotFound(storeID, err)
	}

	taskID, err := requireTaskID(task, "update Agent Memory store")
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for Agent Memory store %s to finish being updated", storeID)

	return a.taskWaiter.Wait(ctx, taskID)
}

func (a *API) Delete(ctx context.Context, storeID string) error {
	var task internal.TaskResponse
	path := storePath(storeID)
	if err := a.client.Delete(ctx, "delete Agent Memory store", path, nil, &task); err != nil {
		return wrapNotFound(storeID, err)
	}

	taskID, err := requireTaskID(task, "delete Agent Memory store")
	if err != nil {
		return err
	}

	a.logger.Printf("Waiting for Agent Memory store %s to finish being deleted", storeID)

	return a.taskWaiter.Wait(ctx, taskID)
}

func (a *API) ListAPIKeys(ctx context.Context, storeID string) ([]APIKeyInfo, error) {
	var response ListAPIKeysResponse
	path := storeAPIKeysPath(storeID)
	if err := a.client.Get(ctx, "list Agent Memory API keys", path, &response); err != nil {
		return nil, wrapNotFound(storeID, err)
	}

	return response.APIKeys, nil
}

// CreateAPIKey creates a named data-plane API key and resolves the ID that the
// create response omits by comparing the API key lists before and after POST.
func (a *API) CreateAPIKey(ctx context.Context, storeID, keyName string) (*CreatedAPIKey, error) {
	existing, err := a.ListAPIKeys(ctx, storeID)
	if err != nil {
		return nil, err
	}
	existingIDs := make(map[string]struct{}, len(existing))
	for _, key := range existing {
		existingIDs[key.APIKeyID] = struct{}{}
	}

	var response CreateAPIKeyResponse
	path := storeAPIKeysPath(storeID)
	if err := a.client.Post(ctx, "create Agent Memory API key", path, CreateAPIKeyRequest{KeyName: keyName}, &response); err != nil {
		return nil, wrapNotFound(storeID, err)
	}
	if response.APIKey == "" {
		return nil, errors.New("create Agent Memory API key response did not contain an API key")
	}

	keys, err := a.ListAPIKeys(ctx, storeID)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if _, existed := existingIDs[key.APIKeyID]; !existed && key.KeyName == keyName {
			return &CreatedAPIKey{APIKey: response.APIKey, Info: key}, nil
		}
	}

	return nil, fmt.Errorf("created Agent Memory API key %q for store %q but could not resolve its ID", keyName, storeID)
}

func (a *API) GetAPIKey(ctx context.Context, storeID, apiKeyID string) (*APIKeyInfo, error) {
	keys, err := a.ListAPIKeys(ctx, storeID)
	if err != nil {
		return nil, err
	}
	for i := range keys {
		if keys[i].APIKeyID == apiKeyID {
			return &keys[i], nil
		}
	}

	return nil, &APIKeyNotFound{StoreID: storeID, APIKeyID: apiKeyID}
}

func (a *API) DeleteAPIKey(ctx context.Context, storeID, apiKeyID string) error {
	path := storeAPIKeyPath(storeID, apiKeyID)
	if err := a.client.Delete(ctx, "delete Agent Memory API key", path, nil, nil); err != nil {
		var httpErr *internal.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
			return &APIKeyNotFound{StoreID: storeID, APIKeyID: apiKeyID}
		}
		return wrapNotFound(storeID, err)
	}

	return nil
}

func storePath(storeID string) string {
	return fmt.Sprintf("%s/%s", memoryStoresPath, url.PathEscape(storeID))
}

func storeAPIKeysPath(storeID string) string {
	return fmt.Sprintf("%s/api-keys", storePath(storeID))
}

func storeAPIKeyPath(storeID, apiKeyID string) string {
	return fmt.Sprintf("%s/%s", storeAPIKeysPath(storeID), url.PathEscape(apiKeyID))
}

func requireTaskID(task internal.TaskResponse, operation string) (string, error) {
	if task.ID == nil || *task.ID == "" {
		return "", fmt.Errorf("%s response did not contain a task ID", operation)
	}

	return *task.ID, nil
}

func wrapNotFound(storeID string, err error) error {
	var httpErr *internal.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
		return &NotFound{StoreID: storeID}
	}
	return err
}
