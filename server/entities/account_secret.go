package entities

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type AccountSecret struct {
	UUID      uuid.UUID `sql:"type:uuid;" gorm:"primary_key"`
	UserUUID  uuid.UUID `sql:"type:uuid;" gorm:"index:idx_user_uuid"`
	Secret    string
	CreatedAt time.Time // TODO: create desc index on this
}
