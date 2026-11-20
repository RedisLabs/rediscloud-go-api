package rediscloud_api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/irisagentmemory"
	"github.com/RedisLabs/rediscloud-go-api/service/irisidentity"
)

type irisExchange struct {
	method, path, body, response string
	query                        url.Values
	status                       int
	headers                      http.Header
	responseHeaders              http.Header
}

func irisServer(t *testing.T, token string, exchanges ...irisExchange) *httptest.Server {
	t.Helper()
	var calls atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		i := int(calls.Add(1)) - 1
		if !assert.Less(t, i, len(exchanges), "unexpected request: %s %s", r.Method, r.URL) {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		want := exchanges[i]
		assert.Equal(t, want.method, r.Method)
		assert.Equal(t, want.path, r.URL.EscapedPath())
		assert.Equal(t, len(want.query), len(r.URL.Query()))
		for key, values := range want.query {
			assert.Equal(t, values, r.URL.Query()[key])
		}
		assert.Equal(t, "Bearer "+token, r.Header.Get("Authorization"))
		assert.Empty(t, r.Header.Get("X-Api-Key"))
		assert.Empty(t, r.Header.Get("X-Api-Secret-Key"))
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		for key := range want.headers {
			assert.Equal(t, want.headers.Get(key), r.Header.Get(key))
		}
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		if want.body != "" {
			assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
			assert.JSONEq(t, want.body, string(body))
		} else {
			assert.Empty(t, body)
		}
		for key, values := range want.responseHeaders {
			w.Header()[key] = values
		}
		status := want.status
		if status == 0 {
			status = http.StatusOK
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, want.response)
	}))
	t.Cleanup(func() {
		server.Close()
		assert.EqualValues(t, len(exchanges), calls.Load(), "unconsumed Iris requests")
	})
	return server
}

func TestIrisAgentMemoryStoreLifecycle(t *testing.T) {
	const store = `{"storeId":"store-1","name":"support","state":"PROVISIONING","createdAt":"2026-11-19T00:00:00Z","database":{"databaseId":"db-1"}}`
	server := irisServer(t, "am-token",
		irisExchange{method: "POST", path: "/proxy/iris/v1/agent-memory/stores", body: `{"name":"support","database":{"databaseId":"db-1"},"shortMemory":{"ttlSeconds":0}}`, response: store, status: 201},
		irisExchange{method: "GET", path: "/proxy/iris/v1/agent-memory/stores", query: url.Values{"limit": {"1"}}, response: `{"items":[` + store + `],"nextCursor":"a+/=&"}`},
		irisExchange{method: "GET", path: "/proxy/iris/v1/agent-memory/stores", query: url.Values{"cursor": {"a+/=&"}}, response: `{"items":[]}`},
		irisExchange{method: "PATCH", path: "/proxy/iris/v1/agent-memory/stores/store-1", body: `{"name":"renamed"}`, response: store},
		irisExchange{method: "DELETE", path: "/proxy/iris/v1/agent-memory/stores/store-1", query: url.Values{"flush": {"false"}}, status: 204},
	)
	client, err := NewIrisClient(IrisConfig{AgentMemory: IrisServiceConfig{BaseURL: server.URL + "/proxy/", Token: "am-token"}})
	require.NoError(t, err)
	assert.Nil(t, client.Identity)
	ctx := context.Background()
	ttl := int32(0)
	created, err := client.AgentMemory.CreateStore(ctx, nil, irisagentmemory.CreateStoreJSONRequestBody{
		Name: "support", Database: irisagentmemory.DatabaseRef{DatabaseId: "db-1"},
		ShortMemory: &irisagentmemory.ShortMemoryConfigInput{TtlSeconds: &ttl},
	})
	require.NoError(t, err)
	assert.Equal(t, "store-1", created.StoreId)
	assert.Equal(t, irisagentmemory.PROVISIONING, created.State)
	limit := int32(1)
	first, err := client.AgentMemory.ListStores(ctx, &irisagentmemory.ListStoresParams{Limit: &limit})
	require.NoError(t, err)
	require.Len(t, first.Items, 1)
	second, err := client.AgentMemory.ListStores(ctx, &irisagentmemory.ListStoresParams{Cursor: first.NextCursor})
	require.NoError(t, err)
	assert.Empty(t, second.Items)
	_, err = client.AgentMemory.UpdateStore(ctx, created.StoreId, nil, irisagentmemory.UpdateStoreJSONRequestBody{Name: redis.String("renamed")})
	require.NoError(t, err)
	require.NoError(t, client.AgentMemory.DeleteStore(ctx, created.StoreId, &irisagentmemory.DeleteStoreParams{Flush: redis.Bool(false)}))
}

