//TODO actually make a dockerfile
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate -path migrations -database 'postgres://localhost:5432/tasker?sslmode=disable' up