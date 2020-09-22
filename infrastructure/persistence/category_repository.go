package persistence

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"database/sql"
	"fmt"
	"time"
)

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepo {
	return &CategoryRepo{db}
}

//CategoryRepo implements the repository.CategoryRepository interface
var _ repositories.CategoryRepository = &CategoryRepo{}

func (r *CategoryRepo) SaveCategory(category *model.Category) (*model.Category, error) {
	category.Prepare()
	queryInsert := fmt.Sprintf("INSERT INTO %s (name, created_at, updated_at, deleted_at) "+
		"VALUES ($1, $2, $3, $4) RETURNING id", "categories")
	prepare, err := r.db.Prepare(queryInsert)
	if err != nil {
		return category, err
	}
	err = prepare.QueryRow(&category.Name, &category.CreatedAt, &category.UpdatedAt, nil).Scan(&category.ID)
	if err != nil {
		return category, err
	}
	return category, err
}

func (r *CategoryRepo) GetCategory(id uint64) (*model.Category, error) {
	var category model.Category
	querySelect := fmt.Sprint("SELECT id, name, created_at, updated_at FROM categories WHERE deleted_at IS NULL AND id=$1")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return &category, err
	}
	err = prepare.QueryRow(id).Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		fmt.Println(err)
	}
	return &category, nil
}

func (r *CategoryRepo) GetAllCategory() ([]model.Category, error) {
	var categories []model.Category
	queryGetUsers := fmt.Sprintf("SELECT id, name, created_at, updated_at FROM categories WHERE deleted_at IS NULL")
	rows, err := r.db.Query(queryGetUsers)
	if err != nil {
		return categories, err
	}
	for rows.Next() {
		var category model.Category
		err := rows.Scan(&category.ID, &category.Name, &category.CreatedAt, &category.UpdatedAt)
		if err != nil {
			return categories, err
		}
		categories = append(categories, category)
	}
	defer rows.Close()
	return categories, nil
}

func (r *CategoryRepo) UpdateCategory(category *model.Category) (*model.Category, error) {
	category.UpdatedAt = time.Now()
	queryUpdate := fmt.Sprint("UPDATE categories SET name=$1, updated_at=$2 WHERE id=$3")
	prepare, err := r.db.Prepare(queryUpdate)
	if err != nil {
		return category, err
	}
	_, err = prepare.Exec(category.Name, category.UpdatedAt, category.ID)
	return category, err
}

func (r *CategoryRepo) DeleteCategory(id uint64) error {
	var category model.Category
	category.DeletedAt = time.Now()
	querySoftDelete := fmt.Sprint("UPDATE categories SET deleted_at=$1 WHERE id=$2")
	prepare, err := r.db.Prepare(querySoftDelete)
	if err != nil {
		return err
	}
	_, err = prepare.Exec(category.DeletedAt, id)
	if err != nil {
		return err
	}
	return nil
}