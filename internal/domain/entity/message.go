package entity

import (
	"time"

	"github.com/google/uuid"
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
	msg := &Message{
		ID:        uuid.New().String(),
		Role:      role,
		Content:   content,
		Tokens:    0,
		Model:     model,
		CreatedAt: time.Now(),
	}
	return msg, nil
}
