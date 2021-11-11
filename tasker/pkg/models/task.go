package models

const (
	Assigned = "assigned"
	Done     = "done"
)

type Task struct {
	Id          uint64 `json:"_id"`
	PublicId    string `json:"id"`
	AssigneeId  string `json:"assignee_id"`
	Description string `json:"description"`
	Fee         int    `json:"fee"`
	Reward      int    `json:"reward"`
	Status      string `json:"status"`
}
