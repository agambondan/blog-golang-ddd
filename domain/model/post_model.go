package model

import (
	"github.com/google/uuid"
	"time"
)

type Post struct {
	ID        uuid.UUID `sql:"primary_key" json:"id,omitempty"`
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	AuthorID  uuid.UUID `json:"author_id,omitempty"`
	Author    User      `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `sql:"index" json:"deleted_at,omitempty"`
}
