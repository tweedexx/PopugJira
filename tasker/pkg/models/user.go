package models

const (
	Admin      = "admin"
	Accountant = "accountant"
	Developer  = "developer"
)

type User struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}

func (u *User) CanBeAssignee() bool {
	return u.Role == Developer
}
