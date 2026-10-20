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
		{name: "embedding only", memory: &LongTermMemoryConfig{Embedding: &ModelConfig{Provider: redis.String("openai"), Model: redis.String("embedding-model")}}, want: `{"embedding":{"provider":"openai","model":"embedding-model"}}`},
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

func TestOptionalStringRequestPresence(t *testing.T) {
	for _, field := range []struct {
		name string
		key  string
		body func(*string) any
	}{
		{"key name", "keyName", func(value *string) any { return CreateAPIKeyRequest{KeyName: value} }},
		{"credential type", "type", func(value *string) any { return ModelCredentials{Type: value} }},
		{"custom strategy type name", "typeName", func(value *string) any { return MemoryTypeStrategyUpdate{TypeName: value} }},
		{"created API key", "apiKey", func(value *string) any { return CreatedAPIKey{APIKey: value} }},
		{"extraction strategy", "extractionStrategy", func(value *string) any { return CreateStore{ExtractionStrategy: value} }},
		{"model provider", "provider", func(value *string) any { return ModelConfig{Provider: value} }},
		{"model name", "model", func(value *string) any { return ModelConfig{Model: value} }},
		{"model API key", "apiKey", func(value *string) any { return ModelCredentials{Type: redis.String("apiKey"), APIKey: value} }},
		{"custom strategy prompt", "prompt", func(value *string) any {
			return MemoryTypeStrategyUpdate{TypeName: redis.String("preference"), Prompt: value}
		}},
	} {
		t.Run(field.name, func(t *testing.T) {
			for _, tc := range []struct {
				name  string
				value *string
			}{
				{"omitted", nil},
				{"explicit empty", redis.String("")},
				{"configured", redis.String("configured-value")},
			} {
				t.Run(tc.name, func(t *testing.T) {
					encoded, err := json.Marshal(field.body(tc.value))
					require.NoError(t, err)
					var body map[string]json.RawMessage
					require.NoError(t, json.Unmarshal(encoded, &body))
					actual, present := body[field.key]
					if tc.value == nil {
						assert.False(t, present)
						return
					}
					require.True(t, present)
					var value string
					require.NoError(t, json.Unmarshal(actual, &value))
					assert.Equal(t, *tc.value, value)
				})
			}
		})
	}
}

func TestScalarJSONPresence(t *testing.T) {
	for _, tc := range []struct {
		name string
		body any
		want string
	}{
		{name: "omitted session TTL", body: ShortMemoryConfig{}, want: `{}`},
		{name: "zero session TTL", body: ShortMemoryConfig{TTLSeconds: redis.Int(0)}, want: `{"ttlSeconds":0}`},
		{name: "omitted exclusions enabled", body: LongTermMemoryExclusions{}, want: `{}`},
		{name: "disabled exclusions", body: LongTermMemoryExclusions{Enabled: redis.Bool(false)}, want: `{"enabled":false}`},
		{name: "omitted summarization threshold", body: SummarizationThresholdConfig{}, want: `{}`},
		{name: "zero summarization threshold", body: SummarizationThresholdConfig{Threshold: redis.Int(0), RetainCount: redis.Int(0)}, want: `{"threshold":0,"retainCount":0}`},
		{name: "omitted cadence", body: ExtractionCadenceConfig{}, want: `{}`},
		{name: "zero cadence", body: ExtractionCadenceConfig{ActiveIntervalSeconds: redis.Int(0)}, want: `{"activeIntervalSeconds":0}`},
		{name: "omitted key metadata", body: APIKeyInfo{}, want: `{}`},
		{name: "zero key timestamp", body: APIKeyInfo{CreatedAt: redis.Int64(0)}, want: `{"createdAt":0}`},
		{name: "omitted endpoint accessibility", body: Endpoint{}, want: `{}`},
		{name: "inaccessible endpoint", body: Endpoint{IsAccessible: redis.Bool(false)}, want: `{"isAccessible":false}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encoded, err := json.Marshal(tc.body)
			require.NoError(t, err)
			assert.JSONEq(t, tc.want, string(encoded))
		})
	}
}
