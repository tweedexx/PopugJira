package transport

import (
	"anku/popug-jira/auth/pkg/models"
	"anku/popug-jira/auth/pkg/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Endpoints struct {
	auth *service.Auth
}

func New(auth *service.Auth) *Endpoints {
	return &Endpoints{auth: auth}
}

func (e *Endpoints) CreateUser(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("CreateUser called\n")

	var t models.User

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		fmt.Printf("can't create task: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = e.auth.CreateUser(t)
	if err != nil {
		fmt.Printf("can't create user: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (e *Endpoints) ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ChangeUserRole called\n")

	params := mux.Vars(r)
	userId := params["id"]

	roleChangeRequest := struct {
		Role string
	}{}

	err := json.NewDecoder(r.Body).Decode(&roleChangeRequest)
	if err != nil {
		fmt.Printf("can't change role: %+v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = e.auth.ChangeRole(userId, roleChangeRequest.Role)
	if err != nil {
		fmt.Printf("can't create user: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (e *Endpoints) Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Login called\n")

	r.Body.Close()
	w.WriteHeader(http.StatusOK)
}
