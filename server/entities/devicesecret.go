package entities

import uuid "github.com/satori/go.uuid"

type DeviceSecret struct {
	UUID     uuid.UUID `sql:"type:uuid;"`
	UserUUID uuid.UUID `sql:"type:uuid;"`
	Secret   string
}
