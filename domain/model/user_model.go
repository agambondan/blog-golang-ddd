package model

import (
	"Repository-Pattern/infrastructure/security"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"html"
	"strings"
	"time"
)

type User struct {
	UUID        uuid.UUID  `gorm:"type:uuid;primary_key;" json:"id"`
	FirstName   string     `gorm:"size:100;not null;" json:"first_name"`
	LastName    string     `gorm:"size:100;not null;" json:"last_name"`
	Email       string     `gorm:"size:100;not null;unique" json:"email"`
	Username    string     `gorm:"size:100;not null;" json:"username,omitempty"`
	Password    string     `gorm:"size:100;not null;" json:"password"`
	PhoneNumber string     `gorm:"size:100;not null;" json:"phone_number,omitempty"`
	Posts       []Post     `json:"posts,omitempty"`
	Role        Role       `json:"role,omitempty"`
	RoleId      int        `gorm:"not null;" json:"role_id,omitempty"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at,omitempty"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at,omitempty"`
	DeletedAt   *time.Time `gorm:"size:100;not null;" json:"deleted_at,omitempty"`
}

type PublicUser struct {
	UUID      uuid.UUID `sql:"primary_key" json:"id,omitempty"`
	FirstName string    `gorm:"size:100;not null;" json:"first_name"`
	LastName  string    `gorm:"size:100;not null;" json:"last_name"`
}

// BeforeCreate will set a UUID rather than numeric ID.
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	uuidV4 := uuid.NewV4()
	return scope.SetColumn("ID", uuidV4)
}

//BeforeSave is a gorm hook
func (u *User) BeforeSave() error {
	hashPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashPassword)
	return nil
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
		UUID:      u.UUID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func (u *User) Prepare() {
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
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
