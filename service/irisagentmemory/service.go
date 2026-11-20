package irisagentmemory

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/RedisLabs/rediscloud-go-api/internal/iris"
)

// HttpClient performs a synchronous Iris request. It is shared by both connectors.
type HttpClient interface {
	Call(context.Context, string, string, url.Values, http.Header, interface{}, interface{}) (http.Header, error)
}

type API struct{ client HttpClient }

func NewAPI(client HttpClient) *API { return &API{client: client} }

// ListBuiltInDetectors calls GET /iris/v1/agent-memory/detectors.
func (a *API) ListBuiltInDetectors(ctx context.Context) (*ListBuiltInDetectorsResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/iris/v1/agent-memory/detectors"
	result, _, err := iris.Result[ListBuiltInDetectorsResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// GetServiceLiveness calls GET /iris/v1/agent-memory/health/live.
func (a *API) GetServiceLiveness(ctx context.Context) (*GetServiceLivenessResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/iris/v1/agent-memory/health/live"
	result, _, err := iris.Result[GetServiceLivenessResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// GetServiceReadiness calls GET /iris/v1/agent-memory/health/ready.
func (a *API) GetServiceReadiness(ctx context.Context) (*GetServiceReadinessResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/iris/v1/agent-memory/health/ready"
	result, _, err := iris.Result[GetServiceReadinessResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// GetMemoryTypeReferenceData calls GET /iris/v1/agent-memory/memory-types/reference.
func (a *API) GetMemoryTypeReferenceData(ctx context.Context) (*GetMemoryTypeReferenceDataResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/iris/v1/agent-memory/memory-types/reference"
	result, _, err := iris.Result[GetMemoryTypeReferenceDataResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// ListStores calls GET /iris/v1/agent-memory/stores.
func (a *API) ListStores(ctx context.Context, params *ListStoresParams) (*ListStoresResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
		if params.Limit != nil {
			query.Set("limit", strconv.FormatInt(int64(*params.Limit), 10))
		}
		if params.Cursor != nil {
			query.Set("cursor", *params.Cursor)
		}
	}
	path := "/iris/v1/agent-memory/stores"
	result, _, err := iris.Result[ListStoresResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// CreateStore calls POST /iris/v1/agent-memory/stores.
func (a *API) CreateStore(ctx context.Context, params *CreateStoreParams, body CreateStoreJSONRequestBody) (*CreateStoreResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores"
	result, _, err := iris.Result[CreateStoreResponseContent](ctx, a.client, http.MethodPost, path, query, headers, body)
	return result, err
}

// DeleteStore calls DELETE /iris/v1/agent-memory/stores/{storeID}.
func (a *API) DeleteStore(ctx context.Context, storeID string, params *DeleteStoreParams) error {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
		if params.Flush != nil {
			query.Set("flush", strconv.FormatBool(*params.Flush))
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID)
	_, err := a.client.Call(ctx, http.MethodDelete, path, query, headers, nil, nil)
	return err
}

// GetStore calls GET /iris/v1/agent-memory/stores/{storeID}.
func (a *API) GetStore(ctx context.Context, storeID string, params *GetStoreParams) (*GetStoreResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID)
	result, _, err := iris.Result[GetStoreResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// UpdateStore calls PATCH /iris/v1/agent-memory/stores/{storeID}.
func (a *API) UpdateStore(ctx context.Context, storeID string, params *UpdateStoreParams, body UpdateStoreJSONRequestBody) (*UpdateStoreResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID)
	result, _, err := iris.Result[UpdateStoreResponseContent](ctx, a.client, http.MethodPatch, path, query, headers, body)
	return result, err
}

// ListExtractors calls GET /iris/v1/agent-memory/stores/{storeID}/extractors.
func (a *API) ListExtractors(ctx context.Context, storeID string, params *ListExtractorsParams) (*ListExtractorsResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
		if params.Limit != nil {
			query.Set("limit", strconv.FormatInt(int64(*params.Limit), 10))
		}
		if params.Cursor != nil {
			query.Set("cursor", *params.Cursor)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/extractors"
	result, _, err := iris.Result[ListExtractorsResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// CreateExtractor calls POST /iris/v1/agent-memory/stores/{storeID}/extractors.
func (a *API) CreateExtractor(ctx context.Context, storeID string, params *CreateExtractorParams, body CreateExtractorJSONRequestBody) (*CreateExtractorResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/extractors"
	result, _, err := iris.Result[CreateExtractorResponseContent](ctx, a.client, http.MethodPost, path, query, headers, body)
	return result, err
}

// RemoveExtractor calls DELETE /iris/v1/agent-memory/stores/{storeID}/extractors/{extractorID}.
func (a *API) RemoveExtractor(ctx context.Context, storeID string, extractorID string, params *RemoveExtractorParams) error {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
		query.Set("expectedRevision", params.ExpectedRevision)
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/extractors/" + url.PathEscape(extractorID)
	_, err := a.client.Call(ctx, http.MethodDelete, path, query, headers, nil, nil)
	return err
}

// GetExtractor calls GET /iris/v1/agent-memory/stores/{storeID}/extractors/{extractorID}.
func (a *API) GetExtractor(ctx context.Context, storeID string, extractorID string, params *GetExtractorParams) (*GetExtractorResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/extractors/" + url.PathEscape(extractorID)
	result, _, err := iris.Result[GetExtractorResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// UpdateExtractor calls PATCH /iris/v1/agent-memory/stores/{storeID}/extractors/{extractorID}.
func (a *API) UpdateExtractor(ctx context.Context, storeID string, extractorID string, params *UpdateExtractorParams, body UpdateExtractorJSONRequestBody) (*UpdateExtractorResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/extractors/" + url.PathEscape(extractorID)
	result, _, err := iris.Result[UpdateExtractorResponseContent](ctx, a.client, http.MethodPatch, path, query, headers, body)
	return result, err
}

// ListMemoryTypes calls GET /iris/v1/agent-memory/stores/{storeID}/memory-types.
func (a *API) ListMemoryTypes(ctx context.Context, storeID string, params *ListMemoryTypesParams) (*ListMemoryTypesResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
		if params.Limit != nil {
			query.Set("limit", strconv.FormatInt(int64(*params.Limit), 10))
		}
		if params.Cursor != nil {
			query.Set("cursor", *params.Cursor)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/memory-types"
	result, _, err := iris.Result[ListMemoryTypesResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// CreateMemoryType calls POST /iris/v1/agent-memory/stores/{storeID}/memory-types.
func (a *API) CreateMemoryType(ctx context.Context, storeID string, params *CreateMemoryTypeParams, body CreateMemoryTypeJSONRequestBody) (*CreateMemoryTypeResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/memory-types"
	result, _, err := iris.Result[CreateMemoryTypeResponseContent](ctx, a.client, http.MethodPost, path, query, headers, body)
	return result, err
}

// RemoveMemoryType calls DELETE /iris/v1/agent-memory/stores/{storeID}/memory-types/{memoryTypeName}.
func (a *API) RemoveMemoryType(ctx context.Context, storeID string, memoryTypeName string, params *RemoveMemoryTypeParams) error {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
		if params.ExpectedRevision != nil {
			query.Set("expectedRevision", *params.ExpectedRevision)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/memory-types/" + url.PathEscape(memoryTypeName)
	_, err := a.client.Call(ctx, http.MethodDelete, path, query, headers, nil, nil)
	return err
}

// GetMemoryType calls GET /iris/v1/agent-memory/stores/{storeID}/memory-types/{memoryTypeName}.
func (a *API) GetMemoryType(ctx context.Context, storeID string, memoryTypeName string, params *GetMemoryTypeParams) (*GetMemoryTypeResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/memory-types/" + url.PathEscape(memoryTypeName)
	result, _, err := iris.Result[GetMemoryTypeResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// UpdateMemoryType calls PATCH /iris/v1/agent-memory/stores/{storeID}/memory-types/{memoryTypeName}.
func (a *API) UpdateMemoryType(ctx context.Context, storeID string, memoryTypeName string, params *UpdateMemoryTypeParams, body UpdateMemoryTypeJSONRequestBody) (*UpdateMemoryTypeResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/memory-types/" + url.PathEscape(memoryTypeName)
	result, _, err := iris.Result[UpdateMemoryTypeResponseContent](ctx, a.client, http.MethodPatch, path, query, headers, body)
	return result, err
}

// GetIngestionHooks calls GET /iris/v1/agent-memory/stores/{storeID}/ingestion-hooks.
func (a *API) GetIngestionHooks(ctx context.Context, storeID string, params *GetIngestionHooksParams) (*GetIngestionHooksResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/ingestion-hooks"
	result, _, err := iris.Result[GetIngestionHooksResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// PutIngestionHooks calls PUT /iris/v1/agent-memory/stores/{storeID}/ingestion-hooks.
func (a *API) PutIngestionHooks(ctx context.Context, storeID string, params *PutIngestionHooksParams, body PutIngestionHooksJSONRequestBody) (*PutIngestionHooksResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ProjectId != nil {
			query.Set("projectId", *params.ProjectId)
		}
	}
	path := "/iris/v1/agent-memory/stores/" + url.PathEscape(storeID) + "/ingestion-hooks"
	result, _, err := iris.Result[PutIngestionHooksResponseContent](ctx, a.client, http.MethodPut, path, query, headers, body)
	return result, err
}
