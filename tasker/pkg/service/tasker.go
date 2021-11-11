package service

import (
	"anku/popug-jira/tasker/pkg/models"
	"anku/popug-jira/tasker/pkg/pricing"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type UserStorage interface {
	GetUserById(string) (models.User, error)
	StoreUser(models.User) error
	GetAllUsers() ([]models.User, error)
	UpdateUserRole(user models.User) error
}

type TaskStorage interface {
	GetTaskById(string) (models.Task, error)
	GetTasksByUserId(string) ([]models.Task, error)
	StoreTask(models.Task) error
	ChangeAssignee(string, string) error
	Finish(string) error
	GetAllTasks() ([]models.Task, error)
}

type Tasker struct {
	users   UserStorage
	tasks   TaskStorage
	pricing pricing.Pricer
}

func New(u UserStorage, t TaskStorage, p pricing.Pricer) *Tasker {
	return &Tasker{users: u, tasks: t, pricing: p}
}

func (e *Tasker) CreateTask(t models.Task) error {
	fee := e.pricing.Fee(t)
	reward := e.pricing.Reward(t)

	t.Fee = fee
	t.Reward = reward

	id, _ := uuid.NewUUID()
	t.PublicId = id.String()

	return e.tasks.StoreTask(t)
}

func (e *Tasker) ListTasksForUser(userId string) ([]models.Task, error) {
	return e.tasks.GetTasksByUserId(userId)
}

func (e *Tasker) AssignTasks() error {
	users, err := e.users.GetAllUsers()
	if err != nil {
		return err
	}

	var probableAssignees []models.User
	for _, user := range users {
		if user.CanBeAssignee() {
			probableAssignees = append(probableAssignees, user)
		}
	}

	if probableAssignees == nil {
		return nil
	}

	tasks, err := e.tasks.GetAllTasks()
	if err != nil {
		return err
	}

	rand.Seed(time.Now().Unix())

	for _, task := range tasks {
		assignee := probableAssignees[rand.Intn(len(probableAssignees))]

		e.tasks.ChangeAssignee(task.PublicId, assignee.Id)
	}

	return nil
}

func (e *Tasker) FinishTask(id string) error {
	return e.tasks.Finish(id)
}
