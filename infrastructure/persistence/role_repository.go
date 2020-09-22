package persistence

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"database/sql"
	"fmt"
	"time"
)

type RoleRepo struct {
	db *sql.DB
}

func NewRoleRepository(db *sql.DB) *RoleRepo {
	return &RoleRepo{db}
}

//RoleRepo implements the repository.RoleRepository interface
var _ repositories.RoleRepository = &RoleRepo{}

func (r *RoleRepo) SaveRole(role *model.Role) (*model.Role, error) {
	role.Prepare()
	queryInsert := fmt.Sprintf("INSERT INTO %s (name, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4) RETURNING id", "roles")
	prepare, err := r.db.Prepare(queryInsert)
	if err != nil {
		return role, err
	}
	err = prepare.QueryRow(&role.Name, &role.CreatedAt, &role.UpdatedAt, nil).Scan(&role.ID)
	if err != nil {
		return role, err
	}
	return role, err
}

func (r *RoleRepo) GetRole(id uint64) (*model.Role, error) {
	var role model.Role
	querySelect := fmt.Sprint("SELECT id, name, created_at, updated_at FROM roles WHERE deleted_at IS NULL AND id=$1")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return &role, err
	}
	err = prepare.QueryRow(id).Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return &role, nil
}

func (r *RoleRepo) GetAllRole() ([]model.Role, error) {
	var roles []model.Role
	queryGetUsers := fmt.Sprintf("SELECT id, name, created_at, updated_at FROM roles WHERE deleted_at IS NULL")
	rows, err := r.db.Query(queryGetUsers)
	if err != nil {
		return roles, err
	}
	for rows.Next() {
		var role model.Role
		err := rows.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return roles, err
		}
		roles = append(roles, role)
	}
	defer rows.Close()
	return roles, nil
}

func (r *RoleRepo) UpdateRole(role *model.Role) (*model.Role, error) {
	role.UpdatedAt = time.Now()
	queryUpdate := fmt.Sprint("UPDATE roles SET name=$1, updated_at=$2 WHERE id=$3")
	prepare, err := r.db.Prepare(queryUpdate)
	if err != nil {
		return role, err
	}
	_, err = prepare.Exec(role.Name, role.UpdatedAt, role.ID)
	return role, err
}

func (r *RoleRepo) DeleteRole(id uint64) error {
	var role model.Role
	role.DeletedAt = time.Now()
	querySoftDelete := fmt.Sprint("UPDATE roles SET deleted_at=$1 WHERE id=$2")
	prepare, err := r.db.Prepare(querySoftDelete)
	if err != nil {
		return err
	}
	_, err = prepare.Exec(role.DeletedAt, id)
	if err != nil {
		return err
	}
	return nil
}
