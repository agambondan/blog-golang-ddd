package application

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
)

type roleApp struct {
	rr repositories.RoleRepository
}

var _ RoleAppInterface = &roleApp{}

type RoleAppInterface interface {
	SaveRole(*model.Role) (*model.Role, map[string]string)
	GetAllRole() ([]model.Role, error)
	GetRole(uint64) (*model.Role, error)
	UpdateRole(*model.Role) (*model.Role, map[string]string)
	DeleteRole(uint64) error
}

func (r *roleApp) SaveRole(role *model.Role) (*model.Role, map[string]string) {
	return r.rr.SaveRole(role)
}

func (r *roleApp) GetAllRole() ([]model.Role, error) {
	return r.rr.GetAllRole()
}

func (r *roleApp) GetRole(roleId uint64) (*model.Role, error) {
	return r.rr.GetRole(roleId)
}

func (r *roleApp) UpdateRole(role *model.Role) (*model.Role, map[string]string) {
	return r.rr.UpdateRole(role)
}

func (r *roleApp) DeleteRole(roleId uint64) error {
	return r.rr.DeleteRole(roleId)
}
