package channel

import (
	"github.com/hibiken/asynq"

	"github.com/flan6/microservice-messaging/internal/consumer/tasks"
)

type Channel struct {
	worker *asynq.Client
}

func NewChannel(w *asynq.Client) Channel {
	return Channel{w}
}

func (c Channel) EnqueueSendEmail() error {
	task, err := tasks.NewSendEmailTask()
	if err != nil {
		return err
	}

	_, err = c.worker.Enqueue(task)
	return err
}
