package pkg

import (
	"context"
	"io"

	"github.com/gotd/td/tg"
)

type ChatID int64

type ChatRef string

type ChatTypeGroup int

type Chat struct {
	ID   ChatID
	Type ChatTypeGroup
}

type UserID int64

type User struct {
	ID       UserID
	Username string
	IsBot    bool
}

type ChatService interface {
	Resolve(ctx context.Context, ref ChatRef) (Chat, error)
	Join(ctx context.Context, ref ChatRef) (Chat, error)
	Info(ctx context.Context, chatID ChatID) (Chat, error)
	DownloadPhoto(ctx context.Context, chatID ChatID) (io.ReadCloser, error)
	IterParticipants(ctx context.Context, chatID ChatID, handler func(user User) error) error
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

func (s *chatService) Info(ctx context.Context, chatID ChatID) (Chat, error) {
	// TODO
	panic("implement me")
}

func (s *chatService) DownloadPhoto(ctx context.Context, chatID ChatID) (io.ReadCloser, error) {
	// TODO
	panic("implement me")
}

func (s *chatService) IterParticipants(ctx context.Context, chatID ChatID, handler func(user User) error) error {
	// TODO
	panic("implement me")
}
