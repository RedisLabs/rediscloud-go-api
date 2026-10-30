package agentmemory

import (
	"fmt"
	"time"
)

const (
	StatusProvisioning = "PROVISIONING"
	StatusReady        = "READY"
	StatusUnavailable  = "UNAVAILABLE"

	ExtractionStrategyInstruct = "INSTRUCT"

	SummarizationTriggerStrategyEventCount = "event_count"
)

type CreateStore struct {
	Name                     *string                   `json:"name,omitempty"`
	DatabaseID               *int                      `json:"databaseId,omitempty"`
	ShortMemory              *ShortMemoryConfig        `json:"shortMemory,omitempty"`
	LongTermMemory           *LongTermMemoryConfig     `json:"longTermMemory,omitempty"`
	LLM                      *ModelConfig              `json:"llm,omitempty"`
	ExtractionStrategy       *string                   `json:"extractionStrategy,omitempty"`
	Summarization            *SummarizationConfig      `json:"summarization,omitempty"`
	LongTermMemoryExclusions *LongTermMemoryExclusions `json:"longTermMemoryExclusions,omitempty"`
	CustomMemoryTypes        []CustomMemoryType        `json:"customMemoryTypes,omitempty"`
	ExtractionCadence        *ExtractionCadenceConfig  `json:"extractionCadence,omitempty"`
}

type UpdateStore struct {
	Name                             *string                    `json:"name,omitempty"`
	ShortMemory                      *ShortMemoryConfig         `json:"shortMemory,omitempty"`
	LongTermMemory                   *LongTermMemoryConfig      `json:"longTermMemory,omitempty"`
	LLM                              *ModelConfig               `json:"llm,omitempty"`
	Summarization                    *SummarizationConfig       `json:"summarization,omitempty"`
	LongTermMemoryExclusions         *LongTermMemoryExclusions  `json:"longTermMemoryExclusions,omitempty"`
	AddCustomMemoryTypes             []CustomMemoryType         `json:"addCustomMemoryTypes,omitempty"`
	UpdateCustomMemoryTypeStrategies []MemoryTypeStrategyUpdate `json:"updateCustomMemoryTypeStrategies,omitempty"`
	ExtractionCadence                *ExtractionCadenceConfig   `json:"extractionCadence,omitempty"`
}

type CreateStoreResponse struct {
	StoreID *string `json:"storeId,omitempty"`
}

type CreateAPIKeyRequest struct {
	KeyName *string `json:"keyName,omitempty"`
}

type CreateAPIKeyResponse struct {
	APIKey *string `json:"apiKey,omitempty"`
}

type ListAPIKeysResponse struct {
	APIKeys []APIKeyInfo `json:"apiKeys"`
}

type APIKeyInfo struct {
	APIKeyID        *string `json:"apiKeyId,omitempty"`
	KeyName         *string `json:"keyName,omitempty"`
	ObfuscatedToken *string `json:"obfuscatedToken,omitempty"`
	CreatedAt       *int64  `json:"createdAt,omitempty"`
}

type CreatedAPIKey struct {
	APIKey *string     `json:"apiKey,omitempty"`
	Info   *APIKeyInfo `json:"info,omitempty"`
}

type ListStoresResponse struct {
	Data []StoreSummary `json:"data"`
}

type StoreSummary struct {
	StoreID            *string              `json:"storeId,omitempty"`
	Name               *string              `json:"name,omitempty"`
	CreatedAt          *time.Time           `json:"createdAt,omitempty"`
	DatabaseID         *string              `json:"databaseId,omitempty"`
	DatabaseName       *string              `json:"databaseName,omitempty"`
	SubscriptionID     *string              `json:"subscriptionId,omitempty"`
	ExtractionStrategy *string              `json:"extractionStrategy,omitempty"`
	Summarization      *SummarizationConfig `json:"summarization,omitempty"`
	LLM                *ModelStatus         `json:"llm,omitempty"`
	Status             *string              `json:"status,omitempty"`
}

type Store struct {
	StoreID                  *string                   `json:"storeId,omitempty"`
	Name                     *string                   `json:"name,omitempty"`
	CreatedAt                *time.Time                `json:"createdAt,omitempty"`
	DatabaseID               *string                   `json:"databaseId,omitempty"`
	DatabaseName             *string                   `json:"databaseName,omitempty"`
	DatabaseProvider         *string                   `json:"databaseProvider,omitempty"`
	DatabaseRegion           *string                   `json:"databaseRegion,omitempty"`
	DatabaseAccessCIDRs      []string                  `json:"databaseAccessCidrs,omitempty"`
	SubscriptionID           *string                   `json:"subscriptionId,omitempty"`
	ShortMemory              *ShortMemoryConfig        `json:"shortMemory,omitempty"`
	LongTermMemory           *LongTermMemoryStatus     `json:"longTermMemory,omitempty"`
	LLM                      *ModelStatus              `json:"llm,omitempty"`
	Endpoint                 *string                   `json:"endpoint,omitempty"`
	Endpoints                []Endpoint                `json:"endpoints,omitempty"`
	ExtractionStrategy       *string                   `json:"extractionStrategy,omitempty"`
	Summarization            *SummarizationConfig      `json:"summarization,omitempty"`
	LongTermMemoryExclusions *LongTermMemoryExclusions `json:"longTermMemoryExclusions,omitempty"`
	CustomMemoryTypes        []CustomMemoryType        `json:"customMemoryTypes,omitempty"`
	ExtractionCadence        *ExtractionCadenceConfig  `json:"extractionCadence,omitempty"`
	Status                   *string                   `json:"status,omitempty"`
	ErrorMessage             *string                   `json:"errorMessage,omitempty"`
	ErrorMessageTimestamp    *float64                  `json:"errorMessageTimestamp,omitempty"`
}

