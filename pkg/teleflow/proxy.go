package teleflow

import "github.com/gotd/td/telegram"

type Proxy interface {
	Configure(opts *telegram.Options) error
}
