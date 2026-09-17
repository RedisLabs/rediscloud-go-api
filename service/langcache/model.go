package langcache

import "fmt"

const (
	SearchStrategyExact    = "exact"
	SearchStrategySemantic = "semantic"

	StatusProvisioning = "PROVISIONING"
	StatusReady        = "READY"
	StatusUnavailable  = "UNAVAILABLE"
)

// CreateCache contains the configuration accepted by the LangCache Admin API.
type CreateCache struct {
	Name                   string         `json:"name"`
	Database               DatabaseConfig `json:"database"`
	EmbeddingModel         EmbeddingModel `json:"embeddingModel"`
	DefaultSearchThreshold float64        `json:"defaultSearchThreshold"`
	DefaultTTLMillis       int64          `json:"defaultTtlMillis"`
	Attributes             []string       `json:"attributes"`
	SearchStrategies       []string       `json:"searchStrategies,omitempty"`
}

type DatabaseConfig struct {
	DatabaseID string `json:"databaseId"`
}

type EmbeddingModel struct {
	Provider      string                       `json:"provider"`
	Name          string                       `json:"name"`
	APIKey        *string                      `json:"apiKey,omitempty"`
	CustomOptions *CustomEmbeddingModelOptions `json:"customOptions,omitempty"`
}

type CustomEmbeddingModelOptions struct {
	BaseURL    string `json:"baseUrl"`
	Dimensions int32  `json:"dimensions"`
}

type CreateCacheResponse struct {
	CacheID string `json:"cacheId"`
}

type UpdateCache struct {
	Name                   *string         `json:"name,omitempty"`
	DefaultSearchThreshold *float64        `json:"defaultSearchThreshold,omitempty"`
	DefaultTTLMillis       *int64          `json:"defaultTtlMillis,omitempty"`
	EmbeddingModel         *EmbeddingModel `json:"embeddingModel,omitempty"`
	Attributes             *[]string       `json:"attributes,omitempty"`
	SearchStrategies       *[]string       `json:"searchStrategies,omitempty"`
}

type AddAPIKeyRequest struct {
	KeyName string `json:"keyName"`
}

type AddAPIKeyResponse struct {
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

// Cache is the configuration returned by the LangCache Admin API.
type Cache struct {
	CacheID                string    `json:"cacheId"`
	CacheName              string    `json:"cacheName"`
	DatabaseID             string    `json:"databaseId"`
	DefaultSearchThreshold float64   `json:"defaultSearchThreshold"`
	DefaultTTLMillis       int64     `json:"defaultTtlMillis"`
	Attributes             []string  `json:"attributes"`
	SearchStrategies       []string  `json:"searchStrategies"`
	Endpoint               string    `json:"endpoint"`
	EmbeddingProvider      Embedding `json:"embeddingProvider"`
	Status                 string    `json:"status"`
}

type Embedding struct {
	Provider string `json:"provider"`
	Model    string `json:"model"`
}

type NotFound struct {
	CacheID string
}

type APIKeyNotFound struct {
	CacheID  string
	APIKeyID string
}

func (e *APIKeyNotFound) Error() string {
	return fmt.Sprintf("API key %q for LangCache %q not found", e.APIKeyID, e.CacheID)
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("LangCache %q not found", e.CacheID)
}
