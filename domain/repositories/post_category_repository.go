package repositories

import "Repository-Pattern/domain/model"

type PostCategoryRepository interface {
	SavePostCategory(*model.PostCategory) (*model.PostCategory, error)
	GetPostCategory(uint64) (*model.PostCategory, error)
	GetAllPostCategory() ([]model.PostCategory, error)
	UpdatePostCategory(*model.PostCategory) (*model.PostCategory, error)
	DeletePostCategory(uint64) error
}