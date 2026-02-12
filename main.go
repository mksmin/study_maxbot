package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		os.Interrupt,
	)
	defer stop()
	defer fmt.Println("Stopping app")

	api, _ := maxbot.New("")

	info, err := api.Bots.GetBot(ctx)
	fmt.Printf("Get me: %#v %#v", info, err)

	for upd := range api.GetUpdates(ctx) {
		switch upd := upd.(type) {
		case *schemes.MessageCreatedUpdate:
			message := maxbot.NewMessage().
				SetChat(upd.Message.Recipient.ChatId).
				SetText("Hello from bot")

			err := api.Messages.Send(ctx, message)
			if err != nil {
			}
		}
	}
}
