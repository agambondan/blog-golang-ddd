package application

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
)

type labelApp struct {
	lr repositories.LabelRepository
}

var _ LabelAppInterface = &labelApp{}

type LabelAppInterface interface {
	SaveLabel(*model.Label) (*model.Label, error)
	GetAllLabel() ([]model.Label, error)
	GetLabel(uint64) (*model.Label, error)
	UpdateLabel(*model.Label) (*model.Label, error)
	DeleteLabel(uint64) error
}

func (r *labelApp) SaveLabel(label *model.Label) (*model.Label, error) {
	return r.lr.SaveLabel(label)
}

func (r *labelApp) GetAllLabel() ([]model.Label, error) {
	return r.lr.GetAllLabel()
}

func (r *labelApp) GetLabel(labelId uint64) (*model.Label, error) {
	return r.lr.GetLabel(labelId)
}

func (r *labelApp) UpdateLabel(label *model.Label) (*model.Label, error) {
	return r.lr.UpdateLabel(label)
}

func (r *labelApp) DeleteLabel(labelId uint64) error {
	return r.lr.DeleteLabel(labelId)
}
