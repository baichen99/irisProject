package models

type profile struct {
	ID      int    `json:"id" gorm:"primary_key"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}
