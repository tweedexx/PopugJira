package main

import (
	"anku/popug-jira/auth/pkg/adapters"
	"anku/popug-jira/auth/pkg/config"
	"anku/popug-jira/auth/pkg/repos"
	"anku/popug-jira/auth/pkg/service"
	"anku/popug-jira/auth/pkg/transport"
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
	done := make(chan struct{})

	db, err := sqlx.Connect("postgres", cfg.DbUrl())
	if err != nil {
		log.Fatalln(err)
	}

	kafka := adapters.NewKafka(cfg, done)
	pg := adapters.NewPostgres(db)

	userRepo := repos.New(pg, kafka)

	taskerSvc := service.New(userRepo)

	users := transport.New(taskerSvc)

	router.Methods(http.MethodPost).Path("/users/create").HandlerFunc(users.CreateUser)
	router.Methods(http.MethodPut).Path("/users/{id}/role").HandlerFunc(users.ChangeUserRole)
	router.Methods(http.MethodPost).Path("/login").HandlerFunc(users.Login)

	fmt.Printf("start listening on %v\n", ":9123")

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:9123",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
	fmt.Printf("stop listening on %v\n", ":9123")

	close(done)

	time.Sleep(time.Second)

	fmt.Println("exiting")
}
