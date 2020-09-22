package repositories

import "Repository-Pattern/domain/model"

type LabelRepository interface {
	SaveLabel(*model.Label) (*model.Label, error)
	GetLabel(uint64) (*model.Label, error)
	GetAllLabel() ([]model.Label, error)
	UpdateLabel(*model.Label) (*model.Label, error)
	DeleteLabel(uint64) error
}