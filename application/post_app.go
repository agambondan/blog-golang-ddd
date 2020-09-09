package application

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
)

type PostApp struct {
	fr repositories.PostRepository
}

var _ PostAppInterface = &PostApp{}

type PostAppInterface interface {
	SavePost(*model.Post) (*model.Post, map[string]string)
	GetAllPost() ([]model.Post, error)
	GetPost(uint64) (*model.Post, error)
	UpdatePost(*model.Post) (*model.Post, map[string]string)
	DeletePost(uint64) error
}

func (f *PostApp) SavePost(Post *model.Post) (*model.Post, map[string]string) {
	return f.fr.SavePost(Post)
}

func (f *PostApp) GetAllPost() ([]model.Post, error) {
	return f.fr.GetAllPost()
}

func (f *PostApp) GetPost(PostId uint64) (*model.Post, error) {
	return f.fr.GetPost(PostId)
}

func (f *PostApp) UpdatePost(Post *model.Post) (*model.Post, map[string]string) {
	return f.fr.UpdatePost(Post)
}

func (f *PostApp) DeletePost(PostId uint64) error {
	return f.fr.DeletePost(PostId)
}
