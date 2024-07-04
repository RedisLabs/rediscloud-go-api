package plan_subscriptions

import (
	"context"
	"fmt"

	"github.com/RedisLabs/rediscloud-go-api/service/fixed/plans"
)

const root = "/fixed/plans/subscriptions"

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

// List will list all plans upgradable from a given subscription
func (a *API) List(ctx context.Context, id int) ([]*plans.GetPlanResponse, error) {
	var response plans.ListPlansResponse

	path := fmt.Sprintf("%s/%d", root, id)
	err := a.client.Get(ctx, "list plans for subscription plans", path, &response)
	if err != nil {
		return nil, err
	}

	a.logger.Printf("Listing fixed plans applicable to subscription %d, there are %d available", id, len(response.Plans))

	return response.Plans, nil
}
