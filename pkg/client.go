package pkg

import (
	"context"
	"errors"
	"sync"

	"github.com/gotd/td/telegram"
)

func Connect(ctx context.Context, appID int, appHash string, method LoginMethod) (Client, error) {
	ctx, cancel := context.WithCancel(ctx)

	c := &client{
		cancel:  cancel,
		done:    make(chan struct{}),
		errChan: make(chan error, 16),
		ready:   make(chan error, 1),
	}

	go c.run(ctx, appID, appHash, method)

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

	Close() error
	Errors() <-chan error
	OnError(ErrorHandler)
}

type client struct {
	cancel  context.CancelFunc
	errChan chan error

	tgClient    *telegram.Client
	chatService *chatService

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

func (c *client) run(ctx context.Context, appID int, appHash string, method LoginMethod) {
	defer close(c.done)

	c.tgClient = telegram.NewClient(appID, appHash, telegram.Options{})

	if err := c.tgClient.Run(ctx, func(ctx context.Context) error {
		if method != nil {
			if err := method.Login(ctx, c.tgClient); err != nil {
				c.signalReady(err)
				return err
			}
		}

		api := c.tgClient.API()
		c.chatService = &chatService{api: api}
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
