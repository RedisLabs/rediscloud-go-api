package agentmemory

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStoreResponseDistinguishesOmittedAndZero(t *testing.T) {
	var omitted Store
	require.NoError(t, json.Unmarshal([]byte(`{"shortMemory":{},"longTermMemory":{},"summarization":{"eventCount":{}},"longTermMemoryExclusions":{"semantic":{}},"endpoints":[{}]}`), &omitted))
	require.Nil(t, omitted.StoreID)
	require.Nil(t, omitted.ShortMemory.TTLSeconds)
	require.Nil(t, omitted.LongTermMemory.TTLSeconds)
	require.Nil(t, omitted.Summarization.Enabled)
	require.Nil(t, omitted.Summarization.EventCount.Threshold)
	require.Nil(t, omitted.LongTermMemoryExclusions.Semantic.Prompt)
	require.Nil(t, omitted.Endpoints[0].IsAccessible)

	var explicit Store
	require.NoError(t, json.Unmarshal([]byte(`{"storeId":"","shortMemory":{"ttlSeconds":0},"longTermMemory":{"ttlSeconds":0},"summarization":{"enabled":false,"eventCount":{"threshold":0,"retainCount":0}},"longTermMemoryExclusions":{"semantic":{"prompt":""}},"endpoints":[{"isAccessible":false}]}`), &explicit))
	require.NotNil(t, explicit.StoreID)
	require.Empty(t, *explicit.StoreID)
	require.NotNil(t, explicit.ShortMemory.TTLSeconds)
	require.Zero(t, *explicit.ShortMemory.TTLSeconds)
	require.NotNil(t, explicit.LongTermMemory.TTLSeconds)
	require.Zero(t, *explicit.LongTermMemory.TTLSeconds)
	require.NotNil(t, explicit.Summarization.Enabled)
	require.False(t, *explicit.Summarization.Enabled)
	require.NotNil(t, explicit.Summarization.EventCount.Threshold)
	require.Zero(t, *explicit.Summarization.EventCount.Threshold)
	require.NotNil(t, explicit.LongTermMemoryExclusions.Semantic.Prompt)
	require.Empty(t, *explicit.LongTermMemoryExclusions.Semantic.Prompt)
	require.NotNil(t, explicit.Endpoints[0].IsAccessible)
	require.False(t, *explicit.Endpoints[0].IsAccessible)
}

func TestAPIKeyResponseDistinguishesOmittedTimestamp(t *testing.T) {
	var omitted, explicit APIKeyInfo
	require.NoError(t, json.Unmarshal([]byte(`{}`), &omitted))
	require.Nil(t, omitted.CreatedAt)
	require.NoError(t, json.Unmarshal([]byte(`{"createdAt":0}`), &explicit))
	require.NotNil(t, explicit.CreatedAt)
	require.Zero(t, *explicit.CreatedAt)
}
