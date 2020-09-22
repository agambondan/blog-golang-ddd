package application

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
)

type categoryApp struct {
	cr repositories.CategoryRepository
}

var _ CategoryAppInterface = &categoryApp{}

type CategoryAppInterface interface {
	SaveCategory(*model.Category) (*model.Category, error)
	GetAllCategory() ([]model.Category, error)
	GetCategory(uint64) (*model.Category, error)
	UpdateCategory(*model.Category) (*model.Category, error)
	DeleteCategory(uint64) error
}

func (r *categoryApp) SaveCategory(category *model.Category) (*model.Category, error) {
	return r.cr.SaveCategory(category)
}

func (r *categoryApp) GetAllCategory() ([]model.Category, error) {
	return r.cr.GetAllCategory()
}

func (r *categoryApp) GetCategory(categoryId uint64) (*model.Category, error) {
	return r.cr.GetCategory(categoryId)
}

func (r *categoryApp) UpdateCategory(category *model.Category) (*model.Category, error) {
	return r.cr.UpdateCategory(category)
}

func (r *categoryApp) DeleteCategory(categoryId uint64) error {
	return r.cr.DeleteCategory(categoryId)
}
