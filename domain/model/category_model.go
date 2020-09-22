package model

import (
	"html"
	"strings"
	"time"
)

type Category struct {
	ID        int       `sql:"primary_key" json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Posts     []Post    `json:"posts,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `sql:"index" json:"deleted_at,omitempty"`
}

func (r *Category) Prepare() {
	r.Name = html.EscapeString(strings.TrimSpace(r.Name))
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
}

func (r *Category) Validate(action string) map[string]string {
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
