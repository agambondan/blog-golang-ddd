package repositories

import "Repository-Pattern/domain/model"

type PostLabelRepository interface {
	SavePostLabel(*model.PostLabel) (*model.PostLabel, error)
	GetPostLabel(uint64) (*model.PostLabel, error)
	GetAllPostLabel() ([]model.PostLabel, error)
	UpdatePostLabel(*model.PostLabel) (*model.PostLabel, error)
	DeletePostLabel(uint64) error
}