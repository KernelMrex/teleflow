package pkg

import (
	"context"
	"errors"
	"fmt"

	tgmessage "github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

var ErrUnknownChat = errors.New("unknown chat")

type Message struct {
	ID   int64
	Text string
}

type MessageService interface {
	SendMessage(ctx context.Context, chatID ChatID, msg *Message) error
}

type messageService struct {
	api   *tg.Client
	peers *peerStore
}

func (s *messageService) SendMessage(ctx context.Context, chatID ChatID, msg *Message) error {
	if msg == nil {
		return errors.New("nil message")
	}

	peer, ok := s.peers.Get(chatID)
	if !ok {
		return fmt.Errorf("%w: %d", ErrUnknownChat, chatID)
	}

	_, err := tgmessage.NewSender(s.api).
		To(peer).
		Text(ctx, msg.Text)

	return err
}