type ShortMemoryConfig struct {
	TTLSeconds *int `json:"ttlSeconds,omitempty"`
}

type LongTermMemoryConfig struct {
	TTLSeconds *int         `json:"ttlSeconds,omitempty"`
	Embedding  *ModelConfig `json:"embedding,omitempty"`
}

type LongTermMemoryStatus struct {
	TTLSeconds *int         `json:"ttlSeconds,omitempty"`
	Embedding  *ModelStatus `json:"embedding,omitempty"`
}

type ModelConfig struct {
	Provider    *string           `json:"provider,omitempty"`
	Model       *string           `json:"model,omitempty"`
	Credentials *ModelCredentials `json:"credentials,omitempty"`
}

type ModelCredentials struct {
	Type   *string `json:"type,omitempty"`
	APIKey *string `json:"apiKey,omitempty"`
}

type ModelStatus struct {
	Provider *string `json:"provider,omitempty"`
	Model    *string `json:"model,omitempty"`
}

type SummarizationConfig struct {
	Enabled         *bool                         `json:"enabled,omitempty"`
	TriggerStrategy *string                       `json:"triggerStrategy,omitempty"`
	EventCount      *SummarizationThresholdConfig `json:"eventCount,omitempty"`
}

type SummarizationThresholdConfig struct {
	Threshold   *int `json:"threshold,omitempty"`
	RetainCount *int `json:"retainCount,omitempty"`
}

type LongTermMemoryExclusions struct {
	Enabled          *bool                   `json:"enabled,omitempty"`
	Semantic         *SemanticExclusion      `json:"semantic,omitempty"`
	BuiltInDetectors *BuiltInDetectorsConfig `json:"builtInDetectors,omitempty"`
	CustomDetectors  *CustomDetectorsConfig  `json:"customDetectors,omitempty"`
}

type SemanticExclusion struct {
	Enabled *bool   `json:"enabled,omitempty"`
	Prompt  *string `json:"prompt,omitempty"`
}

type BuiltInDetectorsConfig struct {
	Enabled *bool `json:"enabled,omitempty"`
	// Empty selections clear the list; nil is sent as null and preserved by CAPI.
	Detectors []DetectorSelection `json:"detectors"`
}

type DetectorSelection struct {
	ID      *string `json:"id,omitempty"`
	Enabled *bool   `json:"enabled,omitempty"`
	Action  *string `json:"action,omitempty"`
}

type CustomDetectorsConfig struct {
	Enabled *bool `json:"enabled,omitempty"`
	// Empty selections clear the list; nil is sent as null and preserved by CAPI.
	Detectors []CustomDetector `json:"detectors"`
}

type CustomDetector struct {
	Name    *string      `json:"name,omitempty"`
	Enabled *bool        `json:"enabled,omitempty"`
	Action  *string      `json:"action,omitempty"`
	Matcher *MatcherSpec `json:"matcher,omitempty"`
}

type MatcherSpec struct {
	Kind  *string    `json:"kind,omitempty"`
	Regex *RegexSpec `json:"regex,omitempty"`
}

type RegexSpec struct {
	Pattern *string `json:"pattern,omitempty"`
}

type CustomMemoryType struct {
	Name               *string                   `json:"name,omitempty"`
	Description        *string                   `json:"description,omitempty"`
	Fields             []CustomField             `json:"fields,omitempty"`
	ExtractionStrategy *CustomExtractionStrategy `json:"extractionStrategy,omitempty"`
}

type CustomField struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Type        *string `json:"type,omitempty"`
}

type CustomExtractionStrategy struct {
	Prompt  *string `json:"prompt,omitempty"`
	Enabled *bool   `json:"enabled,omitempty"`
}

type MemoryTypeStrategyUpdate struct {
	TypeName *string `json:"typeName,omitempty"`
	Prompt   *string `json:"prompt,omitempty"`
	Enabled  *bool   `json:"enabled,omitempty"`
}

type ExtractionCadenceConfig struct {
	ActiveIntervalSeconds *int `json:"activeIntervalSeconds,omitempty"`
}

type Endpoint struct {
	URL          *string  `json:"url,omitempty"`
	Provider     *string  `json:"provider,omitempty"`
	Region       *string  `json:"region,omitempty"`
	EgressIPs    []string `json:"egressIps,omitempty"`
	IsAccessible *bool    `json:"isAccessible,omitempty"`
}

type NotFound struct {
	StoreID string
}

type APIKeyNotFound struct {
	StoreID  string
	APIKeyID string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("Agent Memory store %q not found", e.StoreID)
}

func (e *APIKeyNotFound) Error() string {
	return fmt.Sprintf("API key %q for Agent Memory store %q not found", e.APIKeyID, e.StoreID)
}
