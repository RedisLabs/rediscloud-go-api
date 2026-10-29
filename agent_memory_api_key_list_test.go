package rediscloud_api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/RedisLabs/rediscloud-go-api/redis"
	"github.com/RedisLabs/rediscloud-go-api/service/agentmemory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAgentMemory_ListAPIKeysResponseShape(t *testing.T) {
	for _, tc := range []struct {
		name, payload string
		malformed     bool
		populated     bool
	}{
		{name: "missing array", payload: `{}`, malformed: true},
		{name: "null array", payload: `{"apiKeys":null}`, malformed: true},
		{name: "empty array", payload: `{"apiKeys":[]}`},
		{name: "populated array", payload: `{"apiKeys":[{"apiKeyId":"key-1","keyName":"runtime"}]}`, populated: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, http.MethodGet, r.Method)
				assert.Equal(t, "/memory-stores/store-1/api-keys", r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				_, err := w.Write([]byte(tc.payload))
				assert.NoError(t, err)
			}))
			defer server.Close()
			subject, err := clientFromTestServer(server, "fake-key", "fake-secret")
			require.NoError(t, err)

			keys, err := subject.AgentMemory.ListAPIKeys(context.Background(), "store-1")
			if tc.malformed {
				require.Error(t, err)
				assert.Nil(t, keys)
			} else {
				require.NoError(t, err)
				require.NotNil(t, keys)
				if tc.populated {
					require.Len(t, keys, 1)
					assert.Equal(t, "key-1", redis.StringValue(keys[0].APIKeyID))
				} else {
					assert.Empty(t, keys)
				}
			}

			key, err := subject.AgentMemory.GetAPIKey(context.Background(), "store-1", "key-1")
			var notFound *agentmemory.APIKeyNotFound
			switch {
			case tc.malformed:
				require.Error(t, err)
				assert.Nil(t, key)
				assert.False(t, errors.As(err, &notFound), "malformed responses must not establish absence")
			case tc.populated:
				require.NoError(t, err)
				require.NotNil(t, key)
				assert.Equal(t, "key-1", redis.StringValue(key.APIKeyID))
			default:
				assert.Nil(t, key)
				require.True(t, errors.As(err, &notFound))
				assert.Equal(t, "store-1", notFound.StoreID)
				assert.Equal(t, "key-1", notFound.APIKeyID)
			}
		})
	}
}

func TestAgentMemory_CreateAPIKeyRejectsMalformedPreList(t *testing.T) {
	for _, payload := range []string{`{}`, `{"apiKeys":null}`} {
		t.Run(payload, func(t *testing.T) {
			var posts atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPost {
					posts.Add(1)
					http.Error(w, "unexpected POST", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				_, err := w.Write([]byte(payload))
				assert.NoError(t, err)
			}))
			defer server.Close()
			subject, err := clientFromTestServer(server, "fake-key", "fake-secret")
			require.NoError(t, err)

			key, err := subject.AgentMemory.CreateAPIKey(context.Background(), "store-1", "runtime")
			require.Error(t, err)
			assert.Nil(t, key)
			assert.Zero(t, posts.Load(), "must not create a key without a valid pre-creation listing")
		})
	}
}
