package tags

import "fmt"

type AllTags struct {
	Tags *[]*Tag `json:"tags,omitempty"`
}

type Tag struct {
	Key   *string `json:"key,omitempty"`
	Value *string `json:"value,omitempty"`
}

type NotFound struct {
	subId int
	dbId  int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("database %d in subscription %d not found", f.dbId, f.subId)
}
