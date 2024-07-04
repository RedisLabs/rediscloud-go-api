package plan_subscriptions

import (
	"context"
	"fmt"
	"github.com/RedisLabs/rediscloud-go-api/service/fixed/plans"
	"net/url"
)

const root = "/fixed/plans/subscriptions"

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

// List will list all subscriptions with a specific plan
func (a *API) List(ctx context.Context, id int) ([]*plans.GetPlanResponse, error) {
	var response plans.ListPlansResponse

	err := a.client.Get(ctx, "list plans for subscription plans", fmt.Sprintf("%s/%d", root, id), &response)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Listing fixed plans, all cloud providers, there are %d available", len(response.Plans))

	return response.Plans, nil
}
