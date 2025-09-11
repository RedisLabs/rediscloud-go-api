package privatelink

import "fmt"

type CreatePrivateLink struct {
	ShareName      *string `json:"shareName,omitempty"`
	Principal      *string `json:"principal,omitempty"`
	PrincipalType  *string `json:"type,omitempty"`
	PrincipalAlias *string `json:"alias,omitempty"`
}

type PrivateLink struct {
	Status                   *string                  `json:"status,omitempty"`
	Principals               []*PrivateLinkPrincipal  `json:"principals,omitempty"`
	ResourceConfigurationId  *int                     `json:"resourceConfigurationId,omitempty"`
	ResourceConfigurationArn *string                  `json:"resourceConfigurationArn,omitempty"`
	ShareArn                 *string                  `json:"shareArn,omitempty"`
	ShareName                *string                  `json:"shareName,omitempty"`
	Connections              []*PrivateLinkConnection `json:"connections,omitempty"`
	Databases                []*PrivateLinkDatabase   `json:"databases,omitempty"`
	SubscriptionId           *int                     `json:"subscriptionId,omitempty"`
	RegionId                 *int                     `json:"regionId,omitempty"`
	ErrorMessage             *string                  `json:"errorMessage,omitempty"`
}

type PrivateLinkPrincipal struct {
	Principal *string `json:"principal,omitempty"`
	Type      *string `json:"type,omitempty"`
	Alias     *string `json:"alias,omitempty"`
	Status    *string `json:"status,omitempty"`
}

type PrivateLinkConnection struct {
	AssociationId   *string `json:"associationId,omitempty"`
	ConnectionId    *int    `json:"connectionId,omitempty"`
	Type            *string `json:"type,omitempty"`
	OwnerId         *int    `json:"ownerId,omitempty"`
	AssociationDate *string `json:"associationDate,omitempty"`
}

type PrivateLinkDatabase struct {
	DatabaseId           *int    `json:"databaseId,omitempty"`
	Port                 *int    `json:"port,omitempty"`
	ResourceLinkEndpoint *string `json:"rlEndpoint,omitempty"`
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
	// PrivateLinkStatusInitializing when PrivateLink is initialising
	PrivateLinkStatusInitializing = "initializing"
	// PrivateLinkStatusDeleted when PrivateLink has been deleted
	PrivateLinkStatusDeleted = "deleting"
	// PrivateLinkStatusActive when PrivateLink is ready
	PrivateLinkStatusActive = "active"
)
