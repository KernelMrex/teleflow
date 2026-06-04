package pkg

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/gotd/td/constant"
	"github.com/gotd/td/telegram/deeplink"
	tgpeers "github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/peers/members"
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

func (s *chatService) Info(ctx context.Context, chatID ChatID) (Chat, error) {
	peer, err := s.peerByChatID(ctx, chatID)
	if err != nil {
		return Chat{}, err
	}

	return s.chatFromPeer(peer), nil
}

func (s *chatService) peerByChatID(ctx context.Context, chatID ChatID) (tgpeers.Peer, error) {
	if inputPeer, ok := s.peers.Get(chatID); ok {
		peer, err := s.peerManager.FromInputPeer(ctx, inputPeer)
		if err != nil {
			return nil, err
		}
		return peer, nil
	}

	peer, err := s.peerManager.ResolveTDLibID(ctx, constant.TDLibPeerID(chatID))
	if err != nil {
		return nil, err
	}

	return peer, nil
}

func (s *chatService) DownloadPhoto(ctx context.Context, chatID ChatID) (io.ReadCloser, error) {
	// TODO
	panic("implement me")
}

func (s *chatService) IterParticipants(ctx context.Context, chatID ChatID, handler func(user User) error) error {
	if handler == nil {
		return fmt.Errorf("nil participant handler")
	}

	peer, err := s.peerByChatID(ctx, chatID)
	if err != nil {
		return err
	}

	var iter members.Members
	switch p := peer.(type) {
	case tgpeers.Chat:
		iter = members.Chat(p)
	case tgpeers.Channel:
		iter = members.Channel(p)
	case tgpeers.User:
		return fmt.Errorf("peer %d is a user, not a chat", chatID)
	default:
		return fmt.Errorf("unsupported peer type %T", peer)
	}

	return iter.ForEach(ctx, func(member members.Member) error {
		return handler(userFromPeer(member.User()))
	})
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

func userFromPeer(user tgpeers.User) User {
	raw := user.Raw()

	username, _ := user.Username()
	return User{
		ID:       UserID(user.ID()),
		Username: username,
		IsBot:    raw.Bot,
	}
}
