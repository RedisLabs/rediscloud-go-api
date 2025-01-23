package psc

import (
	"fmt"
)

type PrivateServiceConnectService struct {
	ID                    *int    `json:"id,omitempty"`
	ConnectionHostName    *string `json:"connectionHostName,omitempty"`
	ServiceAttachmentName *string `json:"serviceAttachmentName,omitempty"`
	Status                *string `json:"status,omitempty"`
}

type CreatePrivateServiceConnectEndpoint struct {
	GCPProjectID           *string `json:"gcpProjectId,omitempty"`
	GCPVPCName             *string `json:"gcpVpcName,omitempty"`
	GCPVPCSubnetName       *string `json:"gcpVpcSubnetName,omitempty"`
	EndpointConnectionName *string `json:"endpointConnectionName,omitempty"`
}

type UpdatePrivateServiceConnectEndpoint struct {
	GCPProjectID           *string `json:"gcpProjectId,omitempty"`
	GCPVPCName             *string `json:"gcpVpcName,omitempty"`
	GCPVPCSubnetName       *string `json:"gcpVpcSubnetName,omitempty"`
	EndpointConnectionName *string `json:"endpointConnectionName,omitempty"`
	Action                 *string `json:"action,omitempty"`
}

type PrivateServiceConnectEndpoints struct {
	PSCServiceID *int                             `json:"pscServiceId,omitempty"`
	Endpoints    []*PrivateServiceConnectEndpoint `json:"endpoints,omitempty"`
}
type PrivateServiceConnectEndpoint struct {
	ID                     *int    `json:"id,omitempty"`
	GCPProjectID           *string `json:"gcpProjectId,omitempty"`
	GCPVPCName             *string `json:"gcpVpcName,omitempty"`
	GCPVPCSubnetName       *string `json:"gcpVpcSubnetName,omitempty"`
	EndpointConnectionName *string `json:"endpointConnectionName,omitempty"`
	Status                 *string `json:"status,omitempty"`
}

type CreationScript struct {
	Script *GCPCreationScript `json:"script,omitempty"`
}

type DeletionScript struct {
	Script *GCPDeletionScript `json:"script,omitempty"`
}

type GCPCreationScript struct {
	Bash         *string       `json:"bash,omitempty"`
	Powershell   *string       `json:"powershell,omitempty"`
	TerraformGcp *TerraformGCP `json:"terraformGcp,omitempty"`
}

type GCPDeletionScript struct {
	Bash       *string `json:"bash,omitempty"`
	Powershell *string `json:"powershell,omitempty"`
}

type TerraformGCP struct {
	ServiceAttachments []TerraformGCPServiceAttachment `json:"serviceAttachments,omitempty"`
}

type TerraformGCPServiceAttachment struct {
	Name               *string `json:"name,omitempty"`
	DNSRecord          *string `json:"dnsRecord,omitempty"`
	IPAddressName      *string `json:"ipAddressName,omitempty"`
	ForwardingRuleName *string `json:"forwardingRuleName,omitempty"`
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

	// ServiceStatusCreateQueued when PSC service creation is queued
	ServiceStatusCreateQueued = "create-queued"
	// ServiceStatusDeleteQueued when PSC service deletion is queued
	ServiceStatusDeleteQueued = "delete-queued"
	// ServiceStatusInitialized when PSC service provisioning started
	ServiceStatusInitialized = "initialized"
	// ServiceStatusCreatePending when PSC service provisioning completed but databases are pending update
	ServiceStatusCreatePending = "create-pending"
	// ServiceStatusActive when PSC service is ready
	ServiceStatusActive = "active"
	// ServiceStatusDeletePending when infrastructure deletion is completed but databases are pending update
	ServiceStatusDeletePending = "delete-pending"
	// ServiceStatusDeleted when PSC service is deleted
	ServiceStatusDeleted = "deleted"
	// ServiceStatusProvisionFailed when PSC service has failed while creation/deletion
	ServiceStatusProvisionFailed = "provision-failed"
	// ServiceStatusFailed when PSC service failed after it's been reported as active
	ServiceStatusFailed = "failed"

	// EndpointStatusInitialized the endpoint was created in the SM but the creation script wasn't yet run
	EndpointStatusInitialized = "initialized"
	// EndpointStatusProcessing Processing the status during deletion or creation of 40 endpoints in cloud provider
	EndpointStatusProcessing = "processing"
	// EndpointStatusPending the endpoint is waiting for the user to accept or reject it
	EndpointStatusPending = "pending"
	// EndpointStatusAcceptPending the user accepted. the endpoint is not yet fully accepted
	EndpointStatusAcceptPending = "accept-pending"
	// EndpointStatusActive the endpoint is ready for use
	EndpointStatusActive = "active"
	// EndpointStatusDeleted the endpoint was successfully deleted
	EndpointStatusDeleted = "deleted"
	// EndpointStatusRejected the endpoint was successfully rejected
	EndpointStatusRejected = "rejected"
	// EndpointStatusRejectPending the user rejected. the endpoint is not yet fully rejected
	EndpointStatusRejectPending = "reject-pending"
	// EndpointStatusFailed endpoint is in error status
	EndpointStatusFailed = "failed"

	// EndpointActionAccept accepts the endpoint
	EndpointActionAccept = "accept"
	// EndpointActionReject rejects the endpoint
	EndpointActionReject = "reject"
)
