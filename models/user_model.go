package models

import (
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"time"
)

type User struct {
	ID          uuid.UUID `sql:"primary_key" json:"id,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	DeletedAt   time.Time `sql:"index" json:"deleted_at,omitempty"`
	FullName    string    `json:"full_name,omitempty"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	Email       string    `json:"email,omitempty"`
	Posts       []Post    `json:"posts,omitempty"`
	Role        Role      `json:"role,omitempty"`
	RoleId      int       `json:"role_id,omitempty"`
}