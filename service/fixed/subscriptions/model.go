package subscriptions

import (
	"fmt"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type FixedSubscription struct {
	ID              *int       `json:"id,omitempty"` // Omit for Create and Update
	Name            *string    `json:"name,omitempty"`
	Status          *string    `json:"status,omitempty"` // Omit for Create and Update
	PlanId          *int       `json:"planId,omitempty"`
	PaymentMethod   *string    `json:"paymentMethodType,omitempty"`
	PaymentMethodID *int       `json:"paymentMethodId,omitempty"`
	CreationDate    *time.Time `json:"creationDate,omitempty"` // Omit for Create and Update
}

func (o FixedSubscription) String() string {
	return internal.ToString(o)
}

type listFixedSubscriptionResponse struct {
	FixedSubscriptions []*FixedSubscription `json:"subscriptions"`
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
