package irisidentity

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

// ListAPIKeys calls GET /v1/api-keys.
func (a *API) ListAPIKeys(ctx context.Context, params *ListAPIKeysParams) (*ListAPIKeysResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.NextToken != nil {
			query.Set("nextToken", *params.NextToken)
		}
		if params.MaxResults != nil {
			query.Set("maxResults", strconv.FormatInt(int64(*params.MaxResults), 10))
		}
		if params.Product != nil {
			query.Set("product", string(*params.Product))
		}
		if params.ResourceType != nil {
			query.Set("resourceType", *params.ResourceType)
		}
		if params.ResourceId != nil {
			query.Set("resourceId", *params.ResourceId)
		}
	}
	path := "/v1/api-keys"
	result, _, err := iris.Result[ListAPIKeysResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// CreateAPIKey calls POST /v1/api-keys.
func (a *API) CreateAPIKey(ctx context.Context, body CreateAPIKeyJSONRequestBody) (*CreateAPIKeyResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/api-keys"
	result, _, err := iris.Result[CreateAPIKeyResponseContent](ctx, a.client, http.MethodPost, path, query, headers, body)
	return result, err
}

// RevokeAPIKey calls DELETE /v1/api-keys/{keyID}.
func (a *API) RevokeAPIKey(ctx context.Context, keyID string) error {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/api-keys/" + url.PathEscape(keyID)
	_, err := a.client.Call(ctx, http.MethodDelete, path, query, headers, nil, nil)
	return err
}

// GetAPIKey calls GET /v1/api-keys/{keyID}.
func (a *API) GetAPIKey(ctx context.Context, keyID string) (*GetAPIKeyResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/api-keys/" + url.PathEscape(keyID)
	result, _, err := iris.Result[GetAPIKeyResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// UpdateAPIKey calls PATCH /v1/api-keys/{keyID}.
func (a *API) UpdateAPIKey(ctx context.Context, keyID string, body UpdateAPIKeyJSONRequestBody) (*UpdateAPIKeyResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/api-keys/" + url.PathEscape(keyID)
	result, _, err := iris.Result[UpdateAPIKeyResponseContent](ctx, a.client, http.MethodPatch, path, query, headers, body)
	return result, err
}

// RotateAPIKey calls POST /v1/api-keys/{keyID}/rotate.
func (a *API) RotateAPIKey(ctx context.Context, keyID string, params *RotateAPIKeyParams) (*RotateAPIKeyResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.GraceSeconds != nil {
			query.Set("graceSeconds", strconv.FormatInt(*params.GraceSeconds, 10))
		}
	}
	path := "/v1/api-keys/" + url.PathEscape(keyID) + "/rotate"
	result, _, err := iris.Result[RotateAPIKeyResponseContent](ctx, a.client, http.MethodPost, path, query, headers, nil)
	return result, err
}

// ListProducts calls GET /v1/products.
func (a *API) ListProducts(ctx context.Context) (*ListProductsResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/products"
	result, _, err := iris.Result[ListProductsResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// DeleteSubjectGrant calls DELETE /v1/subject-grants.
//
// Deprecated: use subject grant set operations.
func (a *API) DeleteSubjectGrant(ctx context.Context, params *DeleteSubjectGrantParams) error {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ManagementId != nil {
			query.Set("managementId", *params.ManagementId)
		}
		query.Set("issuer", params.Issuer)
		query.Set("subject", params.Subject)
		if params.Tenant != nil {
			query.Set("tenant", *params.Tenant)
		}
		query.Set("product", string(params.Product))
		query.Set("resourceType", params.ResourceType)
		query.Set("resourceId", params.ResourceId)
	}
	path := "/v1/subject-grants"
	_, err := a.client.Call(ctx, http.MethodDelete, path, query, headers, nil, nil)
	return err
}

// ListSubjectGrants calls GET /v1/subject-grants.
//
// Deprecated: use subject grant set operations.
func (a *API) ListSubjectGrants(ctx context.Context, params *ListSubjectGrantsParams) (*ListSubjectGrantsResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		query.Set("issuer", params.Issuer)
		query.Set("subject", params.Subject)
	}
	path := "/v1/subject-grants"
	result, _, err := iris.Result[ListSubjectGrantsResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// CreateSubjectGrantSet calls POST /v1/subject-grants.
func (a *API) CreateSubjectGrantSet(ctx context.Context, body CreateSubjectGrantSetJSONRequestBody) (*CreateSubjectGrantSetResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/subject-grants"
	result, _, err := iris.Result[CreateSubjectGrantSetResponseContent](ctx, a.client, http.MethodPost, path, query, headers, body)
	return result, err
}

// PutSubjectGrant calls PUT /v1/subject-grants.
//
// Deprecated: use subject grant set operations.
func (a *API) PutSubjectGrant(ctx context.Context, body PutSubjectGrantJSONRequestBody) (*PutSubjectGrantResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/subject-grants"
	result, _, err := iris.Result[PutSubjectGrantResponseContent](ctx, a.client, http.MethodPut, path, query, headers, body)
	return result, err
}

// ApplySubjectGrantsResult includes the revision required by conditional updates.
type ApplySubjectGrantsResult struct {
	ApplySubjectGrantsResponseContent
	ETag string
}

// ApplySubjectGrants calls POST /v1/subject-grants/apply.
func (a *API) ApplySubjectGrants(ctx context.Context, params *ApplySubjectGrantsParams, body ApplySubjectGrantsJSONRequestBody) (*ApplySubjectGrantsResult, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ManagementId != nil {
			query.Set("managementId", *params.ManagementId)
		}
	}
	path := "/v1/subject-grants/apply"
	result, responseHeaders, err := iris.Result[ApplySubjectGrantsResponseContent](ctx, a.client, http.MethodPost, path, query, headers, body)
	if err != nil {
		return nil, err
	}
	return &ApplySubjectGrantsResult{ApplySubjectGrantsResponseContent: *result, ETag: responseHeaders.Get("ETag")}, nil
}

// SearchSubjectGrantSets calls POST /v1/subject-grants/search.
func (a *API) SearchSubjectGrantSets(ctx context.Context, body SearchSubjectGrantSetsJSONRequestBody) (*SearchSubjectGrantSetsResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/subject-grants/search"
	result, _, err := iris.Result[SearchSubjectGrantSetsResponseContent](ctx, a.client, http.MethodPost, path, query, headers, body)
	return result, err
}

// DeleteSubjectGrantSet calls DELETE /v1/subject-grants/{id}.
func (a *API) DeleteSubjectGrantSet(ctx context.Context, id string) error {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/subject-grants/" + url.PathEscape(id)
	_, err := a.client.Call(ctx, http.MethodDelete, path, query, headers, nil, nil)
	return err
}

// GetSubjectGrantSet calls GET /v1/subject-grants/{id}.
func (a *API) GetSubjectGrantSet(ctx context.Context, id string) (*GetSubjectGrantSetResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/subject-grants/" + url.PathEscape(id)
	result, _, err := iris.Result[GetSubjectGrantSetResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// DeleteProductSubjectGrants calls DELETE /v1/subject-grants/{id}/{product}.
func (a *API) DeleteProductSubjectGrants(ctx context.Context, id string, product Product) error {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/subject-grants/" + url.PathEscape(id) + "/" + url.PathEscape(string(product))
	_, err := a.client.Call(ctx, http.MethodDelete, path, query, headers, nil, nil)
	return err
}

// GetProductSubjectGrantsResult includes the revision required by conditional updates.
type GetProductSubjectGrantsResult struct {
	GetProductSubjectGrantsResponseContent
	ETag string
}

// GetProductSubjectGrants calls GET /v1/subject-grants/{id}/{product}.
func (a *API) GetProductSubjectGrants(ctx context.Context, id string, product Product) (*GetProductSubjectGrantsResult, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/subject-grants/" + url.PathEscape(id) + "/" + url.PathEscape(string(product))
	result, responseHeaders, err := iris.Result[GetProductSubjectGrantsResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	if err != nil {
		return nil, err
	}
	return &GetProductSubjectGrantsResult{GetProductSubjectGrantsResponseContent: *result, ETag: responseHeaders.Get("ETag")}, nil
}

// PatchProductSubjectGrantsResult includes the revision required by conditional updates.
type PatchProductSubjectGrantsResult struct {
	PatchProductSubjectGrantsResponseContent
	ETag string
}

// PatchProductSubjectGrants calls PATCH /v1/subject-grants/{id}/{product}.
func (a *API) PatchProductSubjectGrants(ctx context.Context, id string, product Product, params *PatchProductSubjectGrantsParams, body PatchProductSubjectGrantsJSONRequestBody) (*PatchProductSubjectGrantsResult, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.ManagementId != nil {
			query.Set("managementId", *params.ManagementId)
		}
		if params.IfMatch != nil {
			headers.Set("If-Match", *params.IfMatch)
		}
	}
	path := "/v1/subject-grants/" + url.PathEscape(id) + "/" + url.PathEscape(string(product))
	result, responseHeaders, err := iris.Result[PatchProductSubjectGrantsResponseContent](ctx, a.client, http.MethodPatch, path, query, headers, body)
	if err != nil {
		return nil, err
	}
	return &PatchProductSubjectGrantsResult{PatchProductSubjectGrantsResponseContent: *result, ETag: responseHeaders.Get("ETag")}, nil
}

// ListIdentityProviders calls GET /v1/tenants/{tenant}/identity-providers.
func (a *API) ListIdentityProviders(ctx context.Context, tenant string, params *ListIdentityProvidersParams) (*ListIdentityProvidersResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	if params != nil {
		if params.Limit != nil {
			query.Set("limit", strconv.FormatInt(int64(*params.Limit), 10))
		}
		if params.Cursor != nil {
			query.Set("cursor", *params.Cursor)
		}
	}
	path := "/v1/tenants/" + url.PathEscape(tenant) + "/identity-providers"
	result, _, err := iris.Result[ListIdentityProvidersResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// CreateIdentityProvider calls POST /v1/tenants/{tenant}/identity-providers.
func (a *API) CreateIdentityProvider(ctx context.Context, tenant string, body CreateIdentityProviderJSONRequestBody) (*CreateIdentityProviderResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/tenants/" + url.PathEscape(tenant) + "/identity-providers"
	result, _, err := iris.Result[CreateIdentityProviderResponseContent](ctx, a.client, http.MethodPost, path, query, headers, body)
	return result, err
}

// ReplaceIdentityProviders calls PUT /v1/tenants/{tenant}/identity-providers.
func (a *API) ReplaceIdentityProviders(ctx context.Context, tenant string, body ReplaceIdentityProvidersJSONRequestBody) (*ReplaceIdentityProvidersResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/tenants/" + url.PathEscape(tenant) + "/identity-providers"
	result, _, err := iris.Result[ReplaceIdentityProvidersResponseContent](ctx, a.client, http.MethodPut, path, query, headers, body)
	return result, err
}

// DeleteIdentityProvider calls DELETE /v1/tenants/{tenant}/identity-providers/{name}.
func (a *API) DeleteIdentityProvider(ctx context.Context, tenant string, name string) error {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/tenants/" + url.PathEscape(tenant) + "/identity-providers/" + url.PathEscape(name)
	_, err := a.client.Call(ctx, http.MethodDelete, path, query, headers, nil, nil)
	return err
}

// GetIdentityProvider calls GET /v1/tenants/{tenant}/identity-providers/{name}.
func (a *API) GetIdentityProvider(ctx context.Context, tenant string, name string) (*GetIdentityProviderResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/tenants/" + url.PathEscape(tenant) + "/identity-providers/" + url.PathEscape(name)
	result, _, err := iris.Result[GetIdentityProviderResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// PutIdentityProvider calls PUT /v1/tenants/{tenant}/identity-providers/{name}.
func (a *API) PutIdentityProvider(ctx context.Context, tenant string, name string, body PutIdentityProviderJSONRequestBody) (*PutIdentityProviderResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/tenants/" + url.PathEscape(tenant) + "/identity-providers/" + url.PathEscape(name)
	result, _, err := iris.Result[PutIdentityProviderResponseContent](ctx, a.client, http.MethodPut, path, query, headers, body)
	return result, err
}

// GetTenantSettings calls GET /v1/tenants/{tenant}/settings.
func (a *API) GetTenantSettings(ctx context.Context, tenant string) (*GetTenantSettingsResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/tenants/" + url.PathEscape(tenant) + "/settings"
	result, _, err := iris.Result[GetTenantSettingsResponseContent](ctx, a.client, http.MethodGet, path, query, headers, nil)
	return result, err
}

// PutTenantSettings calls PUT /v1/tenants/{tenant}/settings.
func (a *API) PutTenantSettings(ctx context.Context, tenant string, body PutTenantSettingsJSONRequestBody) (*PutTenantSettingsResponseContent, error) {
	query := url.Values{}
	headers := http.Header{}
	path := "/v1/tenants/" + url.PathEscape(tenant) + "/settings"
	result, _, err := iris.Result[PutTenantSettingsResponseContent](ctx, a.client, http.MethodPut, path, query, headers, body)
	return result, err
}
