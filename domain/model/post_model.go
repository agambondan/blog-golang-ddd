package model

import (
	"github.com/google/uuid"
	"html"
	"strings"
	"time"
)

type Post struct {
	ID          uint64     `json:"id"`
	UserUUID    uuid.UUID  `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	PostImage   string     `json:"post_image"`
	Author      *User      `json:"author,omitempty"`
	Labels      []Label    `json:"labels,omitempty"`
	Categories  []Category `json:"categories,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   time.Time  `sql:"index" json:"deleted_at,omitempty"`
}

type PublicPost struct {
	ID          uint64    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	PostImage   string    `json:"post_image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Posts []Post

//So that we dont expose the user's email address and password to the world
func (posts Posts) PublicPosts() []interface{} {
	result := make([]interface{}, len(posts))
	for index, post := range posts {
		result[index] = post.PublicPost()
	}
	return result
}

func (p *Post) PublicPost() interface{} {
	return &PublicPost{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		PostImage:   p.PostImage,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (p *Post) Prepare() {
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Description = html.EscapeString(strings.TrimSpace(p.Description))
	p.PostImage = html.EscapeString(strings.TrimSpace(p.PostImage))
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
