package models

import (
	uuid "github.com/satori/go.uuid"
	"time"
)

type Base struct {
	ID			uuid.UUID	`gorm:"type:uuid;primary_key" json:"id"`
	CreateAt	time.Time	`json:"create_at"`
}
