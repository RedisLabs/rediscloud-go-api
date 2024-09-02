package latest_backups

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

func fromInternal(t internal.Task) LatestBackupStatus {
	lbs := LatestBackupStatus{
		CommandType: t.CommandType,
		Description: t.Description,
		Status:      t.Status,
		ID:          t.ID,
	}
	if t.Response == nil {
		lbs.Response = nil
	} else {
		response := Response{
			ID: t.Response.ID,
		}

		if t.Response.Resource == nil {
			response.Resource = nil
		} else {
			j, err := t.Response.Resource.MarshalJSON()
			if err != nil {
				panic(nil)
			}

			var res Resource
			if j != nil && len(j) > 0 {
				err = json.Unmarshal(j, &res)
				if err != nil {
					panic(nil)
				}
			}

			response.Resource = &res
		}

		response.Error = t.Response.Error
		lbs.Response = &response
	}
	return lbs
}

type LatestBackupStatus struct {
	CommandType *string   `json:"commandType,omitempty"`
	Description *string   `json:"description,omitempty"`
	Status      *string   `json:"status,omitempty"`
	ID          *string   `json:"taskId,omitempty"`
	Response    *Response `json:"response,omitempty"`
}

func (o LatestBackupStatus) String() string {
	return internal.ToString(o)
}

type Response struct {
	ID       *int            `json:"resourceId,omitempty"`
	Resource *Resource       `json:"resource,omitempty"`
	Error    *internal.Error `json:"error,omitempty"`
}

func (o Response) String() string {
	return internal.ToString(o)
}

type Resource struct {
	Status         *string    `json:"status,omitempty"`
	LastBackupTime *time.Time `json:"lastBackupTime,omitempty"`
	FailureReason  *string    `json:"failureReason,omitempty"`
}

func (o Resource) String() string {
	return internal.ToString(o)
}

type NotFound struct {
	subId int
	dbId  int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("database %d in subscription %d not found", f.dbId, f.subId)
}

type NotFoundActiveActive struct {
	subId  int
	dbId   int
	region string
}

func (f *NotFoundActiveActive) Error() string {
	return fmt.Sprintf("database %d in subscription %d in region %s not found", f.dbId, f.subId, f.region)
}
