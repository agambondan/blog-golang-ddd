package model

import (
	"github.com/google/uuid"
	"html"
	"strings"
	"time"
)

type Post struct {
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	UserUUID      uuid.UUID  `gorm:"type:uuid;not null;" json:"user_id"`
	Title       string     `gorm:"size:100;not null;unique" json:"title"`
	Description string     `gorm:"text;not null;" json:"description"`
	PostImage   string     `gorm:"size:255;null;" json:"Post_image"`
	Author      User       `json:"author,omitempty"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt   *time.Time `sql:"index" json:"deleted_at,omitempty"`
}

func (p *Post) BeforeSave() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
}

func (p *Post) Prepare() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Post) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)

	switch strings.ToLower(action) {
	case "update":
		if p.Title == "" || p.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if p.Description == "" || p.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	default:
		if p.Title == "" || p.Title == "null" {
			errorMessages["title_required"] = "title is required"
		}
		if p.Description == "" || p.Description == "null" {
			errorMessages["desc_required"] = "description is required"
		}
	}
	return errorMessages
}
