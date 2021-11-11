package adapters

import (
	"anku/popug-jira/tasker/pkg/config"
	"anku/popug-jira/tasker/pkg/models"
	"context"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

const (
	UserCreated     = "user.created"
	UserRoleChanged = "user.role_changed"
	TaskCreated     = "task.created"
	TaskAssigned    = "task.assigned"
	TaskDone        = "task.done"
)

type UserEvent struct {
	models.User
	EventType string
}

type Kafka struct {
	userEvents chan UserEvent
	done       chan struct{}
	users      *kafka.Reader
	tasks      *kafka.Writer
}

func NewKafka(cfg config.App, ue chan UserEvent, done chan struct{}) *Kafka {

	users := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaAddress},
		GroupID:  cfg.GroupID,
		Topic:    cfg.UsersTopic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	tasks := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaAddress),
		Topic:    cfg.TasksTopic,
		Balancer: &kafka.LeastBytes{},
	}

	k := &Kafka{
		userEvents: ue,
		users:      users,
		tasks:      tasks,
		done:       done,
	}

	go func() {
		<-done
		k.tasks.Close()
	}()

	k.Listen()

	return k
}

func (k *Kafka) Listen() {
	go func() {
		for {
			select {
			case <-k.done:
				if err := k.users.Close(); err != nil {
					fmt.Printf("failed to close reader: %+v", errors.WithStack(err))
				}
			default:
				m, err := k.users.ReadMessage(context.Background())
				if err != nil {
					break
				}
				eventType := string(m.Key)
				var payload models.User

				err = json.Unmarshal(m.Value, &payload)
				if err != nil {
					fmt.Printf("can't unmarshal %v to models.user\n", string(m.Value))
					continue
				}

				k.userEvents <- UserEvent{
					User:      payload,
					EventType: eventType,
				}
			}
		}
	}()
}

func (k *Kafka) Send(eventType string, event interface{}) {

	payload, _ := json.Marshal(event)

	err := k.tasks.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte(eventType),
			Value: payload,
		},
	)

	if err != nil {
		fmt.Printf("can't sent %v because %+v\n", string(payload), errors.WithStack(err))
	}
}
