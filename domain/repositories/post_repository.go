package repositories

import (
	"Repository-Pattern/domain/model"
	"github.com/google/uuid"
)

type PostRepository interface {
	SavePost(*model.Post) (*model.Post, error)
	GetPost(uint64) (*model.Post, error)
	GetPostByIdUser(userUuid uuid.UUID) ([]model.Post, error)
	GetAllPost() ([]model.Post, error)
	UpdatePost(*model.Post) (*model.Post, error)
	DeletePost(uint64) error
}
