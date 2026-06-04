package teleflow

import (
	"context"
	"errors"
	"fmt"

	"github.com/gotd/td/tg"
)

type MessageHandler func(ctx context.Context, chatID ChatID, msg *Message) error

type UpdateMux interface {
	OnChatMessage(chatID ChatID, handler MessageHandler) error

	register(dispatcher *tg.UpdateDispatcher, peerStore *peerStore)
}

func NewUpdateMux() UpdateMux {
	return &updateMux{
		messageHandlers: make(map[ChatID][]MessageHandler),
	}
}

type updateMux struct {
	dispatcher      *tg.UpdateDispatcher
	messageHandlers map[ChatID][]MessageHandler
	peerStore       *peerStore
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

func (r *updateMux) register(dispatcher *tg.UpdateDispatcher, peerStore *peerStore) {
	r.peerStore = peerStore
	r.dispatcher = dispatcher

	r.dispatcher.OnNewMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewMessage) error {
		return r.dispatchMessage(ctx, e, update.Message)
	})
	r.dispatcher.OnNewChannelMessage(func(ctx context.Context, e tg.Entities, update *tg.UpdateNewChannelMessage) error {
		return r.dispatchMessage(ctx, e, update.Message)
	})
}

func (r *updateMux) dispatchMessage(ctx context.Context, e tg.Entities, raw tg.MessageClass) error {
	tgMsg, ok := raw.(*tg.Message)
	if !ok {
		return nil
	}

	chatID, ok := chatIDFromPeer(tgMsg.PeerID)
	if !ok {
		return nil
	}

	if inputPeer, ok := inputPeerFromPeer(tgMsg.PeerID, e); ok {
		r.peerStore.Put(chatID, inputPeer)
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

	return handlerErrors
}
