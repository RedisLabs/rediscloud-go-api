package latest_imports

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/RedisLabs/rediscloud-go-api/internal"
	"github.com/RedisLabs/rediscloud-go-api/redis"
)

type LatestImportStatus struct {
	CommandType *string   `json:"commandType,omitempty"`
	Description *string   `json:"description,omitempty"`
	Status      *string   `json:"status,omitempty"`
	ID          *string   `json:"taskId,omitempty"`
	Response    *Response `json:"response,omitempty"`
}

func (o LatestImportStatus) String() string {
	return internal.ToString(o)
}

type Response struct {
	ID       *int             `json:"resourceId,omitempty"`
	Resource *json.RawMessage `json:"resource,omitempty"`
	Error    *Error           `json:"error,omitempty"`
}

func (o Response) String() string {
	return internal.ToString(o)
}

type Error struct {
	Type        *string `json:"type,omitempty"`
	Description *string `json:"description,omitempty"`
	Status      *string `json:"status,omitempty"`
}

func (e Error) String() string {
	return internal.ToString(e)
}

func (e *Error) StatusCode() string {
	matches := errorStatusCode.FindStringSubmatch(redis.StringValue(e.Status))
	if len(matches) == 2 {
		return matches[1]
	}
	return ""
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s - %s: %s", redis.StringValue(e.Status), redis.StringValue(e.Type), redis.StringValue(e.Description))
}

var errorStatusCode = regexp.MustCompile("^(\\d*).*$")

func NewLatestImportStatus(task *internal.Task) *LatestImportStatus {
	latestImportStatus := LatestImportStatus{
		CommandType: task.CommandType,
		Description: task.Description,
		Status:      task.Status,
		ID:          task.ID,
	}

	if task.Response != nil {
		r := Response{
			ID:       task.Response.ID,
			Resource: task.Response.Resource,
		}

		if task.Response.Error != nil {
			e := Error{
				Type:        task.Response.Error.Type,
				Description: task.Response.Error.Description,
				Status:      task.Response.Error.Status,
			}

			r.Error = &e
		}

		latestImportStatus.Response = &r
	}

	return &latestImportStatus
}

type NotFound struct {
	subId int
	dbId  int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("database %d in subscription %d not found", f.dbId, f.subId)
}
