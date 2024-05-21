package latest_imports

import (
	"fmt"
	"regexp"
	"time"

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
	ID       *int      `json:"resourceId,omitempty"`
	Resource *Resource `json:"resource,omitempty"`
	Error    *Error    `json:"error,omitempty"`
}

func (o Response) String() string {
	return internal.ToString(o)
}

type Resource struct {
	Status         *string    `json:"status,omitempty"`
	LastImportTime *time.Time `json:"lastImportTime,omitempty"`
	FailureReason  *string    `json:"failureReason,omitempty"`
	// FailureReasonParams // Type unknown
}

func (o Resource) String() string {
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

type NotFound struct {
	subId int
	dbId  int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("database %d in subscription %d not found", f.dbId, f.subId)
}
