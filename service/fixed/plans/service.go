package plans

import (
	"context"
	"net/url"
)

const root = "/fixed/plans"

type Log interface {
	Printf(format string, args ...interface{})
}

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
	GetWithQuery(ctx context.Context, name, path string, query url.Values, responseBody interface{}) error
}

type API struct {
	client HttpClient
	logger Log
}

func NewAPI(client HttpClient, logger Log) *API {
	return &API{client: client, logger: logger}
}

// List will list all the plans available to the current account
func (a *API) List(ctx context.Context) ([]*GetPlanResponse, error) {
	var response ListPlansResponse

	err := a.client.Get(ctx, "list fixed plans", root, &response)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Listing fixed plans, all cloud providers, there are %d available", len(response.Plans))

	return response.Plans, nil
}

// ListWithProvider will list all the plans available to the current account, filtered by provider
func (a *API) ListWithProvider(ctx context.Context, provider string) ([]*GetPlanResponse, error) {
	var response ListPlansResponse

	q := map[string][]string{
		"provider": {provider},
	}

	err := a.client.GetWithQuery(ctx, "list fixed plans", root, q, &response)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Listing fixed plans for cloud provider %s, there are %d available", provider, len(response.Plans))

	return response.Plans, nil
}
