package gateway

import (
	"context"

	"unifacema.chatservice/internal/domain/entity"
)

type ChatGateway interface {
	CreateChat(context context.Context, chat *entity.Chat) error
	FindChatByID(context context.Context, chatID string) (*entity.Chat, error)
	SaveChat(context context.Context, chat *entity.Chat) error
}
