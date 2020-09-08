package repositories

import "Repository-Pattern/domain/model"

type UserRepository interface {
	SaveUser(*model.User) (*model.User, map[string]string)
	GetUser(uint64) (*model.User, error)
	GetUsers() ([]model.User, error)
	GetUserByEmailAndPassword(*model.User) (*model.User, map[string]string)
}
