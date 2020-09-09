package application

import (
	"Repository-Pattern/domain/model"
	"Repository-Pattern/domain/repositories"
	"github.com/google/uuid"
)

type userApp struct {
	us repositories.UserRepository
}

//UserApp implements the UserAppInterface
var _ UserAppInterface = &userApp{}

type UserAppInterface interface {
	SaveUser(*model.User) (*model.User, map[string]string)
	GetUsers() ([]model.User, error)
	GetUser(uuid uuid.UUID) (*model.User, error)
	GetUserByEmailAndPassword(*model.User) (*model.User, map[string]string)
}

func (u *userApp) SaveUser(user *model.User) (*model.User, map[string]string) {
	return u.us.SaveUser(user)
}

func (u *userApp) GetUser(userId uuid.UUID) (*model.User, error) {
	return u.us.GetUser(userId)
}

func (u *userApp) GetUsers() ([]model.User, error) {
	return u.us.GetUsers()
}

func (u *userApp) GetUserByEmailAndPassword(user *model.User) (*model.User, map[string]string) {
	return u.us.GetUserByEmailAndPassword(user)
}
