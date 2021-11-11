package pricing

import (
	"anku/popug-jira/tasker/pkg/models"
	"math/rand"
)

type Pricer interface {
	Fee(task models.Task) int
	Reward(task models.Task) int
}

type randomPricer struct {
}

func New() Pricer {
	return &randomPricer{}
}

func (r *randomPricer) Fee(_ models.Task) int {
	return rand.Intn(11) + 10
}

func (r *randomPricer) Reward(_ models.Task) int {
	return rand.Intn(21) + 20
}
