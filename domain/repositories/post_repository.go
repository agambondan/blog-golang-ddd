package repositories

import "Repository-Pattern/domain/model"

type PostRepository interface {
	SavePost(*model.Post) (*model.Post, map[string]string)
	GetPost(uint64) (*model.Post, error)
	GetAllPost() ([]model.Post, error)
	UpdatePost(*model.Post) (*model.Post, map[string]string)
	DeletePost(uint64) error
}
