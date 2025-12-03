package attachments

import "fmt"

type GetAttachmentsTask struct {
	CommandType *string   `json:"commandType,omitempty"`
	Description *string   `json:"description,omitempty"`
	Status      *string   `json:"status,omitempty"`
	ID          *string   `json:"taskId,omitempty"`
	Response    *Response `json:"response,omitempty"`
}

type Response struct {
	ResourceId *int      `json:"resourceId,omitempty"`
	Resource   *Resource `json:"resource,omitempty"`
}

type Resource struct {
	TransitGatewayAttachment []*TransitGatewayAttachment `json:"tgws,omitempty"`
}

type TransitGatewayAttachment struct {
	Id               *int    `json:"id,omitempty"`
	AwsTgwUid        *string `json:"awsTgwUid,omitempty"`
	AttachmentUid    *string `json:"attachmentUid,omitempty"`
	Status           *string `json:"status,omitempty"`
	AttachmentStatus *string `json:"attachmentStatus,omitempty"`
	AwsAccountId     *string `json:"awsAccountId,omitempty"`
	Cidrs            []*Cidr `json:"cidrs,omitempty"`
}

type Cidr struct {
	CidrAddress *string `json:"cidrAddress,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type updateCidrs struct {
	Cidrs *[]*string `json:"cidrs,omitempty"`
}

type TransitGatewayInvitation struct {
	Id               *int    `json:"id,omitempty"`
	Name             *string `json:"name,omitempty"`
	ResourceShareUid *string `json:"resourceShareUid,omitempty"`
	AwsAccountId     *string `json:"awsAccountId,omitempty"`
	Status           *string `json:"status,omitempty"`
	SharedDate       *string `json:"sharedDate,omitempty"`
}

type InvitationsResource struct {
	Resources []*TransitGatewayInvitation `json:"resources,omitempty"`
}

type InvitationsResponse struct {
	CommandType *string              `json:"commandType,omitempty"`
	Description *string              `json:"description,omitempty"`
	Status      *string              `json:"status,omitempty"`
	ID          *string              `json:"taskId,omitempty"`
	Response    *InvitationResponseData `json:"response,omitempty"`
}

type InvitationResponseData struct {
	ResourceId *int                 `json:"resourceId,omitempty"`
	Resource   *InvitationsResource `json:"resource,omitempty"`
}

type NotFound struct {
	subId int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("subscription %d not found", f.subId)
}

type NotFoundActiveActive struct {
	subId    int
	regionId int
}

func (f *NotFoundActiveActive) Error() string {
	return fmt.Sprintf("subscription %d in region %d not found", f.subId, f.regionId)
}
