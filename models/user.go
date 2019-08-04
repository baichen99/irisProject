package models

import (
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

// User model
type User struct {
	Base
	Role	 string		`json:"role"`	// user, admin
	Username string    `json:"username" gorm:"NOT NULL"`
	Password string    `json:"-"` // ignore this filed for safety
}

// BeforeCreate generate user uuid
func (user *User) BeforeCreate(scope *gorm.Scope) (err error){
	err = scope.SetColumn("ID", uuid.NewV4())
	return
}
