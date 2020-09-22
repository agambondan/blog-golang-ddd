package application

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
)

type PostCategoryApp struct {
	lr repositories.PostCategoryRepository
}

var _ PostCategoryAppInterface = &PostCategoryApp{}

type PostCategoryAppInterface interface {
	SavePostCategory(*model.PostCategory) (*model.PostCategory, error)
	GetAllPostCategory() ([]model.PostCategory, error)
	GetPostCategory(uint64) (*model.PostCategory, error)
	UpdatePostCategory(*model.PostCategory) (*model.PostCategory, error)
	DeletePostCategory(uint64) error
}

func (r *PostCategoryApp) SavePostCategory(postCategory *model.PostCategory) (*model.PostCategory, error) {
	return r.lr.SavePostCategory(postCategory)
}

func (r *PostCategoryApp) GetAllPostCategory() ([]model.PostCategory, error) {
	return r.lr.GetAllPostCategory()
}

func (r *PostCategoryApp) GetPostCategory(postCategoryId uint64) (*model.PostCategory, error) {
	return r.lr.GetPostCategory(postCategoryId)
}

func (r *PostCategoryApp) UpdatePostCategory(postCategory *model.PostCategory) (*model.PostCategory, error) {
	return r.lr.UpdatePostCategory(postCategory)
}

func (r *PostCategoryApp) DeletePostCategory(postCategoryId uint64) error {
	return r.lr.DeletePostCategory(postCategoryId)
}
