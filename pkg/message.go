package pkg

import (
	"context"

	"github.com/gotd/td/tg"
)

type Message struct {
	Text string
}

type MessageHandler func(ctx context.Context, chatID ChatID, msg *Message) error

type MessageService interface {
	OnMessage(chatID ChatID, handler MessageHandler)
	SendMessage(ctx context.Context, chatID ChatID, msg *Message) error
}

type messageService struct {
	api *tg.Client
}

func (s *messageService) OnMessage(chatID ChatID, handler MessageHandler) {
	// TODO: implement
	panic("implement me")
}

func (s *messageService) SendMessage(ctx context.Context, chatID ChatID, msg *Message) error {
	// TODO: implement
	panic("implement me")
}
