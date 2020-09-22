package model

import (
	"html"
	"strings"
	"time"
)

type Role struct {
	ID        uint64    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `sql:"index" json:"deleted_at,omitempty"`
}

func (r *Role) Prepare() {
	r.Name = html.EscapeString(strings.TrimSpace(r.Name))
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
}

func (r *Role) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if r.Name == "" || r.Name == "null" {
			errorMessages["title_required"] = "title is required"
		}
	default:
		if r.Name == "" || r.Name == "null" {
			errorMessages["title_required"] = "title is required"
		}
	}
	return errorMessages
}
