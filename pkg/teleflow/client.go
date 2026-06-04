package teleflow

import (
	"context"
	"errors"
	"sync"

	"github.com/gotd/td/telegram"
	tgpeers "github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

type ClientOption func(*client)

func WithLoginMethod(method LoginMethod) ClientOption {
	return func(c *client) {
		c.loginMethod = method
	}
}

func WithUpdateMux(mux UpdateMux) ClientOption {
	return func(c *client) {
		c.updateMux = mux
	}
}

func Connect(ctx context.Context, appID int, appHash string, opts ...ClientOption) (Client, error) {
	ctx, cancel := context.WithCancel(ctx)

	c := &client{
		cancel:  cancel,
		done:    make(chan struct{}),
		errChan: make(chan error, 16),
		ready:   make(chan error, 1),
	}

	for _, opt := range opts {
		opt(c)
	}

	go c.run(ctx, appID, appHash)

	select {
	case err := <-c.ready:
		if err != nil {
			cancel()
			<-c.done
			return nil, err
		}
		return c, nil
	case <-ctx.Done():
		cancel()
		<-c.done
		return nil, ctx.Err()
	}
}

type Client interface {
	Chats() (ChatService, error)
	Messages() (MessageService, error)

	Close() error
	Errors() <-chan error
	OnError(ErrorHandler)
}

type client struct {
	cancel  context.CancelFunc
	errChan chan error

	tgClient       *telegram.Client
	chatService    *chatService
	messageService *messageService

	loginMethod LoginMethod
	updateMux   UpdateMux

	done      chan struct{}
	ready     chan error
	readyOnce sync.Once
}

func (c *client) Close() error {
	c.cancel()
	<-c.done
	return nil
}

func (c *client) Errors() <-chan error {
	return c.errChan
}

func (c *client) OnError(handler ErrorHandler) {
	go func() {
		for {
			select {
			case err := <-c.errChan:
				handler(err)
			case <-c.done:
				return
			}
		}
	}()
}

func (c *client) Chats() (ChatService, error) {
	if c.chatService == nil {
		return nil, errors.New("client is not ready")
	}
	return c.chatService, nil
}

func (c *client) Messages() (MessageService, error) {
	if c.messageService == nil {
		return nil, errors.New("client is not ready")
	}
	return c.messageService, nil
}

func (c *client) run(ctx context.Context, appID int, appHash string) {
	defer close(c.done)

	peerStore := newPeerStore()

	tgOptions := telegram.Options{}

	if c.updateMux != nil {
		dispatcher := tg.NewUpdateDispatcher()
		c.updateMux.register(&dispatcher, peerStore)
		tgOptions.UpdateHandler = dispatcher
	} else {
		tgOptions.NoUpdates = true
	}

	c.tgClient = telegram.NewClient(appID, appHash, tgOptions)

	if err := c.tgClient.Run(ctx, func(ctx context.Context) error {
		if c.loginMethod != nil {
			if err := c.loginMethod.Login(ctx, c.tgClient); err != nil {
				c.signalReady(err)
				return err
			}
		}

		api := c.tgClient.API()
		peerManager := tgpeers.Options{}.Build(api)

		c.chatService = &chatService{
			api:         api,
			peers:       peerStore,
			peerManager: peerManager,
		}
		c.messageService = &messageService{
			api:   api,
			peers: peerStore,
		}
		c.signalReady(nil)

		<-ctx.Done()
		return ctx.Err()
	}); err != nil {
		c.signalReady(err)
		c.errChan <- err
	}
}

func (c *client) signalReady(err error) {
	c.readyOnce.Do(func() {
		c.ready <- err
	})
}

type ErrorHandler func(error)
