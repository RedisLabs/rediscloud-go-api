package rediscloud_api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/service/agentmemory"
	"github.com/stretchr/testify/require"
)

func TestAgentMemoryMalformedKeyMetadataDoesNotBecomeNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"apiKeys":[{"keyName":"runtime"}]}`))
	}))
	defer server.Close()
	client, err := NewClient(BaseURL(server.URL), Auth("fake-key", "fake-secret"))
	require.NoError(t, err)
	_, err = client.AgentMemory.GetAPIKey(context.Background(), "store-1", "key-1")
	require.Error(t, err)
	var notFound *agentmemory.APIKeyNotFound
	require.False(t, errors.As(err, &notFound), "invalid metadata must not remove Terraform ownership")
}

func TestAgentMemoryKeyCreationRejectsUnresolvableMetadataWithoutRetry(t *testing.T) {
	for _, corruptBeforeCreate := range []bool{true, false} {
		t.Run(map[bool]string{true: "before-create", false: "after-create"}[corruptBeforeCreate], func(t *testing.T) {
			var posts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				if r.Method == http.MethodPost {
					posts.Add(1)
					_, _ = w.Write([]byte(`{"apiKey":"fake-created-secret"}`))
					return
				}
				if corruptBeforeCreate || posts.Load() > 0 {
					_, _ = w.Write([]byte(`{"apiKeys":[{"keyName":"runtime"}]}`))
				} else {
					_, _ = w.Write([]byte(`{"apiKeys":[]}`))
				}
			}))
			defer server.Close()
			client, err := NewClient(BaseURL(server.URL), Auth("fake-key", "fake-secret"))
			require.NoError(t, err)
			key, err := client.AgentMemory.CreateAPIKey(context.Background(), "store-1", "runtime")
			require.Error(t, err)
			require.Nil(t, key)
			if corruptBeforeCreate {
				require.Zero(t, posts.Load())
			} else {
				require.EqualValues(t, 1, posts.Load())
			}
		})
	}
}
