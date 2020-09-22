package persistence

import (
	"Repository-Pattern/domain/repositories"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type Repositories struct {
	User         repositories.UserRepository
	Post         repositories.PostRepository
	Role         repositories.RoleRepository
	Label        repositories.LabelRepository
	Category     repositories.CategoryRepository
	PostLabel    repositories.PostLabelRepository
	PostCategory repositories.PostCategoryRepository
	db           *sql.DB
}

func NewRepositories(DbDriver, DbUser, DbPassword, DbPort, DbHost, DbName string) (*Repositories, error) {
	DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
	db, err := sql.Open(DbDriver, DBURL)
	if err != nil {
		return nil, err
	}
	//db.LogMode(true)

	return &Repositories{
		User:     NewUserRepository(db),
		Post:     NewPostRepository(db),
		Role:     NewRoleRepository(db),
		Label:    NewLabelRepository(db),
		Category: NewCategoryRepository(db),
		PostLabel: NewPostLabelRepository(db),
		PostCategory: NewPostCategoryRepository(db),
		db:       db,
	}, nil
}

//closes the  database connection
func (s *Repositories) Close() error {
	return s.db.Close()
}

//This migrate all tables
func (s *Repositories) Seeder() error {
	var err error
	var result sql.Result
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS users (id uuid PRIMARY KEY, first_name VARCHAR(55) not null, last_name VARCHAR(55) not null, email VARCHAR(55) unique not null, " +
		"phone_number VARCHAR(15) not null, username VARCHAR(55) unique not null, password VARCHAR(255) not null, role_id int not null, created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS roles (id serial PRIMARY KEY, name VARCHAR(15) not null unique, " +
		"created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS posts (id serial PRIMARY KEY, title VARCHAR(100) not null unique, description text not null, post_images text, user_id uuid not null, " +
		"created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS labels (id serial PRIMARY KEY, name VARCHAR(35) not null unique, " +
		"created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS categories (id serial PRIMARY KEY, name VARCHAR(35) not null unique, " +
		"created_at timestamp, updated_at timestamp, deleted_at timestamp)")
	log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS post_labels " +
		"(post_id serial not null, label_id int not null)")
	log.Println(result, err)
	result, err = s.db.Exec("CREATE TABLE IF NOT EXISTS post_categories " +
		"(post_id serial not null, category_id int not null)")
	log.Println(result, err)
	return err
}

func (s *Repositories) AddForeignKey() error {
	var err error
	return err
}
