package pkg

import (
	"context"
	"io"

	tgpeers "github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

type ChatID int64

type ChatRef string

type ChatType int

const (
	ChatTypeUnknown ChatType = iota
	ChatTypeUser
	ChatTypeGroup
	ChatTypeChannel
	ChatTypeSupergroup
)

type Chat struct {
	ID   ChatID
	Type ChatType
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
	api         *tg.Client
	peers       *peerStore
	peerManager *tgpeers.Manager
}

func (s *chatService) Resolve(ctx context.Context, ref ChatRef) (Chat, error) {
	peer, err := s.peerManager.Resolve(ctx, string(ref))
	if err != nil {
		return Chat{}, err
	}

	chatID := ChatID(peer.TDLibPeerID())
	s.peers.Put(chatID, peer.InputPeer())

	return Chat{
		ID:   chatID,
		Type: chatTypeFromPeer(peer),
	}, nil
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
