package application

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
)

type PostLabelApp struct {
	lr repositories.PostLabelRepository
}

var _ PostLabelAppInterface = &PostLabelApp{}

type PostLabelAppInterface interface {
	SavePostLabel(*model.PostLabel) (*model.PostLabel, error)
	GetAllPostLabel() ([]model.PostLabel, error)
	GetPostLabel(uint64) (*model.PostLabel, error)
	UpdatePostLabel(*model.PostLabel) (*model.PostLabel, error)
	DeletePostLabel(uint64) error
}

func (r *PostLabelApp) SavePostLabel(postLabel *model.PostLabel) (*model.PostLabel, error) {
	return r.lr.SavePostLabel(postLabel)
}

func (r *PostLabelApp) GetAllPostLabel() ([]model.PostLabel, error) {
	return r.lr.GetAllPostLabel()
}

func (r *PostLabelApp) GetPostLabel(postLabelId uint64) (*model.PostLabel, error) {
	return r.lr.GetPostLabel(postLabelId)
}

func (r *PostLabelApp) UpdatePostLabel(postLabel *model.PostLabel) (*model.PostLabel, error) {
	return r.lr.UpdatePostLabel(postLabel)
}

func (r *PostLabelApp) DeletePostLabel(postLabelId uint64) error {
	return r.lr.DeletePostLabel(postLabelId)
}
