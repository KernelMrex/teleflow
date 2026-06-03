package pkg

import (
	"context"

	"github.com/gotd/td/tg"
)

type ChatID int64

type ChatRef string

type ChatTypeGroup int

type Chat struct {
	ID   ChatID
	Type ChatTypeGroup
}

type ChatService interface {
	Resolve(ctx context.Context, ref ChatRef) (Chat, error)
	Join(ctx context.Context, ref ChatRef) (Chat, error)
}

type chatService struct {
	api *tg.Client
}

func (s *chatService) Resolve(ctx context.Context, ref ChatRef) (Chat, error) {
	// TODO
	panic("implement me")
}

func (s *chatService) Join(ctx context.Context, ref ChatRef) (Chat, error) {
	// TODO
	panic("implement me")
}
