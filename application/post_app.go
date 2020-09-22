package application

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"github.com/google/uuid"
)

type postApp struct {
	pr repositories.PostRepository
}

var _ PostAppInterface = &postApp{}

type PostAppInterface interface {
	SavePost(*model.Post) (*model.Post, error)
	GetAllPost() ([]model.Post, error)
	GetPostByIdUser(userUuid uuid.UUID) ([]model.Post, error)
	GetPost(uint64) (*model.Post, error)
	UpdatePost(*model.Post) (*model.Post, error)
	DeletePost(uint64) error
}

func (p *postApp) SavePost(post *model.Post) (*model.Post, error) {
	return p.pr.SavePost(post)
}

func (p *postApp) GetAllPost() ([]model.Post, error) {
	return p.pr.GetAllPost()
}

func (p *postApp) GetPostByIdUser(userUuid uuid.UUID) ([]model.Post, error) {
	return p.pr.GetPostByIdUser(userUuid)
}

func (p *postApp) GetPost(postId uint64) (*model.Post, error) {
	return p.pr.GetPost(postId)
}

func (p *postApp) UpdatePost(post *model.Post) (*model.Post, error) {
	return p.pr.UpdatePost(post)
}

func (p *postApp) DeletePost(postId uint64) error {
	return p.pr.DeletePost(postId)
}
