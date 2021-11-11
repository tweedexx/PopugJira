package adapters

import (
	"anku/popug-jira/auth/pkg/models"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

func (p *Postgres) GetUserById(u uint64) (models.User, error) {
	usr := models.User{}

	err := p.db.Select(&usr, fmt.Sprintf("select * from users where id = %d", u))

	if err != nil {
		return models.User{}, errors.Wrap(err, fmt.Sprintf("can't get user by id: %v", u))
	}

	return usr, nil
}

func (p *Postgres) GetUserByPublicId(id string) (models.User, error) {
	usr := models.User{}

	err := p.db.Select(&usr, fmt.Sprintf("select * from users where public_id = %v", id))

	if err != nil {
		return models.User{}, errors.Wrap(err, fmt.Sprintf("can't get user by id: %v", id))
	}

	return usr, nil
}

func (p *Postgres) StoreUser(user models.User) error {
	_, err := p.db.Exec("insert into users (role, email, name, password, public_id) values ($1, $2, $3, $4, $5)", user.Role, user.Email, user.Name, user.Credentials, user.PublicId)
	return errors.Wrap(err, fmt.Sprintf("can't store user %+v", user))
}

func (p *Postgres) ChangeUserRole(userId string, newRole string) error {
	_, err := p.db.Exec("update users set role=$1 where public_id=$2", newRole, userId)
	return errors.Wrap(err, fmt.Sprintf("can't change user role %v, %v", userId, newRole))
}
