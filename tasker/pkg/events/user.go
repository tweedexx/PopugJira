package events

import (
	"anku/popug-jira/tasker/pkg/adapters"
	"anku/popug-jira/tasker/pkg/models"
	"anku/popug-jira/tasker/pkg/service"
	"fmt"
	"math/rand"
	"time"
)

type UserHandler struct {
	in    chan adapters.UserEvent
	users service.UserStorage
	tasks service.TaskStorage
}

func New(in chan adapters.UserEvent, users service.UserStorage, tasks service.TaskStorage) *UserHandler {
	return &UserHandler{
		in:    in,
		users: users,
		tasks: tasks,
	}
}

func (u *UserHandler) Route() {
	for {
		for usr := range u.in {
			if usr.EventType == adapters.UserCreated {
				err := u.HandleUserCreated(usr.User)
				fmt.Printf("can't create user %+v\n", err)
				continue
			}

			if usr.EventType == adapters.UserRoleChanged {
				err := u.HandleUserRoleChanged(usr.User)
				fmt.Printf("can't change user's role %+v\n", err)
				continue
			}

			fmt.Printf("unknown user type %v: %+v\n", usr.EventType, usr.User)
		}
	}
}

func (u *UserHandler) HandleUserCreated(user models.User) error {
	return u.users.StoreUser(user)
}

func (u *UserHandler) HandleUserRoleChanged(user models.User) error {
	oldUser, err := u.users.GetUserById(user.Id)
	if err != nil {
		return err
	}

	if oldUser.Id != "" && oldUser.CanBeAssignee() && !user.CanBeAssignee() {
		tasks, err := u.tasks.GetTasksByUserId(user.Id)
		if err != nil {
			return err
		}

		users, err := u.users.GetAllUsers()
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

		rand.Seed(time.Now().Unix())

		for _, task := range tasks {
			assignee := probableAssignees[rand.Intn(len(probableAssignees))]

			u.tasks.ChangeAssignee(task.PublicId, assignee.Id)
		}
	}

	return u.users.UpdateUserRole(user)
}
