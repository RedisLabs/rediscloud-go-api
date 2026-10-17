package agentmemory

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/RedisLabs/rediscloud-go-api/redis"
)

func TestLongTermMemoryTTLRequestPresence(t *testing.T) {
	for _, tc := range []struct {
		name   string
		memory *LongTermMemoryConfig
		want   string
	}{
		{name: "memory omitted"},
		{name: "TTL omitted", memory: &LongTermMemoryConfig{}, want: `{}`},
		{name: "zero", memory: &LongTermMemoryConfig{TTLSeconds: redis.Int(0)}, want: `{"ttlSeconds":0}`},
		{name: "negative", memory: &LongTermMemoryConfig{TTLSeconds: redis.Int(-1)}, want: `{"ttlSeconds":-1}`},
		{name: "positive", memory: &LongTermMemoryConfig{TTLSeconds: redis.Int(86400)}, want: `{"ttlSeconds":86400}`},
		{name: "above previous maximum", memory: &LongTermMemoryConfig{TTLSeconds: redis.Int(31536001)}, want: `{"ttlSeconds":31536001}`},
		{name: "embedding only", memory: &LongTermMemoryConfig{Embedding: &ModelConfig{Provider: "openai", Model: "embedding-model"}}, want: `{"embedding":{"provider":"openai","model":"embedding-model"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, request := range []struct {
				name string
				body any
			}{
				{name: "create", body: CreateStore{LongTermMemory: tc.memory}},
				{name: "update", body: UpdateStore{LongTermMemory: tc.memory}},
			} {
				t.Run(request.name, func(t *testing.T) {
					encoded, err := json.Marshal(request.body)
					require.NoError(t, err)
					var body map[string]json.RawMessage
					require.NoError(t, json.Unmarshal(encoded, &body))
					memory, present := body["longTermMemory"]
					if tc.memory == nil {
						assert.False(t, present)
						return
					}
					require.True(t, present)
					assert.JSONEq(t, tc.want, string(memory))
				})
			}
		})
	}
}
