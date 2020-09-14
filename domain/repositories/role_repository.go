package repositories

import "Repository-Pattern/domain/model"

type RoleRepository interface {
	SaveRole(*model.Role) (*model.Role, map[string]string)
	GetRole(uint64) (*model.Role, error)
	GetAllRole() ([]model.Role, error)
	UpdateRole(*model.Role) (*model.Role, map[string]string)
	DeleteRole(uint64) error
}