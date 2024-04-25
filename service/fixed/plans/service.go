package plans

import (
	"context"
	"fmt"
)

const root = "/fixed/plans"

type Log interface {
	Printf(format string, args ...interface{})
}

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
}

type API struct {
	client HttpClient
	logger Log
}

func NewAPI(client HttpClient, logger Log) *API {
	return &API{client: client, logger: logger}
}

func (a *API) List(ctx context.Context) ([]*GetPlanResponse, error) {
	return a.list(ctx, root)
}

func (a *API) ListWithProvider(ctx context.Context, provider string) ([]*GetPlanResponse, error) {
	address := fmt.Sprintf("%s?provider=%s", root, provider)
	return a.list(ctx, address)
}

// List will list all the plans available to the current account (filtered by provider if given).
func (a *API) list(ctx context.Context, address string) ([]*GetPlanResponse, error) {
	var response ListPlansResponse

	err := a.client.Get(ctx, "list fixed plans", address, &response)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Listing fixed plans, there are %d available", len(response.Plans))

	return response.Plans, nil
}
