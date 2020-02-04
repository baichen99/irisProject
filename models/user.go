package models

import (
	"github.com/lib/pq"
)

// User model
type User struct {
	Base
	Role       string         `json:"role" gorm:"DEFAULT:'user'"` // user, admin
	Username   string         `json:"username" gorm:"NOT NULL"`
	Password   string         `json:"-"` // ignore this filed for safety
	TeachersID pq.StringArray `gorm:"TYPE:uuid[];NOT NULL;DEFAULT:array[]::uuid[];"`
}
