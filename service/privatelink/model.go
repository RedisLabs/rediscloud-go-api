package privatelink

import "fmt"

type PrivateLinkConfig struct {
	ID                    *int    `json:"id,omitempty"`
	ConnectionHostName    *string `json:"connectionHostName,omitempty"`
	ServiceAttachmentName *string `json:"serviceAttachmentName,omitempty"`
	Status                *string `json:"status,omitempty"`
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
