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
	Id               *string `json:"id,omitempty"`
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
