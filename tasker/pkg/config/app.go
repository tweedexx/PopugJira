package config

import "fmt"

type App struct {
	KafkaAddress string
	AuthAddress  string

	DbHost     string
	DbPort     int
	DbUser     string
	DbPassword string
	Database   string

	UsersTopic string
	TasksTopic string

	GroupID string
}

func New() App {
	return App{
		KafkaAddress: "localhost:9092",
		DbHost:       "localhost",
		DbUser:       "postgres",
		DbPassword:   "postgres",
		DbPort:       5432,
		Database:     "tasker",
		UsersTopic:   "users",
		TasksTopic:   "tasks_lifecycle",
		GroupID:      "tasker",
	}
}

func (cfg App) DbUrl() string {
	return fmt.Sprintf("user=%v password=%v dbname=%v host=%v port=%d sslmode=disable", cfg.DbUser, cfg.DbPassword, cfg.Database, cfg.DbHost, cfg.DbPort)
}
