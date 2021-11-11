package transport

import (
	"anku/popug-jira/tasker/pkg/models"
	"anku/popug-jira/tasker/pkg/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Endpoints struct {
	tasker *service.Tasker
}

func New(tasker *service.Tasker) *Endpoints {
	return &Endpoints{tasker: tasker}
}

func (e *Endpoints) CreateTask(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("CreateTask called\n")

	var t models.Task

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		fmt.Printf("can't unmarshal task: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = e.tasker.CreateTask(t)
	if err != nil {
		fmt.Printf("can't create task: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (e *Endpoints) ListTasksForUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	userId := params["id"]

	tasks, err := e.tasker.ListTasksForUser(userId)
	if err != nil {
		fmt.Printf("can't list tasks: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(tasks)
	w.WriteHeader(http.StatusOK)
}

func (e *Endpoints) AssignTasks(w http.ResponseWriter, r *http.Request) {
	err := e.tasker.AssignTasks()
	if err != nil {
		fmt.Printf("can't assign tasks: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (e *Endpoints) FinishTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	taskId := params["id"]

	err := e.tasker.FinishTask(taskId)
	if err != nil {
		fmt.Printf("can't finish task: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