func TestIrisIdentityKeyLifecycle(t *testing.T) {
	server := irisServer(t, "ids-token",
		irisExchange{method: "POST", path: "/v1/api-keys", body: `{"name":"runtime","expiresAt":0}`, response: `{"keyId":"key-1","token":"secret-once","createdAt":123}`, status: 201},
		irisExchange{method: "GET", path: "/v1/api-keys", query: url.Values{"product": {"memory"}, "nextToken": {"x+/=&"}}, response: `{"apiKeys":[{"keyId":"key-1","name":"runtime","tenantId":"tenant-1","createdAt":123,"grants":[]}],"nextToken":"page-2"}`},
		irisExchange{method: "PATCH", path: "/v1/api-keys/key-1", body: `{"expiresAt":0}`, response: `{"keyId":"key-1","expiresAt":0}`},
		irisExchange{method: "POST", path: "/v1/api-keys/key-1/rotate", query: url.Values{"graceSeconds": {"0"}}, response: `{"keyId":"key-2","token":"rotated-secret","createdAt":124,"oldExpiresAt":124}`, status: 201},
		irisExchange{method: "DELETE", path: "/v1/api-keys/key-2", status: 204},
	)
	client, err := NewIrisClient(IrisConfig{Identity: IrisServiceConfig{BaseURL: server.URL, Token: "ids-token"}})
	require.NoError(t, err)
	ctx := context.Background()
	created, err := client.Identity.CreateAPIKey(ctx, irisidentity.CreateAPIKeyJSONRequestBody{Name: "runtime", ExpiresAt: redis.Int64(0)})
	require.NoError(t, err)
	assert.Equal(t, "secret-once", created.Token)
	product := irisidentity.Memory
	listed, err := client.Identity.ListAPIKeys(ctx, &irisidentity.ListAPIKeysParams{Product: &product, NextToken: redis.String("x+/=&")})
	require.NoError(t, err)
	require.Len(t, listed.ApiKeys, 1)
	assert.Equal(t, "tenant-1", listed.ApiKeys[0].TenantId)
	assert.Equal(t, "page-2", *listed.NextToken)
	_, err = client.Identity.UpdateAPIKey(ctx, created.KeyId, irisidentity.UpdateAPIKeyJSONRequestBody{ExpiresAt: redis.Int64(0)})
	require.NoError(t, err)
	rotated, err := client.Identity.RotateAPIKey(ctx, created.KeyId, &irisidentity.RotateAPIKeyParams{GraceSeconds: redis.Int64(0)})
	require.NoError(t, err)
	assert.Equal(t, "rotated-secret", rotated.Token)
	require.NoError(t, client.Identity.RevokeAPIKey(ctx, rotated.KeyId))
}

func TestIrisIdentityConditionalGrantPatch(t *testing.T) {
	server := irisServer(t, "ids-token",
		irisExchange{method: "GET", path: "/v1/subject-grants/set-1/memory", response: `{"product":"memory","grants":{}}`, responseHeaders: http.Header{"Etag": {`"rev-1"`}}},
		irisExchange{method: "PATCH", path: "/v1/subject-grants/set-1/memory", query: url.Values{"managementId": {"controller-1"}}, headers: http.Header{"If-Match": {`"rev-1"`}}, body: `{"grants":{"store":{"remove-me":null,"store-1":{"actions":["read"]}}}}`, response: `{"product":"memory","grants":{"store":{"store-1":{"actions":["read"]}}}}`, responseHeaders: http.Header{"Etag": {`"rev-2"`}}},
		irisExchange{method: "PATCH", path: "/v1/subject-grants/set-1/memory", headers: http.Header{"If-Match": {`"rev-1"`}}, body: `{"grants":{}}`, status: 412, response: `{"type":"/errors/precondition-failed","status":412}`},
	)
	client, err := NewIrisClient(IrisConfig{Identity: IrisServiceConfig{BaseURL: server.URL, Token: "ids-token"}})
	require.NoError(t, err)
	ctx := context.Background()
	current, err := client.Identity.GetProductSubjectGrants(ctx, "set-1", irisidentity.Memory)
	require.NoError(t, err)
	assert.Equal(t, `"rev-1"`, current.ETag)
	updated, err := client.Identity.PatchProductSubjectGrants(ctx, "set-1", irisidentity.Memory, &irisidentity.PatchProductSubjectGrantsParams{IfMatch: &current.ETag, ManagementId: redis.String("controller-1")}, irisidentity.PatchProductSubjectGrantsJSONRequestBody{Grants: irisidentity.ResourceTypeSubjectGrantPatches{"store": {"remove-me": nil, "store-1": irisidentity.SubjectGrantEntry{Actions: []irisidentity.Action{irisidentity.Read}}}}})
	require.NoError(t, err)
	assert.Equal(t, `"rev-2"`, updated.ETag)
	assert.Equal(t, []irisidentity.Action{irisidentity.Read}, updated.Grants["store"]["store-1"].Actions)
	_, err = client.Identity.PatchProductSubjectGrants(ctx, "set-1", irisidentity.Memory, &irisidentity.PatchProductSubjectGrantsParams{IfMatch: &current.ETag}, irisidentity.PatchProductSubjectGrantsJSONRequestBody{Grants: irisidentity.ResourceTypeSubjectGrantPatches{}})
	var httpErr *IrisHTTPError
	require.ErrorAs(t, err, &httpErr)
	assert.Equal(t, 412, httpErr.StatusCode)
	assert.JSONEq(t, `{"type":"/errors/precondition-failed","status":412}`, string(httpErr.Body))
}

