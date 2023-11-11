package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
)

const TypeEmail = "email"

func NewSendEmailTask() (*asynq.Task, error) {
	payload, err := json.Marshal(emailPayload{})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEmail, payload), nil
}

type SendEmailProcessor struct{}

func NewSendEmailProcessor() SendEmailProcessor {
	return SendEmailProcessor{}
}

type emailPayload struct{}

func (s SendEmailProcessor) ProcessTask(ctx context.Context, t *asynq.Task) error {
	payload := emailPayload{}

	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v, %w", err, asynq.SkipRetry)
	}

	return nil
}
