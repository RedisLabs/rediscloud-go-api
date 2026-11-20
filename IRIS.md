# Iris control-plane connectors

`NewIrisClient` connects directly to Iris Agent Memory and Identity Service (IdS).
It is independent of `NewClient` and the existing Cloud `AgentMemory` connector.
All existing APIs remain unchanged, and no additional Go dependencies are needed.

```go
import (
    "context"
    "os"

    rediscloud_api "github.com/RedisLabs/rediscloud-go-api"
    "github.com/RedisLabs/rediscloud-go-api/service/irisagentmemory"
)

client, err := rediscloud_api.NewIrisClient(rediscloud_api.IrisConfig{
    AgentMemory: rediscloud_api.IrisServiceConfig{
        BaseURL: "https://memory.example.com",
        Token: os.Getenv("IRIS_MEMORY_TOKEN"),
    },
    Identity: rediscloud_api.IrisServiceConfig{
        BaseURL: "https://identity.example.com",
        Token: os.Getenv("IRIS_IDENTITY_TOKEN"),
    },
})
if err != nil {
    return err
}

store, err := client.AgentMemory.CreateStore(context.Background(), nil,
    irisagentmemory.CreateStoreJSONRequestBody{
        Name: "support",
        Database: irisagentmemory.DatabaseRef{DatabaseId: "database-1"},
    })
```

Configure either or both services; an omitted service leaves its connector nil.
URLs are service origins, optionally with a reverse-proxy prefix. Operations add
`/iris/v1/agent-memory` for AM and `/v1` for IdS. Each service receives its own
`Authorization: Bearer` credential. `IrisConfig.HTTPClient` supports a custom
transport, timeouts, or authentication transport; request contexts support cancellation.
These connectors do not read Redis Cloud credential environment variables.

AM supports store CRUD, extractors, custom memory types, reference data, health,
and on-prem ingestion hooks. IdS supports key creation/read/update/revoke/rotation,
product capabilities, subject grants and grant sets, external identity providers,
and tenant settings. Runtime/data-plane APIs and Cloud-only AM operations are
outside these control-plane connectors.

Methods follow the service operation names and take typed request bodies and
optional `*OperationParams`. Pass nil to use server defaults. Pointer fields
preserve explicit zero/false values. Methods return synchronous service responses;
a newly created store may still have `State == PROVISIONING`. There is no Cloud
task polling or implicit retry of a mutation.

List results expose `NextCursor` (AM/identity providers) or `NextToken` (IdS keys).
Pass that value into the next call's params. Lists return one page at a time.
Grant GET/PATCH/apply results expose `ETag`; send it through `IfMatch` to make a
conditional patch. Revision fields stay in the request bodies or query params
specified by each operation. A nil resource entry in
`irisidentity.ResourceSubjectGrantPatches` encodes JSON `null` to remove a grant.
Use an empty, non-nil slice/map when explicitly clearing required collections.

Use `errors.As(err, &httpErr)` with `var httpErr *rediscloud_api.IrisHTTPError`
to inspect `StatusCode`, `Header`, and the raw RFC 9457 `Body`. The error string
omits the body, and the connectors do not log requests or returned tokens.

The models come from Iris commit
`23ff21223ab3b7a06c37e5fc1c4e1a3564f6d442` and its generated Smithy OpenAPI
outputs. Regenerate them from this checkout's development environment with
`python3 scripts/generate-iris-models.py /path/to/iris`. The script pins
oapi-codegen v2.8.2 and verifies the three source hashes before generating models.
To use an existing generator, set `OAPI_CODEGEN` to its executable path.
When updating contracts, review both model and service method changes before
updating the hashes in the script.
