package models

import (
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type Base struct {
	ID       uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	CreateAt time.Time `json:"create_at"`
}

// BeforeCreate generate user uuid
func (*Base) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("ID", uuid.NewV4())
	err = scope.SetColumn("CreateAt", time.Now())
	return
}
