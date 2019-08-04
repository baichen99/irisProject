package models

type profile struct {
	Base
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
}
