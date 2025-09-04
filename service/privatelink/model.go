package privatelink

import "fmt"

type CreatePrivateLink struct {
	SubscriptionId *int    `json:"subscriptionId"`
	PrincipalId    *int    `json:"principal,omitempty"`
	PrincipalType  *string `json:"type,omitempty"`
	PrincipalAlias *string `json:"alias,omitempty"`
}

type PrivateLink struct {
	SubscriptionId *int `json:"subscriptionId"`
}

type CreatePrivateLinkActiveActive struct {
	SubscriptionId *int    `json:"subscriptionId"`
	PrincipalId    *int    `json:"principal,omitempty"`
	PrincipalType  *string `json:"type,omitempty"`
	PrincipalAlias *string `json:"alias,omitempty"`
}

type PrivateLinkActiveActive struct {
	SubscriptionId *int    `json:"subscriptionId"`
	RegionId       *string `json:"region_id"`
}

type NotFound struct {
	subscriptionID int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("resource not found - subscription %d", f.subscriptionID)
}

type NotFoundActiveActive struct {
	subscriptionID int
	regionID       int
}

func (f *NotFoundActiveActive) Error() string {
	return fmt.Sprintf("resource not found - subscription %d and region %d", f.subscriptionID, f.regionID)
}

const (
	// PrivateLinkStatusCreateQueued when PrivateLink creation is queued
	PrivateLinkStatusCreateQueued = "create-queued"
	// PrivateLinkStatusInitialized when PrivateLink provisioning started
	PrivateLinkStatusInitialized = "initialized"
	// PrivateLinkStatusCreatePending when PrivateLink provisioning is completed but databases are pending update
	PrivateLinkStatusCreatePending = "create-pending"
	// PrivateLinkStatusActive when PrivateLink is ready
	PrivateLinkStatusActive = "active"
)
