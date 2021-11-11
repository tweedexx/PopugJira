package repos

import (
	"anku/popug-jira/tasker/pkg/adapters"
	"anku/popug-jira/tasker/pkg/models"
	"anku/popug-jira/tasker/pkg/service"
)

type TaskRepo struct {
	taskStorage service.TaskStorage
	kafka       *adapters.Kafka
}

func New(ts service.TaskStorage, kafka *adapters.Kafka) *TaskRepo {
	return &TaskRepo{
		taskStorage: ts,
		kafka:       kafka,
	}
}

func (t *TaskRepo) GetTaskById(id string) (models.Task, error) {
	return t.taskStorage.GetTaskById(id)
}

func (t *TaskRepo) GetTasksByUserId(userId string) ([]models.Task, error) {
	return t.taskStorage.GetTasksByUserId(userId)
}

func (t *TaskRepo) StoreTask(task models.Task) error {
	err := t.taskStorage.StoreTask(task)
	if err != nil {
		return err
	}

	t.kafka.Send(adapters.TaskCreated, map[string]interface{}{
		"id":          task.PublicId,
		"fee":         task.Fee,
		"reward":      task.Reward,
		"description": task.Description,
	})

	return err
}

func (t *TaskRepo) GetAllTasks() ([]models.Task, error) {
	return t.taskStorage.GetAllTasks()
}

func (t *TaskRepo) ChangeAssignee(taskId string, assigneeId string) error {
	err := t.taskStorage.ChangeAssignee(taskId, assigneeId)
	if err != nil {
		return nil
	}

	tsk, err := t.GetTaskById(taskId)
	if err != nil {
		return err
	}

	t.kafka.Send(adapters.TaskAssigned, map[string]interface{}{
		"taskId":      taskId,
		"assigneeId":  assigneeId,
		"description": tsk.Description,
		"fee":         tsk.Fee,
	})

	return nil
}

func (t *TaskRepo) Finish(taskId string) error {
	err := t.taskStorage.Finish(taskId)
	if err != nil {
		return nil
	}

	tsk, err := t.GetTaskById(taskId)
	if err != nil {
		return err
	}

	t.kafka.Send(adapters.TaskDone, map[string]interface{}{
		"taskId":      taskId,
		"assigneeId":  tsk.AssigneeId,
		"description": tsk.Description,
		"reward":      tsk.Reward,
	})

	return nil
}
