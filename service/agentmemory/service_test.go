package agentmemory

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAPI_Create(t *testing.T) {
	client := &fakeHTTPClient{
		postFunc: func(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
			assert.Equal(t, "create Agent Memory store", name)
			assert.Equal(t, memoryStoresPath, path)
			assert.Equal(t, CreateStore{Name: "store", DatabaseID: 123}, requestBody)
			setTaskID(t, responseBody, "task-1")
			return nil
		},
	}
	taskWaiter := &fakeTaskWaiter{
		waitForResourceFunc: func(ctx context.Context, id string, resource interface{}) error {
			assert.Equal(t, "task-1", id)
			response, ok := resource.(*CreateStoreResponse)
			require.True(t, ok)
			response.StoreID = "store-1"
			return nil
		},
	}

	subject := NewAPI(client, taskWaiter, &fakeLogger{})

	actual, err := subject.Create(context.TODO(), CreateStore{Name: "store", DatabaseID: 123})

	require.NoError(t, err)
	assert.Equal(t, "store-1", actual)
}

func TestAPI_List(t *testing.T) {
	client := &fakeHTTPClient{
		getFunc: func(ctx context.Context, name, path string, responseBody interface{}) error {
			assert.Equal(t, "list Agent Memory stores", name)
			assert.Equal(t, memoryStoresPath, path)
			response, ok := responseBody.(*ListStoresResponse)
			require.True(t, ok)
			response.Data = []StoreSummary{{StoreID: "store-1", Name: "store"}}
			return nil
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	actual, err := subject.List(context.TODO())

	require.NoError(t, err)
	assert.Equal(t, []StoreSummary{{StoreID: "store-1", Name: "store"}}, actual)
}

func TestAPI_Get(t *testing.T) {
	client := &fakeHTTPClient{
		getFunc: func(ctx context.Context, name, path string, responseBody interface{}) error {
			assert.Equal(t, "get Agent Memory store", name)
			assert.Equal(t, "/memory-stores/store-1", path)
			response, ok := responseBody.(*Store)
			require.True(t, ok)
			response.StoreID = "store-1"
			response.Name = "store"
			return nil
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	actual, err := subject.Get(context.TODO(), "store-1")

	require.NoError(t, err)
	assert.Equal(t, &Store{StoreID: "store-1", Name: "store"}, actual)
}

func TestAPI_Update(t *testing.T) {
	update := UpdateStore{Name: stringPtr("updated")}
	client := &fakeHTTPClient{
		patchFunc: func(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
			assert.Equal(t, "update Agent Memory store", name)
			assert.Equal(t, "/memory-stores/store-1", path)
			assert.Equal(t, update, requestBody)
			setTaskID(t, responseBody, "task-2")
			return nil
		},
	}
	taskWaiter := &fakeTaskWaiter{
		waitFunc: func(ctx context.Context, id string) error {
			assert.Equal(t, "task-2", id)
			return nil
		},
	}

	subject := NewAPI(client, taskWaiter, &fakeLogger{})

	err := subject.Update(context.TODO(), "store-1", update)

	require.NoError(t, err)
}

func TestAPI_Delete(t *testing.T) {
	client := &fakeHTTPClient{
		deleteFunc: func(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
			assert.Equal(t, "delete Agent Memory store", name)
			assert.Equal(t, "/memory-stores/store-1", path)
			assert.Nil(t, requestBody)
			setTaskID(t, responseBody, "task-3")
			return nil
		},
	}
	taskWaiter := &fakeTaskWaiter{
		waitFunc: func(ctx context.Context, id string) error {
			assert.Equal(t, "task-3", id)
			return nil
		},
	}

	subject := NewAPI(client, taskWaiter, &fakeLogger{})

	err := subject.Delete(context.TODO(), "store-1")

	require.NoError(t, err)
}

func TestAPI_Get_wrapsNotFound(t *testing.T) {
	client := &fakeHTTPClient{
		getFunc: func(ctx context.Context, name, path string, responseBody interface{}) error {
			return &internal.HTTPError{StatusCode: http.StatusNotFound}
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	actual, err := subject.Get(context.TODO(), "missing-store")

	assert.Nil(t, actual)
	assert.IsType(t, &NotFound{}, err)
}

func TestAPI_Create_errorsWhenTaskIDIsMissing(t *testing.T) {
	client := &fakeHTTPClient{
		postFunc: func(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
			return nil
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	actual, err := subject.Create(context.TODO(), CreateStore{Name: "store", DatabaseID: 123})

	assert.Empty(t, actual)
	assert.EqualError(t, err, "create Agent Memory store response did not contain a task ID")
}

func TestAPI_CreateAPIKeyResolvesNewKeyID(t *testing.T) {
	listCalls := 0
	client := &fakeHTTPClient{
		getFunc: func(ctx context.Context, name, path string, responseBody interface{}) error {
			assert.Equal(t, "list Agent Memory API keys", name)
			assert.Equal(t, "/memory-stores/store-1/api-keys", path)
			listCalls++
			keys := []APIKeyInfo{{APIKeyID: "existing-id", KeyName: "existing"}}
			if listCalls == 2 {
				keys = append(keys, APIKeyInfo{
					APIKeyID:        "new-id",
					KeyName:         "terraform",
					ObfuscatedToken: "mem1...abcd",
					CreatedAt:       1234,
				})
			}
			response, ok := responseBody.(*ListAPIKeysResponse)
			require.True(t, ok)
			response.APIKeys = keys
			return nil
		},
		postFunc: func(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
			assert.Equal(t, "create Agent Memory API key", name)
			assert.Equal(t, "/memory-stores/store-1/api-keys", path)
			assert.Equal(t, CreateAPIKeyRequest{KeyName: "terraform"}, requestBody)
			response, ok := responseBody.(*CreateAPIKeyResponse)
			require.True(t, ok)
			response.APIKey = "secret-token"
			return nil
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	created, err := subject.CreateAPIKey(context.TODO(), "store-1", "terraform")

	require.NoError(t, err)
	assert.Equal(t, "secret-token", created.APIKey)
	assert.Equal(t, "new-id", created.Info.APIKeyID)
}

func TestAPI_CreateAPIKey_errorsWhenAPIKeyIsMissing(t *testing.T) {
	client := &fakeHTTPClient{
		getFunc: func(ctx context.Context, name, path string, responseBody interface{}) error {
			response, ok := responseBody.(*ListAPIKeysResponse)
			require.True(t, ok)
			response.APIKeys = []APIKeyInfo{}
			return nil
		},
		postFunc: func(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
			return nil
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	created, err := subject.CreateAPIKey(context.TODO(), "store-1", "terraform")

	assert.Nil(t, created)
	assert.EqualError(t, err, "create Agent Memory API key response did not contain an API key")
}

func TestAPI_GetAPIKey(t *testing.T) {
	client := &fakeHTTPClient{
		getFunc: func(ctx context.Context, name, path string, responseBody interface{}) error {
			assert.Equal(t, "list Agent Memory API keys", name)
			assert.Equal(t, "/memory-stores/store-1/api-keys", path)
			response, ok := responseBody.(*ListAPIKeysResponse)
			require.True(t, ok)
			response.APIKeys = []APIKeyInfo{
				{APIKeyID: "first-id", KeyName: "first"},
				{APIKeyID: "wanted-id", KeyName: "wanted"},
			}
			return nil
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	actual, err := subject.GetAPIKey(context.TODO(), "store-1", "wanted-id")

	require.NoError(t, err)
	assert.Equal(t, &APIKeyInfo{APIKeyID: "wanted-id", KeyName: "wanted"}, actual)
}

func TestAPI_GetAPIKey_returnsAPIKeyNotFound(t *testing.T) {
	client := &fakeHTTPClient{
		getFunc: func(ctx context.Context, name, path string, responseBody interface{}) error {
			response, ok := responseBody.(*ListAPIKeysResponse)
			require.True(t, ok)
			response.APIKeys = []APIKeyInfo{{APIKeyID: "other-id", KeyName: "other"}}
			return nil
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	actual, err := subject.GetAPIKey(context.TODO(), "store-1", "missing-id")

	assert.Nil(t, actual)
	assert.IsType(t, &APIKeyNotFound{}, err)
}

func TestAPI_DeleteAPIKey(t *testing.T) {
	client := &fakeHTTPClient{
		deleteFunc: func(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
			assert.Equal(t, "delete Agent Memory API key", name)
			assert.Equal(t, "/memory-stores/store-1/api-keys/key-id", path)
			assert.Nil(t, requestBody)
			assert.Nil(t, responseBody)
			return nil
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	err := subject.DeleteAPIKey(context.TODO(), "store-1", "key-id")

	require.NoError(t, err)
}

func TestAPI_DeleteAPIKey_wrapsNotFound(t *testing.T) {
	client := &fakeHTTPClient{
		deleteFunc: func(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
			return &internal.HTTPError{StatusCode: http.StatusNotFound}
		},
	}

	subject := NewAPI(client, &fakeTaskWaiter{}, &fakeLogger{})

	err := subject.DeleteAPIKey(context.TODO(), "store-1", "missing-id")

	assert.IsType(t, &APIKeyNotFound{}, err)
}

type fakeHTTPClient struct {
	getFunc    func(context.Context, string, string, interface{}) error
	postFunc   func(context.Context, string, string, interface{}, interface{}) error
	patchFunc  func(context.Context, string, string, interface{}, interface{}) error
	deleteFunc func(context.Context, string, string, interface{}, interface{}) error
}

func (f *fakeHTTPClient) Get(ctx context.Context, name, path string, responseBody interface{}) error {
	if f.getFunc == nil {
		return errors.New("unexpected Get call")
	}
	return f.getFunc(ctx, name, path, responseBody)
}

func (f *fakeHTTPClient) Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	if f.postFunc == nil {
		return errors.New("unexpected Post call")
	}
	return f.postFunc(ctx, name, path, requestBody, responseBody)
}

func (f *fakeHTTPClient) Patch(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	if f.patchFunc == nil {
		return errors.New("unexpected Patch call")
	}
	return f.patchFunc(ctx, name, path, requestBody, responseBody)
}

func (f *fakeHTTPClient) Delete(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error {
	if f.deleteFunc == nil {
		return errors.New("unexpected Delete call")
	}
	return f.deleteFunc(ctx, name, path, requestBody, responseBody)
}

type fakeTaskWaiter struct {
	waitForResourceFunc func(context.Context, string, interface{}) error
	waitFunc            func(context.Context, string) error
}

func (f *fakeTaskWaiter) WaitForResource(ctx context.Context, id string, resource interface{}) error {
	if f.waitForResourceFunc == nil {
		return errors.New("unexpected WaitForResource call")
	}
	return f.waitForResourceFunc(ctx, id, resource)
}

func (f *fakeTaskWaiter) Wait(ctx context.Context, id string) error {
	if f.waitFunc == nil {
		return errors.New("unexpected Wait call")
	}
	return f.waitFunc(ctx, id)
}

type fakeLogger struct{}

func (f *fakeLogger) Printf(format string, args ...interface{}) {}

func setTaskID(t *testing.T, responseBody interface{}, id string) {
	t.Helper()
	task, ok := responseBody.(*internal.TaskResponse)
	require.True(t, ok)
	task.ID = &id
}

func stringPtr(value string) *string {
	return &value
}
