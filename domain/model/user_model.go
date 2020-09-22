package model

import (
	"github.com/badoux/checkmail"
	"github.com/google/uuid"
	"html"
	"strings"
	"time"
)

type User struct {
	UUID        uuid.UUID `sql:"primary_key" json:"id"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number,omitempty"`
	Username    string    `json:"username,omitempty"`
	Password    string    `json:"password,omitempty"`
	RoleId      uint64    `json:"role_id,omitempty"`
	Posts       []Post    `json:"posts,omitempty"`
	Role        *Role     `json:"role,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	DeletedAt   time.Time `sql:"index" json:"deleted_at,omitempty"`
}

type PublicUser struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Users []User

//So that we dont expose the user's email address and password to the world
func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublicUser()
	}
	return result
}

//So that we dont expose the user's email address and password to the world
func (u *User) PublicUser() interface{} {
	return &PublicUser{
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (u *User) Prepare() {
	u.UUID, _ = uuid.NewUUID()
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.PhoneNumber = html.EscapeString(strings.TrimSpace(u.PhoneNumber))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if u.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "email email"
			}
		}

	case "login":
		if u.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if u.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	case "forgotpassword":
		if u.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	default:
		if u.FirstName == "" {
			errorMessages["firstname_required"] = "first name is required"
		}
		if u.LastName == "" {
			errorMessages["lastname_required"] = "last name is required"
		}
		if u.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if u.Password != "" && len(u.Password) < 6 {
			errorMessages["invalid_password"] = "password should be at least 6 characters"
		}
		if u.PhoneNumber == "" {
			errorMessages["phone_number_required"] = "phone number is required"
		}
		if u.RoleId == 0 {
			errorMessages["role_id_required"] = "role id is required"
		}
		if u.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	}
	return errorMessages
}
