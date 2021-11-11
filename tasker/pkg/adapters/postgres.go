package adapters

import (
	"anku/popug-jira/tasker/pkg/models"
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

func (p *Postgres) GetUserById(u string) (models.User, error) {
	usr := models.User{}

	err := p.db.Select(&usr, fmt.Sprintf("select * from users where id = %d", u))

	if err != nil {
		return models.User{}, errors.Wrap(err, fmt.Sprintf("can't get user by id: %v", u))
	}

	return usr, nil
}

func (p *Postgres) StoreUser(user models.User) error {
	_, err := p.db.Exec("insert into users (id, role) values ($1, $2)", user.Id, user.Role)
	return errors.Wrap(err, fmt.Sprintf("can't store user %+v", user))
}

func (p *Postgres) GetAllUsers() ([]models.User, error) {
	var usrs []models.User

	err := p.db.Select(&usrs, "select * from users")

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't get all users"))
	}

	return usrs, nil
}

func (p *Postgres) UpdateUserRole(user models.User) error {
	_, err := p.db.Exec("update users set role=$1 where id=$2", user.Role, user.Id)
	return errors.Wrap(err, fmt.Sprintf("can't update user %+v", user))
}

func (p *Postgres) GetTaskById(u string) (models.Task, error) {
	tsk := models.Task{}

	err := p.db.Select(&tsk, fmt.Sprintf("select * from tasks where public_id = %v", u))

	if err != nil {
		return models.Task{}, errors.Wrap(err, fmt.Sprintf("can't get task by id: %v", u))
	}

	return tsk, nil
}

func (p *Postgres) GetTasksByUserId(u string) ([]models.Task, error) {
	var tsks []models.Task

	err := p.db.Select(&tsks, fmt.Sprintf("select * from tasks where assignee_id = %d", u))

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't get task by user id: %v", u))
	}

	return tsks, nil
}

func (p *Postgres) StoreTask(task models.Task) error {
	_, err := p.db.Exec("insert into tasks (assignee_id, status, fee, description, reward, public_id) values ($1, $2, $3, $4, $5, $6)", task.AssigneeId, task.Status, task.Fee, task.Description, task.Reward, task.PublicId)
	return errors.Wrap(err, fmt.Sprintf("can't store task %+v", task))
}

func (p *Postgres) GetAllTasks() ([]models.Task, error) {
	var tsks []models.Task

	err := p.db.Select(&tsks, "select * from tasks")

	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("can't get all tasks"))
	}

	return tsks, nil
}

func (p *Postgres) ChangeAssignee(taskId string, assigneeId string) error {
	_, err := p.db.Exec("update tasks set assignee_id=$1 where public_id=$2", assigneeId, taskId)
	return errors.Wrap(err, fmt.Sprintf("can't change assignee of task %+v", taskId))
}

func (p *Postgres) Finish(taskId string) error {
	_, err := p.db.Exec("update tasks set status=$1 where public_id=$2", models.Done, taskId)
	return errors.Wrap(err, fmt.Sprintf("can't finish task %+v", taskId))
}
