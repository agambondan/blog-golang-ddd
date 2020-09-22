package model

type PostLabel struct {
	PostID  uint64 `json:"post_id,omitempty"`
	LabelID uint64 `json:"label_id,omitempty"`
}
