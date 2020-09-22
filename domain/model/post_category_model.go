package model

type PostCategory struct {
	PostID     uint64 `json:"post_id,omitempty"`
	CategoryID uint64 `json:"category_id,omitempty"`
}
