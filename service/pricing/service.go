package pricing

import (
	"context"
	"fmt"
)

type HttpClient interface {
	Get(ctx context.Context, name, path string, responseBody interface{}) error
}

type API struct {
	client HttpClient
}

func NewAPI(client HttpClient) *API {
	return &API{client: client}
}

// List will return the list of available pricing detail blocks for the provided subscription.
func (a *API) List(ctx context.Context, subscription int) ([]*Pricing, error) {
	var body ListPricingResponse

	message := fmt.Sprintf("get pricing information for subscription %d", subscription)
	address := fmt.Sprintf("/subscriptions/%d/pricing", subscription)

	if err := a.client.Get(ctx, message, address, &body); err != nil {
		return nil, err
	}

	return body.Pricing, nil
}