func TestIrisEscapingAndTenantRevisions(t *testing.T) {
	server := irisServer(t, "ids-token",
		irisExchange{method: "GET", path: "/v1/tenants/tenant%2Fa%25%3Fb/settings", response: `{"tenantId":"tenant/a%?b","revision":"r1","properties":{}}`},
		irisExchange{method: "PUT", path: "/v1/tenants/tenant%2Fa%25%3Fb/settings", body: `{"expectedRevision":"r1","properties":{"memory.ownership.enabled":false}}`, response: `{"tenantId":"tenant/a%?b","revision":"r2","properties":{"memory.ownership.enabled":false}}`},
		irisExchange{method: "PUT", path: "/v1/tenants/tenant-1/identity-providers", body: `{"items":[],"expectedRevision":"r2"}`, response: `{"items":[],"revision":"r3"}`},
	)
	client, err := NewIrisClient(IrisConfig{Identity: IrisServiceConfig{BaseURL: server.URL, Token: "ids-token"}})
	require.NoError(t, err)
	ctx := context.Background()
	settings, err := client.Identity.GetTenantSettings(ctx, "tenant/a%?b")
	require.NoError(t, err)
	updated, err := client.Identity.PutTenantSettings(ctx, "tenant/a%?b", irisidentity.PutTenantSettingsJSONRequestBody{ExpectedRevision: &settings.Revision, Properties: irisidentity.TenantProperties{"memory.ownership.enabled": json.RawMessage(`false`)}})
	require.NoError(t, err)
	assert.Equal(t, "r2", updated.Revision)
	providers, err := client.Identity.ReplaceIdentityProviders(ctx, "tenant-1", irisidentity.ReplaceIdentityProvidersJSONRequestBody{ExpectedRevision: redis.String(updated.Revision), Items: []irisidentity.IdentityProviderInput{}})
	require.NoError(t, err)
	assert.Empty(t, providers.Items)
}

func TestIrisClientValidatesConfiguration(t *testing.T) {
	_, err := NewIrisClient(IrisConfig{})
	require.Error(t, err)
	for _, baseURL := range []string{"localhost:8080", "file:///tmp/iris", "https://example.com?token=secret", "https://example.com/#fragment", "https://user:pass@example.com"} {
		t.Run(baseURL, func(t *testing.T) {
			_, err := NewIrisClient(IrisConfig{Identity: IrisServiceConfig{BaseURL: baseURL}})
			require.Error(t, err)
		})
	}
}

