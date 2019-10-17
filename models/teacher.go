package models

// user teacher 多对多
type Teacher struct {
	Base
	Name     string `gorm:"NOT NULL" json:"name"`
	Students []User `gorm:"MANY2MANY:"student_teacher" json:"students"`
}
