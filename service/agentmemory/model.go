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
	Name                     string                    `json:"name"`
	DatabaseID               int                       `json:"databaseId"`
	ShortMemory              *ShortMemoryConfig        `json:"shortMemory,omitempty"`
	LongTermMemory           *LongTermMemoryConfig     `json:"longTermMemory,omitempty"`
	LLM                      *ModelConfig              `json:"llm,omitempty"`
	ExtractionStrategy       string                    `json:"extractionStrategy,omitempty"`
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
	StoreID string `json:"storeId"`
}

type CreateAPIKeyRequest struct {
	KeyName string `json:"keyName"`
}

type CreateAPIKeyResponse struct {
	APIKey string `json:"apiKey"`
}

type ListAPIKeysResponse struct {
	APIKeys []APIKeyInfo `json:"apiKeys"`
}

type APIKeyInfo struct {
	APIKeyID        string `json:"apiKeyId"`
	KeyName         string `json:"keyName"`
	ObfuscatedToken string `json:"obfuscatedToken"`
	CreatedAt       int64  `json:"createdAt"`
}

type CreatedAPIKey struct {
	APIKey string
	Info   APIKeyInfo
}

type ListStoresResponse struct {
	Data []StoreSummary `json:"data"`
}

type StoreSummary struct {
	StoreID            string               `json:"storeId"`
	Name               string               `json:"name"`
	CreatedAt          *time.Time           `json:"createdAt,omitempty"`
	DatabaseID         string               `json:"databaseId"`
	DatabaseName       string               `json:"databaseName"`
	SubscriptionID     string               `json:"subscriptionId"`
	ExtractionStrategy string               `json:"extractionStrategy"`
	Summarization      *SummarizationConfig `json:"summarization,omitempty"`
	LLM                *ModelStatus         `json:"llm,omitempty"`
	Status             string               `json:"status"`
}

type Store struct {
	StoreID                  string                    `json:"storeId"`
	Name                     string                    `json:"name"`
	CreatedAt                *time.Time                `json:"createdAt,omitempty"`
	DatabaseID               string                    `json:"databaseId"`
	DatabaseName             string                    `json:"databaseName"`
	DatabaseProvider         string                    `json:"databaseProvider"`
	DatabaseRegion           string                    `json:"databaseRegion"`
	DatabaseAccessCIDRs      []string                  `json:"databaseAccessCidrs,omitempty"`
	SubscriptionID           string                    `json:"subscriptionId"`
	ShortMemory              *ShortMemoryConfig        `json:"shortMemory,omitempty"`
	LongTermMemory           *LongTermMemoryStatus     `json:"longTermMemory,omitempty"`
	LLM                      *ModelStatus              `json:"llm,omitempty"`
	Endpoint                 string                    `json:"endpoint"`
	Endpoints                []Endpoint                `json:"endpoints,omitempty"`
	ExtractionStrategy       string                    `json:"extractionStrategy"`
	Summarization            *SummarizationConfig      `json:"summarization,omitempty"`
	LongTermMemoryExclusions *LongTermMemoryExclusions `json:"longTermMemoryExclusions,omitempty"`
	CustomMemoryTypes        []CustomMemoryType        `json:"customMemoryTypes,omitempty"`
	ExtractionCadence        *ExtractionCadenceConfig  `json:"extractionCadence,omitempty"`
	Status                   string                    `json:"status"`
	ErrorMessage             string                    `json:"errorMessage,omitempty"`
	ErrorMessageTimestamp    *float64                  `json:"errorMessageTimestamp,omitempty"`
}

type ShortMemoryConfig struct {
	TTLSeconds int `json:"ttlSeconds"`
}

type LongTermMemoryConfig struct {
	TTLSeconds int          `json:"ttlSeconds,omitempty"`
	Embedding  *ModelConfig `json:"embedding,omitempty"`
}

type LongTermMemoryStatus struct {
	TTLSeconds int          `json:"ttlSeconds,omitempty"`
	Embedding  *ModelStatus `json:"embedding,omitempty"`
}

type ModelConfig struct {
	Provider    string            `json:"provider,omitempty"`
	Model       string            `json:"model,omitempty"`
	Credentials *ModelCredentials `json:"credentials,omitempty"`
}

type ModelCredentials struct {
	Type   string `json:"type"`
	APIKey string `json:"apiKey,omitempty"`
}

type ModelStatus struct {
	Provider string `json:"provider,omitempty"`
	Model    string `json:"model,omitempty"`
}

type SummarizationConfig struct {
	Enabled         bool                          `json:"enabled"`
	TriggerStrategy string                        `json:"triggerStrategy,omitempty"`
	EventCount      *SummarizationThresholdConfig `json:"eventCount,omitempty"`
}

type SummarizationThresholdConfig struct {
	Threshold   int `json:"threshold"`
	RetainCount int `json:"retainCount"`
}

type LongTermMemoryExclusions struct {
	Enabled          bool                    `json:"enabled"`
	Semantic         *SemanticExclusion      `json:"semantic,omitempty"`
	BuiltInDetectors *BuiltInDetectorsConfig `json:"builtInDetectors,omitempty"`
	CustomDetectors  *CustomDetectorsConfig  `json:"customDetectors,omitempty"`
}

type SemanticExclusion struct {
	Enabled bool   `json:"enabled"`
	Prompt  string `json:"prompt,omitempty"`
}

type BuiltInDetectorsConfig struct {
	Enabled   bool                `json:"enabled"`
	Detectors []DetectorSelection `json:"detectors,omitempty"`
}

type DetectorSelection struct {
	ID      string `json:"id"`
	Enabled bool   `json:"enabled"`
	Action  string `json:"action,omitempty"`
}

type CustomDetectorsConfig struct {
	Enabled   bool             `json:"enabled"`
	Detectors []CustomDetector `json:"detectors,omitempty"`
}

type CustomDetector struct {
	Name    string       `json:"name"`
	Enabled bool         `json:"enabled"`
	Action  string       `json:"action,omitempty"`
	Matcher *MatcherSpec `json:"matcher,omitempty"`
}

type MatcherSpec struct {
	Kind  string     `json:"kind"`
	Regex *RegexSpec `json:"regex,omitempty"`
}

type RegexSpec struct {
	Pattern string `json:"pattern"`
}

type CustomMemoryType struct {
	Name               string                    `json:"name"`
	Description        string                    `json:"description"`
	Fields             []CustomField             `json:"fields,omitempty"`
	ExtractionStrategy *CustomExtractionStrategy `json:"extractionStrategy,omitempty"`
}

type CustomField struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type CustomExtractionStrategy struct {
	Prompt  string `json:"prompt,omitempty"`
	Enabled *bool  `json:"enabled,omitempty"`
}

type MemoryTypeStrategyUpdate struct {
	TypeName string `json:"typeName"`
	Prompt   string `json:"prompt,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

type ExtractionCadenceConfig struct {
	ActiveIntervalSeconds int `json:"activeIntervalSeconds"`
}

type Endpoint struct {
	URL          string   `json:"url"`
	Provider     string   `json:"provider"`
	Region       string   `json:"region"`
	EgressIPs    []string `json:"egressIps,omitempty"`
	IsAccessible bool     `json:"isAccessible"`
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
