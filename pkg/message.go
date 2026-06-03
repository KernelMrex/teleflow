package pkg

import (
	"context"

	"github.com/gotd/td/tg"
)

type Message struct {
	ID   int64
	Text string
}

type MessageService interface {
	SendMessage(ctx context.Context, chatID ChatID, msg *Message) error
}

type messageService struct {
	api *tg.Client
}

func (s *messageService) SendMessage(ctx context.Context, chatID ChatID, msg *Message) error {
	// TODO: implement
	panic("implement me")
}
