package pricing

import (
	"fmt"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type Pricing struct {
	DatabaseName        *string  `json:"databaseName,omitempty"`
	Type                *string  `json:"type,omitempty"`
	TypeDetails         *string  `json:"typeDetails,omitempty"`
	Quantity            *int     `json:"quantity,omitempty"`
	QuantityMeasurement *string  `json:"quantityMeasurement,omitempty"`
	PricePerUnit        *float64 `json:"pricePerUnit,omitempty"`
	PriceCurrency       *string  `json:"priceCurrency,omitempty"`
	PricePeriod         *string  `json:"pricePeriod,omitempty"`
	Region              *string  `json:"region,omitempty"`
}

func (o Pricing) String() string {
	return internal.ToString(o)
}

type NotFound struct {
	subId int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("subscription %d not found", f.subId)
}
