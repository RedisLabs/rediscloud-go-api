package latest_backups

import (
	"fmt"
)

// TODO Add LatestBackupStatus model and conversion from internal.task

type NotFound struct {
	subId int
	dbId  int
}

func (f *NotFound) Error() string {
	return fmt.Sprintf("database %d in subscription %d not found", f.dbId, f.subId)
}
