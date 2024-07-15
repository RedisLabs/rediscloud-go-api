package maintenance

import (
	"fmt"

	"github.com/RedisLabs/rediscloud-go-api/internal"
)

type Maintenance struct {
	Mode    *string   `json:"mode,omitempty"`
	Windows []*Window `json:"windows,omitempty"`
}

func (o Maintenance) String() string {
	return internal.ToString(o)
}

type Window struct {
	StartHour       *int      `json:"startHour,omitempty"`
	DurationInHours *int      `json:"durationInHours,omitempty"`
	Days            []*string `json:"days,omitempty"`
}

func (o Window) String() string {
	return internal.ToString(o)
}

type NotFound struct {
	subId int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("maintenance in subscription %d not found", f.subId)
}
