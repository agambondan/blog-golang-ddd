package persistence

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"errors"
	"github.com/jinzhu/gorm"
	"strings"
)

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db}
}

//RoleRepo implements the repository.RoleRepository interface
var _ repositories.RoleRepository = &RoleRepo{}

func (r *RoleRepo) SaveRole(role *model.Role) (*model.Role, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Create(&role).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "Role title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return role, dbErr
}

func (r *RoleRepo) GetRole(id uint64) (*model.Role, error) {
	var Role model.Role
	err := r.db.Debug().Where("id = ?", id).Take(&Role).Error
	if err != nil {
		return nil, errors.New("database error, please try again")
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("role not found")
	}
	return &Role, nil
}

func (r *RoleRepo) GetAllRole() ([]model.Role, error) {
	var Roles []model.Role
	err := r.db.Debug().Limit(100).Order("created_at desc").Find(&Roles).Error
	if err != nil {
		return nil, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return nil, errors.New("role not found")
	}
	return Roles, nil
}

func (r *RoleRepo) UpdateRole(role *model.Role) (*model.Role, map[string]string) {
	dbErr := map[string]string{}
	err := r.db.Debug().Save(&role).Error
	if err != nil {
		//since our title is unique
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "Duplicate") {
			dbErr["unique_title"] = "title already taken"
			return nil, dbErr
		}
		//any other db error
		dbErr["db_error"] = "database error"
		return nil, dbErr
	}
	return role, nil
}

func (r *RoleRepo) DeleteRole(id uint64) error {
	var Role model.Role
	err := r.db.Debug().Where("id = ?", id).Delete(&Role).Error
	if err != nil {
		return errors.New("database error, please try again")
	}
	return nil
}
