package persistence

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepo {
	return &UserRepo{db}
}

//UserRepo implements the repository.UserRepository interface
var _ repositories.UserRepository = &UserRepo{}

func (r *UserRepo) SaveUser(user *model.User) (*model.User, error) {
	user.Prepare()
	queryInsert := fmt.Sprintf("INSERT INTO %s (id, first_name, last_name, email, phone_number, username, password, role_id, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)", "users")
	stmt, err := r.db.Prepare(queryInsert)
	if err != nil {
		return user, err
	}
	_, err = stmt.Exec(user.UUID, user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.Username, user.Password, user.RoleId, user.CreatedAt, user.UpdatedAt, nil)
	if err != nil {
		return user, err
	}
	return user, err
}

func (r *UserRepo) GetUser(id uuid.UUID) (*model.User, error) {
	var user model.User
	querySelect := fmt.Sprint("SELECT id, first_name, last_name, email, phone_number, username, password, role_id, created_at, updated_at FROM users WHERE id=$1 AND deleted_at IS NULL")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return &user, err
	}
	err = prepare.QueryRow(id).Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.RoleId, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return &user, err
	}
	return &user, nil
}

func (r *UserRepo) GetUsers() ([]model.User, error) {
	var users []model.User
	queryGetUsers := fmt.Sprintf("SELECT id, first_name, last_name, email, phone_number, username, password, role_id, created_at, updated_at FROM users WHERE deleted_at IS NULL")
	rows, err := r.db.Query(queryGetUsers)
	if err != nil {
		return users, err
	}
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.RoleId, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	defer rows.Close()
	return users, nil
}

func (r *UserRepo) GetUserByEmailAndPassword(u *model.User) (*model.User, error) {
	var user model.User
	queryLogin := fmt.Sprint("SELECT id, first_name, last_name, email, phone_number, username, password, role_id, created_at, updated_at "+
		"FROM users WHERE email=$1 AND password=$2 AND deleted_at IS NULL")
	err := r.db.QueryRow(queryLogin, u.Email, u.Password).Scan(&user.UUID, &user.FirstName, &user.LastName, &user.Email, &user.PhoneNumber, &user.Username, &user.Password, &user.RoleId, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return &user, err
	}
	return &user, nil
}