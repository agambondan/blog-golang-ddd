package repositories

import "Repository-Pattern/domain/model"

type CategoryRepository interface {
	SaveCategory(*model.Category) (*model.Category, error)
	GetCategory(uint64) (*model.Category, error)
	GetAllCategory() ([]model.Category, error)
	UpdateCategory(*model.Category) (*model.Category, error)
	DeleteCategory(uint64) error
}