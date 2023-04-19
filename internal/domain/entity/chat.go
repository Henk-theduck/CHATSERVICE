package entity

import (
	"errors"

	"github.com/google/uuid"
)

type ChatConfig struct {
	Model            *Model
	Temperature      float32  // 0.0 to 1.0
	TopP             float32  // 0.0 to 1.0 - to a low value, like 0.1, the model will be very conservative in its word choices, and will tend to generate relatively predictable prompts
	N                int      // number of messages to generate
	Stop             []string // list of tokens to stop on
	MaxTokens        int      // number of tokens to generate
	PresencePenalty  float32  // -2.0 to 2.0 - Number between -2.0 and 2.0. Positive values penalize new tokens based on whether they appear in the text so far, increasing the model's likelihood to talk about new topics.
	FrequencyPenalty float32  // -2.0 to 2.0 - Number between -2.0 and 2.0. Positive values penalize new tokens based on their existing frequency in the text so far, increasing the model's likelihood to talk about new topics.
}

type Chat struct {
	ID                   string
	UserID               string
	InitialSystemMessage *Message
	Messages             []*Message
	ErasedMessages       []*Message
	Status               string
	TokenUsage           int
	Config               *ChatConfig
}

func NewChat(userID string, initialSystemMessage *Message, chatConfig *ChatConfig) (*Chat, error) {
	chat := &Chat{
		ID:                   uuid.New().String(),
		UserID:               userID,
		InitialSystemMessage: initialSystemMessage,
		Status:               "active",
		Config:               chatConfig,
		TokenUsage:           0,
	}

	chat.AddMessage(initialSystemMessage)

	if e := chat.Validate(); e != nil {
		return nil, e
	}
	return chat, nil
}

func (chat *Chat) Validate() error {
	if chat.UserID == "" {
		return errors.New("user id is empty")
	}
	if chat.Status != "active" && chat.Status != "ended" {
		return errors.New("invalid status")
	}
	if chat.Config.Temperature < 0 || chat.Config.Temperature > 2 {
		return errors.New("invalid temperature")
	}
	return nil
}

func (chat *Chat) AddMessage(message *Message) error {
	if chat.Status == "ended" {
		return errors.New("chat is ended. no more messages allowed")
	}

	for {
		if chat.Config.Model.GetMaxTokens() >= message.GetQuantitytokens()+chat.TokenUsage {
			chat.Messages = append(chat.Messages, message)
			chat.RefreshTokenUsage()
			break
		}
		chat.ErasedMessages = append(chat.ErasedMessages, chat.Messages[0])
		chat.Messages = chat.Messages[1:]
		chat.RefreshTokenUsage()
	}
	return nil
}

func (chat *Chat) GetMessage() []*Message {
	return chat.Messages

}

func (chat *Chat) CountMessages() int {
	return len(chat.Messages)
}

func (chat *Chat) End() {
	chat.Status = "ended"
}

func (chat *Chat) RefreshTokenUsage() {
	chat.TokenUsage = 0
	for message := range chat.Messages {
		chat.TokenUsage += chat.Messages[message].GetQuantitytokens()
	}
}
