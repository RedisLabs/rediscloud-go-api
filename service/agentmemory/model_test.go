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

func TestCreateStoreRequestPresence(t *testing.T) {
	for _, tc := range []struct {
		name string
		body CreateStore
		want string
	}{
		{name: "omitted", want: `{}`},
		{name: "explicit-empty-zero", body: CreateStore{Name: redis.String(""), DatabaseID: redis.Int(0)}, want: `{"name":"","databaseId":0}`},
		{name: "configured", body: CreateStore{Name: redis.String("store"), DatabaseID: redis.Int(123)}, want: `{"name":"store","databaseId":123}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.body)
			require.NoError(t, err)
			assert.JSONEq(t, tc.want, string(encoded))
		})
	}
}

func TestExclusionDetectorListPresence(t *testing.T) {
	for _, tc := range []struct {
		name string
		body any
		want string
	}{
		{name: "builtin-nil", body: BuiltInDetectorsConfig{Enabled: redis.Bool(true)}, want: `{"enabled":true,"detectors":null}`},
		{name: "builtin-empty", body: BuiltInDetectorsConfig{Enabled: redis.Bool(true), Detectors: []DetectorSelection{}}, want: `{"enabled":true,"detectors":[]}`},
		{name: "custom-nil", body: CustomDetectorsConfig{Enabled: redis.Bool(true)}, want: `{"enabled":true,"detectors":null}`},
		{name: "custom-empty", body: CustomDetectorsConfig{Enabled: redis.Bool(true), Detectors: []CustomDetector{}}, want: `{"enabled":true,"detectors":[]}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.body)
			require.NoError(t, err)
			assert.JSONEq(t, tc.want, string(encoded))
		})
	}
}
