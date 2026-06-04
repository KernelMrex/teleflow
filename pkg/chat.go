package pkg

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/gotd/td/telegram/deeplink"
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

	return s.chatFromPeer(peer), nil
}

func (s *chatService) Join(ctx context.Context, ref ChatRef) (Chat, error) {
	raw := strings.TrimSpace(string(ref))
	if raw == "" {
		return Chat{}, fmt.Errorf("empty chat ref")
	}

	if deeplink.IsDeeplinkLike(raw) {
		link, err := deeplink.Parse(raw)
		if err != nil {
			return Chat{}, err
		}
		switch link.Type {
		case deeplink.Join:
			peer, err := s.peerManager.JoinLink(ctx, raw)
			if err != nil {
				return Chat{}, err
			}
			return s.chatFromPeer(peer), nil
		case deeplink.Resolve:
			return s.joinPublic(ctx, raw)
		default:
			return Chat{}, fmt.Errorf("unsupported chat ref %q", raw)
		}
	}

	if strings.HasPrefix(raw, "+") {
		peer, err := s.peerManager.ImportInvite(ctx, strings.TrimPrefix(raw, "+"))
		if err != nil {
			return Chat{}, err
		}
		return s.chatFromPeer(peer), nil
	}

	return s.joinPublic(ctx, raw)
}

func (s *chatService) joinPublic(ctx context.Context, ref string) (Chat, error) {
	peer, err := s.peerManager.Resolve(ctx, ref)
	if err != nil {
		return Chat{}, err
	}

	switch p := peer.(type) {
	case tgpeers.Channel:
		if p.Left() {
			if err := p.Join(ctx); err != nil {
				return Chat{}, err
			}
		}
	default:
		return Chat{}, fmt.Errorf("peer %q is not joinable", ref)
	}

	return s.chatFromPeer(peer), nil
}

func (s *chatService) chatFromPeer(peer tgpeers.Peer) Chat {
	chatID := ChatID(peer.TDLibPeerID())
	s.peers.Put(chatID, peer.InputPeer())

	return Chat{
		ID:   chatID,
		Type: chatTypeFromPeer(peer),
	}
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
