package adapters

import (
	"anku/popug-jira/auth/pkg/config"
	"anku/popug-jira/auth/pkg/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

const (
	UserCreated     = "user.created"
	UserRoleChanged = "user.role_changed"
)

type UserEvent struct {
	models.User
	EventType string
}

type Kafka struct {
	userEvents chan UserEvent
	done       chan struct{}
	users      *kafka.Writer
}

func NewKafka(cfg config.App, done chan struct{}) *Kafka {

	users := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaAddress),
		Topic:    cfg.UsersTopic,
		Balancer: &kafka.LeastBytes{},
	}

	k := &Kafka{
		users: users,
		done:  done,
	}

	go func() {
		<-done
		k.users.Close()
	}()

	return k
}

func (k *Kafka) Send(eventType string, event interface{}) {

	payload, _ := json.Marshal(event)

	err := k.users.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(eventType),
			Value: payload,
		},
	)

	if err != nil {
		fmt.Printf("can't sent %v because %+v\n", string(payload), errors.WithStack(err))
	}
}
