package models

import uuid "github.com/satori/go.uuid"

// User model
type User struct {
	Base
	Role      string    `json:"role" gorm:"DEFAULT:'user'"` // user, admin
	Teachers  []Teacher `json:"teachers" gorm:"MANY2MANY:student_teacher;JOINTABLE_FOREIGNKEY:TeacherID"`
	Profile   Profile   `json:"profile" gorm:"FOREIGNKEY:ProfileID"`
	ProfileID uuid.UUID `json:"-" `
	Username  string    `json:"username" gorm:"NOT NULL"`
	Password  string    `json:"-"` // ignore this filed for safety
}
