package repositories

import (
	"Repository-Pattern/domain/model"
	"github.com/google/uuid"
)

type UserRepository interface {
	SaveUser(*model.User) (*model.User, error)
	GetUser(uuid uuid.UUID) (*model.User, error)
	GetUsers() ([]model.User, error)
	GetUserByEmailAndPassword(*model.User) (*model.User, error)
}
