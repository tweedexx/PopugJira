package main

import (
	"anku/popug-jira/tasker/pkg/adapters"
	"anku/popug-jira/tasker/pkg/config"
	"anku/popug-jira/tasker/pkg/events"
	"anku/popug-jira/tasker/pkg/pricing"
	"anku/popug-jira/tasker/pkg/repos"
	"anku/popug-jira/tasker/pkg/service"
	"anku/popug-jira/tasker/pkg/transport"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	router := mux.NewRouter()

	cfg := config.New()
	ue := make(chan adapters.UserEvent)
	done := make(chan struct{})

	db, err := sqlx.Connect("postgres", cfg.DbUrl())
	if err != nil {
		log.Fatalln(err)
	}

	kafka := adapters.NewKafka(cfg, ue, done)
	pg := adapters.NewPostgres(db)

	taskRepo := repos.New(pg, kafka)

	usrEvents := events.New(ue, pg, taskRepo)
	go usrEvents.Route()

	pricer := pricing.New()

	taskerSvc := service.New(pg, taskRepo, pricer)

	tasker := transport.New(taskerSvc)

	router.Methods(http.MethodPost).Path("/tasks/create").HandlerFunc(tasker.CreateTask)
	router.Methods(http.MethodGet).Path("/users/{id}/tasks").HandlerFunc(tasker.ListTasksForUser)
	router.Methods(http.MethodPost).Path("/tasks/assign").HandlerFunc(tasker.ListTasksForUser)
	router.Methods(http.MethodPost).Path("/tasks/{id}/finish").HandlerFunc(tasker.FinishTask)

	fmt.Printf("start listening on %v\n", ":9124")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9124",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
	fmt.Printf("stop listening on %v\n", ":9124")

	close(done)

	time.Sleep(time.Second)

	fmt.Println("exiting")
}
