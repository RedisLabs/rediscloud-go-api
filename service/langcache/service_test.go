package langcache

import (
	"context"
	"errors"
	"net/url"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubHTTPClient struct {
	getFn             func(context.Context, string, string, interface{}) error
	postFn            func(context.Context, string, string, interface{}, interface{}) error
	patchFn           func(context.Context, string, string, interface{}, interface{}) error
	deleteWithQueryFn func(context.Context, string, string, url.Values, interface{}, interface{}) error
}

func (s *stubHTTPClient) Get(ctx context.Context, name, path string, response interface{}) error {
	return s.getFn(ctx, name, path, response)
}

func (s *stubHTTPClient) Post(ctx context.Context, name, path string, request, response interface{}) error {
	return s.postFn(ctx, name, path, request, response)
}

func (s *stubHTTPClient) Patch(ctx context.Context, name, path string, request, response interface{}) error {
	return s.patchFn(ctx, name, path, request, response)
}

func (s *stubHTTPClient) DeleteWithQuery(ctx context.Context, name, path string, query url.Values, request, response interface{}) error {
	return s.deleteWithQueryFn(ctx, name, path, query, request, response)
}

func TestCreate(t *testing.T) {
	client := &stubHTTPClient{
		postFn: func(_ context.Context, name, path string, request, response interface{}) error {
			assert.Equal(t, "create LangCache", name)
			assert.Equal(t, cachesPath, path)
			assert.Equal(t, "cache-name", request.(CreateCache).Name)
			response.(*CreateCacheResponse).CacheID = "cache-id"
			return nil
		},
	}

	id, err := NewAPI(client).Create(context.Background(), CreateCache{Name: "cache-name"})
	require.NoError(t, err)
	assert.Equal(t, "cache-id", id)
}

func TestGetWrapsNotFound(t *testing.T) {
	client := &stubHTTPClient{
		getFn: func(_ context.Context, _, path string, _ interface{}) error {
			assert.Equal(t, cachesPath+"/missing", path)
			return &internal.HTTPError{StatusCode: 404}
		},
	}

	_, err := NewAPI(client).Get(context.Background(), "missing")
	var notFound *NotFound
	require.True(t, errors.As(err, &notFound))
	assert.Equal(t, "missing", notFound.CacheID)
}

func TestDelete(t *testing.T) {
	client := &stubHTTPClient{
		deleteWithQueryFn: func(_ context.Context, name, path string, query url.Values, _, _ interface{}) error {
			assert.Equal(t, "delete LangCache", name)
			assert.Equal(t, cachesPath+"/cache-id", path)
			assert.Equal(t, "true", query.Get("flush"))
			return nil
		},
	}

	require.NoError(t, NewAPI(client).Delete(context.Background(), "cache-id", true))
}

func TestUpdate(t *testing.T) {
	newName := "renamed-cache"
	threshold := 0.2
	client := &stubHTTPClient{
		patchFn: func(_ context.Context, name, path string, request, response interface{}) error {
			assert.Equal(t, "update LangCache", name)
			assert.Equal(t, cachesPath+"/cache-id", path)
			assert.Equal(t, UpdateCache{Name: &newName, DefaultSearchThreshold: &threshold}, request)
			assert.Nil(t, response)
			return nil
		},
	}

	require.NoError(t, NewAPI(client).Update(context.Background(), "cache-id", UpdateCache{Name: &newName, DefaultSearchThreshold: &threshold}))
}

func TestCreateAPIKeyResolvesNewKeyID(t *testing.T) {
	listCalls := 0
	client := &stubHTTPClient{
		getFn: func(_ context.Context, name, path string, response interface{}) error {
			assert.Equal(t, "list LangCache API keys", name)
			assert.Equal(t, cachesPath+"/cache-id/api-key", path)
			listCalls++
			keys := []APIKeyInfo{{APIKeyID: "existing-id", KeyName: "existing"}}
			if listCalls == 2 {
				keys = append(keys, APIKeyInfo{
					APIKeyID:        "new-id",
					KeyName:         "terraform",
					ObfuscatedToken: "****abcd",
					CreatedAt:       1234,
				})
			}
			response.(*ListAPIKeysResponse).APIKeys = keys
			return nil
		},
		postFn: func(_ context.Context, name, path string, request, response interface{}) error {
			assert.Equal(t, "create LangCache API key", name)
			assert.Equal(t, cachesPath+"/cache-id/api-key", path)
			assert.Equal(t, AddAPIKeyRequest{KeyName: "terraform"}, request)
			response.(*AddAPIKeyResponse).APIKey = "secret-token"
			return nil
		},
	}

	created, err := NewAPI(client).CreateAPIKey(context.Background(), "cache-id", "terraform")
	require.NoError(t, err)
	assert.Equal(t, "secret-token", created.APIKey)
	assert.Equal(t, "new-id", created.Info.APIKeyID)
}

func TestDeleteAPIKey(t *testing.T) {
	client := &stubHTTPClient{
		deleteWithQueryFn: func(_ context.Context, name, path string, query url.Values, _, _ interface{}) error {
			assert.Equal(t, "delete LangCache API key", name)
			assert.Equal(t, cachesPath+"/cache-id/api-key/key-id", path)
			assert.Nil(t, query)
			return nil
		},
	}

	require.NoError(t, NewAPI(client).DeleteAPIKey(context.Background(), "cache-id", "key-id"))
}
