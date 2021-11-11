package service

import (
	"anku/popug-jira/auth/pkg/models"

	"github.com/google/uuid"
)

type UserStorage interface {
	GetUserById(uint64) (models.User, error)
	GetUserByPublicId(string) (models.User, error)
	StoreUser(models.User) error
	ChangeUserRole(string, string) error
}

type Auth struct {
	users UserStorage
}

func New(u UserStorage) *Auth {
	return &Auth{users: u}
}

func (e *Auth) CreateUser(t models.User) error {
	id, _ := uuid.NewUUID()
	t.PublicId = id.String()
	return e.users.StoreUser(t)
}

func (e *Auth) ChangeRole(userId string, newRole string) error {
	return e.users.ChangeUserRole(userId, newRole)
}

func (e *Auth) Login() error {
	return nil
}
