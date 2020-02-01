package models

import "github.com/lib/pq"

// user teacher 多对多
type Teacher struct {
	Base
	Name       string         `gorm:"NOT NULL" json:"name"`
	StudentsID pq.StringArray `gorm:"TYPE:uuid[];NOT NULL;DEFAULT:array[]::uuid[];"`
}
