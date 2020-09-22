package persistence

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type PostRepo struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepo {
	return &PostRepo{db}
}

//PostRepo implements the repository.PostRepository interface
var _ repositories.PostRepository = &PostRepo{}

func (r *PostRepo) SavePost(post *model.Post) (*model.Post, error) {
	post.Prepare()
	queryInsert := fmt.Sprint("INSERT INTO posts (title, description, post_images, user_id, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id")
	prepare, err := r.db.Prepare(queryInsert)
	if err != nil {
		return post, err
	}
	err = prepare.QueryRow(&post.Title, &post.Description, &post.PostImage, &post.UserUUID, &post.CreatedAt, &post.UpdatedAt, nil).Scan(&post.ID)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (r *PostRepo) GetPost(id uint64) (*model.Post, error) {
	var post model.Post
	querySelect := fmt.Sprint("SELECT id, title, description, post_images, user_id, created_at, updated_at FROM posts WHERE deleted_at is NULL AND id=$1")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return &post, err
	}
	err = prepare.QueryRow(id).Scan(&post.ID, &post.Title, &post.Description, &post.PostImage, &post.UserUUID, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return &post, err
	}
	return &post, nil
}

func (r *PostRepo) GetPostByIdUser(userUuid uuid.UUID) ([]model.Post, error) {
	var posts []model.Post
	querySelect := fmt.Sprint("SELECT id, title, description, post_images, user_id, created_at, updated_at FROM posts WHERE deleted_at is NULL AND user_id=$1")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return posts, err
	}
	rows, err := prepare.Query(userUuid)
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.PostImage, &post.UserUUID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepo) GetAllPost() ([]model.Post, error) {
	var posts []model.Post
	querySelect := fmt.Sprint("SELECT id, title, description, post_images, user_id, created_at, updated_at FROM posts WHERE deleted_at is NULL")
	prepare, err := r.db.Prepare(querySelect)
	if err != nil {
		return posts, err
	}
	rows, err := prepare.Query()
	if err != nil {
		return posts, err
	}
	for rows.Next() {
		var post model.Post
		err := rows.Scan(&post.ID, &post.Title, &post.Description, &post.PostImage, &post.UserUUID, &post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return posts, err
		}
		posts = append(posts, post)
	}
	return posts, err
}

func (r *PostRepo) UpdatePost(post *model.Post) (*model.Post, error) {
	post.UpdatedAt = time.Now()
	queryInsert := fmt.Sprint("UPDATE posts SET title=$1, description=$2, post_images=$3, updated_at=$4 WHERE id=$5")
	prepare, err := r.db.Prepare(queryInsert)
	if err != nil {
		return post, err
	}
	_, err = prepare.Exec(&post.Title, &post.Description, &post.PostImage, &post.CreatedAt, &post.UpdatedAt, nil)
	if err != nil {
		return post, err
	}
	return post, nil
}

func (r *PostRepo) DeletePost(id uint64) error {
	var post model.Post
	post.DeletedAt = time.Now()
	querySoftDelete := fmt.Sprint("UPDATE posts SET deleted_at=$1 WHERE id=$2")
	prepare, err := r.db.Prepare(querySoftDelete)
	if err != nil {
		return err
	}
	_, err = prepare.Exec(post.DeletedAt, id)
	if err != nil {
		return err
	}
	return nil
}
