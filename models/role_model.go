package models

import (
	"time"
)

type Role struct {
	ID        int       `sql:"primary_key" json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `sql:"index" json:"deleted_at,omitempty"`
	Name      string    `json:"name,omitempty"`
}