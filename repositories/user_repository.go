package repositories

import (
	"Repository-Pattern/models"
	"fmt"
	"github.com/google/uuid"
)

type UserRepository interface {
	List() []models.User
	Get(ID uuid.UUID) models.User
	Create() models.User
	Update(ID string) models.User
	Delete(ID string) models.User
}

var (
	user  models.User
	users []models.User
)

func ListUser() ([]models.User, error) {
	rows, err := DB.conn.Query("Select * from users")
	if err != nil {
		return users, err
	}
	for rows.Next() {
		err = rows.Scan(user.ID, user.CreatedAt, user.UpdatedAt, user.DeletedAt, user.FullName, user.PhoneNumber, user.Username,
			user.Password, user.Email, user.RoleId)
		if err != nil {
			return users, err
		}
		var role models.Role
		err = DB.conn.QueryRow("SELECT id, created_at, updated_at, name FROM role WHERE id=$1", &user.RoleId).
			Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt, &role.Name)
		if err != nil {
			fmt.Println(err)
		}
		user.Role = role
		rowsPost, err := DB.conn.Query("SELECT id, created_at, updated_at, title, content, author_id FROM post WHERE author_id=$1", &user.ID)
		if err != nil {
			fmt.Println(err)
		}
		for rowsPost.Next() {
			var post models.Post
			err := rowsPost.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt, &post.Title, &post.Content, &post.AuthorID)
			if err != nil {
				fmt.Println(err.Error())
			}
			user.Posts = append(user.Posts, post)
		}
		users = append(users, user)
	}
	defer rows.Close()
	return users, err
}

func GetUserById() {
}

func CreateUser(data interface{}, query string) {
}

func UpdateUserById(ID string, data interface{}, query string) {
}

func DeleteUserById(ID string, query string) {
}
