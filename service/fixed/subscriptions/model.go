package subscriptions

import (
	"fmt"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

// Represents request information to create/update a fixed subscription
type FixedSubscriptionRequest struct {
	Name            *string `json:"name,omitempty"`
	PlanId          *int    `json:"planId,omitempty"`
	PaymentMethod   *string `json:"paymentMethod,omitempty"`
	PaymentMethodID *int    `json:"paymentMethodId,omitempty"`
}

// Represents subscription info response from the get/list endpoints
type FixedSubscriptionResponse struct {
	ID              *int       `json:"id,omitempty"`
	Name            *string    `json:"name,omitempty"`
	Status          *string    `json:"status,omitempty"`
	PlanId          *int       `json:"planId,omitempty"`
	PaymentMethod   *string    `json:"paymentMethodType,omitempty"`
	PaymentMethodID *int       `json:"paymentMethodId,omitempty"`
	CreationDate    *time.Time `json:"creationDate,omitempty"`
}

func (o FixedSubscriptionRequest) String() string {
	return internal.ToString(o)
}

func (o FixedSubscriptionResponse) String() string {
	return internal.ToString(o)
}

type listFixedSubscriptionResponse struct {
	FixedSubscriptions []*FixedSubscriptionResponse `json:"subscriptions"`
}

type NotFound struct {
	ID int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("fixed subscription %d not found", f.ID)
}

const (
	// FixedSubscriptionStatusActive is the active value of the `Status` field in `Subscription`
	FixedSubscriptionStatusActive = "active"
	// FixedSubscriptionStatusPending is the pending value of the `Status` field in `Subscription`
	FixedSubscriptionStatusPending = "pending"
	// FixedSubscriptionStatusError is the error value of the `Status` field in `Subscription`
	FixedSubscriptionStatusError = "error"
	// FixedSubscriptionStatusDeleting is the deleting value of the `Status` field in `Subscription`
	FixedSubscriptionStatusDeleting = "deleting"
)
