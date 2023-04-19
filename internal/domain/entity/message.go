package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
	tiktoken_go "github.com/j178/tiktoken-go"
)

type Message struct {
	ID        string
	Role      string
	Content   string
	Tokens    int
	Model     *Model
	CreatedAt time.Time
}

func NewMessage(role, content string, model *Model) (*Message, error) {
	totalTokens := tiktoken_go.CountTokens(model.GetModelName())
	msg := &Message{
		ID:        uuid.New().String(),
		Role:      role,
		Content:   content,
		Tokens:    totalTokens,
		Model:     model,
		CreatedAt: time.Now(),
	}
	if err := msg.Validate(); err != nil {
		return nil, err
	}
	return msg, nil
}

func (message *Message) Validate() error {
	if message.Role != "user" && message.Role != "system" && message.Role != "assistant" {
		return errors.New("invalid role")
	}
	if message.Content == "" {
		return errors.New("content is empty")
	}
	if message.CreatedAt.IsZero() {
		return errors.New("invalid created at")
	}
	return nil
}

func (message *Message) GetQuantitytokens() int {
	return message.Tokens
}
