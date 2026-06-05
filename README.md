# Teleflow

Teleflow is a small Go library for building server-side Telegram clients on top of [gotd/td](https://github.com/gotd/td). It is intended for backend adapters with HTTP APIs, user action emulation, and other Golang services that need to work with Telegram chats, messages, and updates.

## Status

> [!WARNING]
> Teleflow is currently in `v0.*`. The API is not stable yet: exported types, methods, package layout, and behavior may change between minor versions until a `v1.0.0` release.

## License

Teleflow is open-source software released under the [MIT License](LICENSE).

## Installation

```sh
go get github.com/kernelmrex/teleflow/pkg/teleflow@latest
```

You can also download the source code from GitHub:

- [Repository](https://github.com/kernelmrex/teleflow)
- [Latest source archive](https://github.com/kernelmrex/teleflow/archive/refs/heads/main.zip)
- [Releases](https://github.com/kernelmrex/teleflow/releases)

## Usage

```go
package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/kernelmrex/teleflow/pkg/teleflow"
)

func main() {
	ctx := context.Background()

	appID, err := strconv.Atoi(os.Getenv("TELEGRAM_APP_ID"))
	if err != nil {
		panic("TELEGRAM_APP_ID must be an integer")
	}

	appHash := os.Getenv("TELEGRAM_APP_HASH")
	phone := os.Getenv("TELEGRAM_PHONE")
	chatRef := os.Getenv("TELEGRAM_CHAT")

	client, err := teleflow.Connect(
		ctx,
		appID,
		appHash,
		teleflow.WithLoginMethod(teleflow.Phone(phone, readCode)),
	)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	chats, err := client.Chats()
	if err != nil {
		panic(err)
	}

	chat, err := chats.Resolve(ctx, teleflow.ChatRef(chatRef))
	if err != nil {
		panic(err)
	}

	messages, err := client.Messages()
	if err != nil {
		panic(err)
	}

	if err := messages.SendMessage(ctx, chat.ID, &teleflow.Message{
		Text: "Hello from Teleflow",
	}); err != nil {
		panic(err)
	}
}

func readCode(ctx context.Context) (string, error) {
	fmt.Print("Telegram code: ")

	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}
```

Run it with:

```sh
export TELEGRAM_APP_ID=123456
export TELEGRAM_APP_HASH=your_app_hash
export TELEGRAM_PHONE=+10000000000
export TELEGRAM_CHAT=username_or_chat_link
go run .
```
