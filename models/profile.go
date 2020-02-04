package models

import uuid "github.com/satori/go.uuid"

type Profile struct {
	Base
	CreatorID uuid.UUID `json:"-"`
	Content   string    `json:"content"`
	Creator   User      `json:"creator"`
}
