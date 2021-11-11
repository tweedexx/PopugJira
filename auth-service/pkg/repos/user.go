package repos

import (
	"anku/popug-jira/auth/pkg/adapters"
	"anku/popug-jira/auth/pkg/models"
	"anku/popug-jira/auth/pkg/service"
)

type UserRepo struct {
	userStorage service.UserStorage
	kafka       *adapters.Kafka
}

func New(us service.UserStorage, kafka *adapters.Kafka) *UserRepo {
	return &UserRepo{
		userStorage: us,
		kafka:       kafka,
	}
}

func (u *UserRepo) GetUserById(id uint64) (models.User, error) {
	return u.userStorage.GetUserById(id)
}

func (u *UserRepo) StoreUser(user models.User) error {
	err := u.userStorage.StoreUser(user)
	if err == nil {
		return err
	}

	u.kafka.Send(adapters.UserCreated, map[string]interface{}{
		"id":   user.PublicId,
		"role": user.Role,
	})

	return nil
}

func (u *UserRepo) ChangeUserRole(userId string, newRole string) error {
	err := u.userStorage.ChangeUserRole(userId, newRole)
	if err != nil {
		return err
	}

	u.kafka.Send(adapters.UserRoleChanged, map[string]interface{}{
		"id":   userId,
		"role": newRole,
	})

	return nil
}

func (u *UserRepo) GetUserByPublicId(id string) (models.User, error) {
	return u.userStorage.GetUserByPublicId(id)
}
