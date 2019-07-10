package models

// User model
type User struct {
	ID			int    `json:"id" gorm:"primary_key"`
	Username	string `json:"username" gorm:"NOT NULL"`
	Password	string `json:"-"`	// ignore this filed for safety
}
