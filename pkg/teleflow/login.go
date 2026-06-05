package teleflow

import (
	"context"
	"errors"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

type CodeFunc func(context.Context) (string, error)

type LoginMethod interface {
	Login(context.Context, *telegram.Client) error
}

func Phone(phone string, codeFn CodeFunc) LoginMethod {
	return &phoneLoginMethod{
		phone:  phone,
		codeFn: codeFn,
	}
}

type phoneLoginMethod struct {
	phone  string
	codeFn CodeFunc
}

func (lm *phoneLoginMethod) Login(ctx context.Context, client *telegram.Client) error {
	flow := auth.NewFlow(auth.CodeOnly(
		lm.phone,
		auth.CodeAuthenticatorFunc(func(ctx context.Context, _ *tg.AuthSentCode) (string, error) {
			if lm.codeFn == nil {
				return "", errors.New("nil code function")
			}
			return lm.codeFn(ctx)
		}),
	), auth.SendCodeOptions{})

	return client.Auth().IfNecessary(ctx, flow)
}

func PhonePassword(phone string, password string, codeFn CodeFunc) LoginMethod {
	return &phonePasswordLoginMethod{
		phone:    phone,
		password: password,
		codeFn:   codeFn,
	}
}

type phonePasswordLoginMethod struct {
	phone    string
	password string
	codeFn   CodeFunc
}

func (lm *phonePasswordLoginMethod) Login(ctx context.Context, client *telegram.Client) error {
	flow := auth.NewFlow(auth.Constant(
		lm.phone,
		lm.password,
		auth.CodeAuthenticatorFunc(func(ctx context.Context, _ *tg.AuthSentCode) (string, error) {
			if lm.codeFn == nil {
				return "", errors.New("nil code function")
			}
			return lm.codeFn(ctx)
		}),
	), auth.SendCodeOptions{})

	return client.Auth().IfNecessary(ctx, flow)
}
