package models

const (
	Admin      = "admin"
	Accountant = "accountant"
	Developer  = "developer"
)

type User struct {
	Id          uint64 `json:"_id"`
	PublicId    string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	Role        string `json:"role"`
	Credentials string `json:"credentials"`
}
