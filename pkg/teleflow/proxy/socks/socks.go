package socks

import (
	"context"
	"errors"
	"net"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/dcs"
	"golang.org/x/net/proxy"
)

type Option func(p *Proxy)

func WithAuth(username, password string) Option {
	return func(p *Proxy) {
		p.username = username
		p.password = password
	}
}

func New(address string, opts ...Option) *Proxy {
	p := &Proxy{
		address: address,
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}

type Proxy struct {
	address  string
	username string
	password string
}

func (p *Proxy) Configure(opts *telegram.Options) error {
	var auth *proxy.Auth
	if p.username != "" || p.password != "" {
		auth = &proxy.Auth{
			User:     p.username,
			Password: p.password,
		}
	}

	dialer, err := proxy.SOCKS5("tcp", p.address, auth, proxy.Direct)
	if err != nil {
		return err
	}

	contextDialer, ok := dialer.(proxy.ContextDialer)
	if !ok {
		return errors.New("socks5 proxy dialer does not support context")
	}

	opts.Resolver = dcs.Plain(dcs.PlainOptions{
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return contextDialer.DialContext(ctx, network, address)
		},
	})

	return nil
}
