package model

import (
	"html"
	"strings"
	"time"
)

type Role struct {
	ID        uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Name      string     `gorm:"size:100;not null;unique" json:"name"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

func (r *Role) BeforeSave() {
	r.Name = html.EscapeString(strings.TrimSpace(r.Name))
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
