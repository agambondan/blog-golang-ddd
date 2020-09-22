package persistence

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"database/sql"
	"fmt"
)

type PostCategoryRepo struct {
	db *sql.DB
}

func NewPostCategoryRepository(db *sql.DB) *PostCategoryRepo {
	return &PostCategoryRepo{db}
}

//PostCategoryRepo implements the repository.PostCategoryRepository interface
var _ repositories.PostCategoryRepository = &PostCategoryRepo{}

func (r *PostCategoryRepo) SavePostCategory(postCategory *model.PostCategory) (*model.PostCategory, error) {
	queryInsert := fmt.Sprintf("INSERT INTO %s (post_id, label_id) "+
		"VALUES ($1, $2)", "post_categories")
	prepare, err := r.db.Prepare(queryInsert)
	if err != nil {
		return postCategory, err
	}
	_, err = prepare.Exec(&postCategory.PostID, &postCategory.CategoryID)
	if err != nil {
		return postCategory, err
	}
	return postCategory, err
}

func (r *PostCategoryRepo) GetPostCategory(id uint64) (*model.PostCategory, error) {
	var postCategory model.PostCategory
	querySelect := fmt.Sprint("SELECT * FROM post_categories WHERE post_id=$1")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return &postCategory, err
	}
	err = prepare.QueryRow(id).Scan(&postCategory.PostID, &postCategory.CategoryID)
	if err != nil {
		fmt.Println(err)
	}
	return &postCategory, nil
}

func (r *PostCategoryRepo) GetAllPostCategory() ([]model.PostCategory, error) {
	var postCategorys []model.PostCategory
	queryGetUsers := fmt.Sprintf("SELECT * FROM post_categories")
	prepare, err := r.db.Prepare(queryGetUsers)
	if err != nil {
		return postCategorys, err
	}
	rows, err := prepare.Query()
	if err != nil {
		return postCategorys, err
	}
	for rows.Next() {
		var postCategory model.PostCategory
		err := rows.Scan(&postCategory.PostID, &postCategory.CategoryID)
		if err != nil {
			return postCategorys, err
		}
		postCategorys = append(postCategorys, postCategory)
	}
	defer rows.Close()
	return postCategorys, nil
}

func (r *PostCategoryRepo) UpdatePostCategory(postCategory *model.PostCategory) (*model.PostCategory, error) {
	queryUpdate := fmt.Sprint("UPDATE post_categories SET post_id=$1, label_id=$2 WHERE post_id=$3")
	prepare, err := r.db.Prepare(queryUpdate)
	if err != nil {
		return postCategory, err
	}
	_, err = prepare.Exec(postCategory.PostID, postCategory.CategoryID, postCategory.PostID)
	return postCategory, err
}

func (r *PostCategoryRepo) DeletePostCategory(id uint64) error {
	querySoftDelete := fmt.Sprint("DELETE FROM post_categories WHERE post_id=$1")
	prepare, err := r.db.Prepare(querySoftDelete)
	if err != nil {
		return err
	}
	_, err = prepare.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
