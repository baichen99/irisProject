package models

type Profile struct {
	Base
	Content string `json:"content"`
}
