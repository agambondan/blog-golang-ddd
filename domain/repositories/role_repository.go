package repositories

import "Repository-Pattern/domain/model"

type RoleRepository interface {
	SaveRole(*model.Role) (*model.Role, error)
	GetRole(uint64) (*model.Role, error)
	GetAllRole() ([]model.Role, error)
	UpdateRole(*model.Role) (*model.Role, error)
	DeleteRole(uint64) error
}