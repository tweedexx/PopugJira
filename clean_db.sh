#!/usr/bin/env bash

migrate -path auth-service/migrations -database 'postgres://localhost:5432/auth?sslmode=disable' drop -f
migrate -path auth-service/migrations -database 'postgres://localhost:5432/auth?sslmode=disable' up

migrate -path tasker/migrations -database 'postgres://localhost:5432/tasker?sslmode=disable' drop -f
migrate -path tasker/migrations -database 'postgres://localhost:5432/tasker?sslmode=disable' up