func TestIrisAgentMemoryRevisions(t *testing.T) {
	server := irisServer(t, "am-token",
		irisExchange{method: "PATCH", path: "/iris/v1/agent-memory/stores/store-1/extractors/ex-1", body: `{"expectedRevision":"r1","enabled":false,"destinationNamespaceId":""}`, response: `{"extractorId":"ex-1","enabled":false,"revision":"r2"}`},
		irisExchange{method: "DELETE", path: "/iris/v1/agent-memory/stores/store-1/extractors/ex-1", query: url.Values{"expectedRevision": {"r2"}}, status: 204},
		irisExchange{method: "PATCH", path: "/iris/v1/agent-memory/stores/store-1/memory-types/profile", body: `{"expectedRevision":"r3","extractionStrategy":{"enabled":false}}`, response: `{"name":"profile","revision":"r4","extractionStrategy":{"enabled":false}}`},
		irisExchange{method: "PUT", path: "/iris/v1/agent-memory/stores/store-1/ingestion-hooks", body: `{"enabled":false,"expectedRevision":"0"}`, response: `{"enabled":false,"revision":"r5"}`},
	)
	client, err := NewIrisClient(IrisConfig{AgentMemory: IrisServiceConfig{BaseURL: server.URL, Token: "am-token"}})
	require.NoError(t, err)
	ctx := context.Background()
	extractor, err := client.AgentMemory.UpdateExtractor(ctx, "store-1", "ex-1", nil, irisagentmemory.UpdateExtractorJSONRequestBody{
		ExpectedRevision: "r1", Enabled: redis.Bool(false), DestinationNamespaceId: redis.String(""),
	})
	require.NoError(t, err)
	assert.False(t, extractor.Enabled)
	require.NoError(t, client.AgentMemory.RemoveExtractor(ctx, "store-1", "ex-1", &irisagentmemory.RemoveExtractorParams{ExpectedRevision: extractor.Revision}))
	memoryType, err := client.AgentMemory.UpdateMemoryType(ctx, "store-1", "profile", nil, irisagentmemory.UpdateMemoryTypeJSONRequestBody{
		ExpectedRevision: redis.String("r3"), ExtractionStrategy: irisagentmemory.CustomExtractionStrategyUpdate{Enabled: redis.Bool(false)},
	})
	require.NoError(t, err)
	assert.Equal(t, "r4", redis.StringValue(memoryType.Revision))
	hooks, err := client.AgentMemory.PutIngestionHooks(ctx, "store-1", nil, irisagentmemory.PutIngestionHooksJSONRequestBody{Enabled: false, ExpectedRevision: "0"})
	require.NoError(t, err)
	assert.False(t, hooks.Enabled)
	assert.Equal(t, "r5", hooks.Revision)
}

func TestIrisServicesUseSeparateCredentials(t *testing.T) {
	memory := irisServer(t, "memory-only", irisExchange{method: "GET", path: "/iris/v1/agent-memory/health/ready", response: `{"status":"ok"}`})
	identity := irisServer(t, "identity-only", irisExchange{method: "GET", path: "/v1/products", response: `{"products":[]}`})
	client, err := NewIrisClient(IrisConfig{
		AgentMemory: IrisServiceConfig{BaseURL: memory.URL, Token: "memory-only"},
		Identity:    IrisServiceConfig{BaseURL: identity.URL, Token: "identity-only"},
		HTTPClient:  &http.Client{},
	})
	require.NoError(t, err)
	_, err = client.AgentMemory.GetServiceReadiness(context.Background())
	require.NoError(t, err)
	_, err = client.Identity.ListProducts(context.Background())
	require.NoError(t, err)
}

func TestIrisHTTPFailuresAreNotRetried(t *testing.T) {
	for _, status := range []int{401, 403, 404, 409, 412, 429, 500} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			server := irisServer(t, "ids-token", irisExchange{method: "POST", path: "/v1/api-keys", body: `{"name":"runtime"}`, status: status, response: `{"detail":"secret-response-content"}`, responseHeaders: http.Header{"Retry-After": {"60"}}})
			client, err := NewIrisClient(IrisConfig{Identity: IrisServiceConfig{BaseURL: server.URL, Token: "ids-token"}})
			require.NoError(t, err)
			result, err := client.Identity.CreateAPIKey(context.Background(), irisidentity.CreateAPIKeyJSONRequestBody{Name: "runtime"})
			assert.Nil(t, result)
			var httpErr *IrisHTTPError
			require.ErrorAs(t, err, &httpErr)
			assert.Equal(t, status, httpErr.StatusCode)
			assert.Equal(t, "60", httpErr.Header.Get("Retry-After"))
			assert.NotContains(t, err.Error(), "secret-response-content")
		})
	}
}

func TestIrisCancellationAndMalformedResponses(t *testing.T) {
	server := irisServer(t, "ids-token", irisExchange{method: "GET", path: "/v1/products", response: `{"products":`})
	client, err := NewIrisClient(IrisConfig{Identity: IrisServiceConfig{BaseURL: server.URL, Token: "ids-token"}})
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	result, err := client.Identity.ListProducts(ctx)
	assert.Nil(t, result)
	require.ErrorIs(t, err, context.Canceled)
	result, err = client.Identity.ListProducts(context.Background())
	assert.Nil(t, result)
	require.ErrorContains(t, err, "decode Iris response")
}
