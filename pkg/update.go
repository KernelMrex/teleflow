package pkg

import (
	"context"
	"errors"
	"fmt"

	"github.com/gotd/td/constant"
	"github.com/gotd/td/tg"
)

type MessageHandler func(ctx context.Context, chatID ChatID, msg *Message) error

type UpdateMux interface {
	OnChatMessage(chatID ChatID, handler MessageHandler) error

	register(dispatcher *tg.UpdateDispatcher)
}

func NewUpdateMux() UpdateMux {
	return &updateMux{
		messageHandlers: make(map[ChatID][]MessageHandler),
	}
}

type updateMux struct {
	dispatcher      *tg.UpdateDispatcher
	messageHandlers map[ChatID][]MessageHandler
}

func (r *updateMux) OnChatMessage(chatID ChatID, handler MessageHandler) error {
	if r.dispatcher != nil {
		return fmt.Errorf("mux already registered")
	}

	handlers, ok := r.messageHandlers[chatID]
	if !ok {
		handlers = make([]MessageHandler, 0, 1)
	}
	r.messageHandlers[chatID] = append(handlers, handler)

	return nil
}

func (r *updateMux) register(dispatcher *tg.UpdateDispatcher) {
	r.dispatcher = dispatcher

	r.dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		return r.dispatchMessage(ctx, update.Message)
	})
	r.dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		return r.dispatchMessage(ctx, update.Message)
	})
}

func (r *updateMux) dispatchMessage(ctx context.Context, raw tg.MessageClass) error {
	tgMsg, ok := raw.(*tg.Message)
	if !ok {
		return nil
	}

	chatID, ok := chatIDFromPeer(tgMsg.PeerID)
	if !ok {
		return nil
	}

	tfMsg := &Message{
		Text: tgMsg.Message,
	}

	handlers, ok := r.messageHandlers[chatID]
	if !ok {
		return nil
	}

	if len(handlers) == 0 {
		return nil
	}

	var handlerErrors error
	for _, h := range handlers {
		if err := h(ctx, chatID, tfMsg); err != nil {
			handlerErrors = errors.Join(handlerErrors, err)
		}
	}

	return nil
}

func chatIDFromPeer(peer tg.PeerClass) (ChatID, bool) {
	var id constant.TDLibPeerID

	switch p := peer.(type) {
	case *tg.PeerUser:
		id.User(p.UserID)
	case *tg.PeerChat:
		id.Chat(p.ChatID)
	case *tg.PeerChannel:
		id.Channel(p.ChannelID)
	default:
		return 0, false
	}

	return ChatID(id), true
}
