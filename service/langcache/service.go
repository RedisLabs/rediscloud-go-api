package langcache

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

const cachesPath = "/iris/ui/caches"

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	Post(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	Patch(ctx context.Context, name, path string, requestBody interface{}, responseBody interface{}) error
	DeleteWithQuery(ctx context.Context, name, path string, query url.Values, requestBody interface{}, responseBody interface{}) error
}

type API struct {
	client HttpClient
}

func NewAPI(client HttpClient) *API {
	return &API{client: client}
}

func (a *API) Create(ctx context.Context, cache CreateCache) (string, error) {
	var response CreateCacheResponse
	if err := a.client.Post(ctx, "create LangCache", cachesPath, cache, &response); err != nil {
		return "", err
	}
	if response.CacheID == "" {
		return "", errors.New("create LangCache response did not contain a cache ID")
	}
	return response.CacheID, nil
}

func (a *API) Get(ctx context.Context, cacheID string) (*Cache, error) {
	var response Cache
	path := fmt.Sprintf("%s/%s", cachesPath, url.PathEscape(cacheID))
	if err := a.client.Get(ctx, "get LangCache", path, &response); err != nil {
		return nil, wrapNotFound(cacheID, err)
	}
	return &response, nil
}

func (a *API) Update(ctx context.Context, cacheID string, cache UpdateCache) error {
	path := fmt.Sprintf("%s/%s", cachesPath, url.PathEscape(cacheID))
	if err := a.client.Patch(ctx, "update LangCache", path, cache, nil); err != nil {
		return wrapNotFound(cacheID, err)
	}
	return nil
}

func (a *API) Delete(ctx context.Context, cacheID string, flush bool) error {
	path := fmt.Sprintf("%s/%s", cachesPath, url.PathEscape(cacheID))
	query := url.Values{"flush": []string{strconv.FormatBool(flush)}}
	if err := a.client.DeleteWithQuery(ctx, "delete LangCache", path, query, nil, nil); err != nil {
		return wrapNotFound(cacheID, err)
	}
	return nil
}

func (a *API) ListAPIKeys(ctx context.Context, cacheID string) ([]APIKeyInfo, error) {
	var response ListAPIKeysResponse
	path := fmt.Sprintf("%s/%s/api-key", cachesPath, url.PathEscape(cacheID))
	if err := a.client.Get(ctx, "list LangCache API keys", path, &response); err != nil {
		return nil, wrapNotFound(cacheID, err)
	}
	return response.APIKeys, nil
}

// CreateAPIKey creates a named data-plane API key and resolves the ID that the
// create response omits by comparing the API key lists before and after POST.
func (a *API) CreateAPIKey(ctx context.Context, cacheID, keyName string) (*CreatedAPIKey, error) {
	existing, err := a.ListAPIKeys(ctx, cacheID)
	if err != nil {
		return nil, err
	}
	existingIDs := make(map[string]struct{}, len(existing))
	for _, key := range existing {
		existingIDs[key.APIKeyID] = struct{}{}
	}

	var response AddAPIKeyResponse
	path := fmt.Sprintf("%s/%s/api-key", cachesPath, url.PathEscape(cacheID))
	if err := a.client.Post(ctx, "create LangCache API key", path, AddAPIKeyRequest{KeyName: keyName}, &response); err != nil {
		return nil, wrapNotFound(cacheID, err)
	}
	if response.APIKey == "" {
		return nil, errors.New("create LangCache API key response did not contain an API key")
	}

	keys, err := a.ListAPIKeys(ctx, cacheID)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if _, existed := existingIDs[key.APIKeyID]; !existed && key.KeyName == keyName {
			return &CreatedAPIKey{APIKey: response.APIKey, Info: key}, nil
		}
	}
	return nil, fmt.Errorf("created LangCache API key %q for cache %q but could not resolve its ID", keyName, cacheID)
}

func (a *API) GetAPIKey(ctx context.Context, cacheID, apiKeyID string) (*APIKeyInfo, error) {
	keys, err := a.ListAPIKeys(ctx, cacheID)
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if key.APIKeyID == apiKeyID {
			return &key, nil
		}
	}
	return nil, &APIKeyNotFound{CacheID: cacheID, APIKeyID: apiKeyID}
}

func (a *API) DeleteAPIKey(ctx context.Context, cacheID, apiKeyID string) error {
	path := fmt.Sprintf("%s/%s/api-key/%s", cachesPath, url.PathEscape(cacheID), url.PathEscape(apiKeyID))
	if err := a.client.DeleteWithQuery(ctx, "delete LangCache API key", path, nil, nil, nil); err != nil {
		var httpErr *internal.HTTPError
		if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
			return &APIKeyNotFound{CacheID: cacheID, APIKeyID: apiKeyID}
		}
		return err
	}
	return nil
}

func wrapNotFound(cacheID string, err error) error {
	var httpErr *internal.HTTPError
	if errors.As(err, &httpErr) && httpErr.StatusCode == http.StatusNotFound {
		return &NotFound{CacheID: cacheID}
	}
	return err
}
