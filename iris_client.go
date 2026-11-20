package rediscloud_api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/RedisLabs/rediscloud-go-api/internal/iris"
	"github.com/RedisLabs/rediscloud-go-api/service/irisagentmemory"
	"github.com/RedisLabs/rediscloud-go-api/service/irisidentity"
)

// IrisServiceConfig identifies an Iris service. BaseURL is the origin (without
// /v1 or /iris/v1/agent-memory); Token is its control-plane bearer credential.
type IrisServiceConfig struct {
	BaseURL string
	Token   string
}

type IrisConfig struct {
	AgentMemory IrisServiceConfig
	Identity    IrisServiceConfig
	HTTPClient  *http.Client
}

// IrisClient connects to Iris control planes independently of Redis Cloud.
type IrisClient struct {
	AgentMemory *irisagentmemory.API
	Identity    *irisidentity.API
}

// IrisHTTPError exposes HTTP failures through errors.As without logging secrets.
type IrisHTTPError = iris.HTTPError

// NewIrisClient configures either or both services. Unconfigured connectors are nil.
func NewIrisClient(config IrisConfig) (*IrisClient, error) {
	if config.AgentMemory.BaseURL == "" && config.Identity.BaseURL == "" {
		return nil, errors.New("at least one Iris service base URL is required")
	}
	client := &IrisClient{}
	var httpClient iris.HTTPClient
	if config.HTTPClient != nil {
		httpClient = config.HTTPClient
	}
	if config.AgentMemory.BaseURL != "" {
		transport, err := iris.NewClient(config.AgentMemory.BaseURL, config.AgentMemory.Token, httpClient)
		if err != nil {
			return nil, fmt.Errorf("configure Iris Agent Memory: %w", err)
		}
		client.AgentMemory = irisagentmemory.NewAPI(transport)
	}
	if config.Identity.BaseURL != "" {
		transport, err := iris.NewClient(config.Identity.BaseURL, config.Identity.Token, httpClient)
		if err != nil {
			return nil, fmt.Errorf("configure Iris Identity: %w", err)
		}
		client.Identity = irisidentity.NewAPI(transport)
	}
	return client, nil
}
