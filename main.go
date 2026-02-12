package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
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

	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	api, _ := maxbot.New(os.Getenv("MAX_BOT_TOKEN"))

	info, err := api.Bots.GetBot(ctx)
	bot_info, _ := json.MarshalIndent(info, "", " ")
	log.Println("Bot info:")
	log.Println(string(bot_info))
	log.Printf("Err %v", err)

	for upd := range api.GetUpdates(ctx) {
		switch u := upd.(type) {
		case *schemes.BotStartedUpdate:
			message := maxbot.NewMessage().
				SetChat(u.ChatId).
				SetText("Добро пожаловать в бота")

			fmt.Printf("User %v from chat %v ", u.User, u.ChatId)

			err := api.Messages.Send(ctx, message)
			if err != nil {
			}
		case *schemes.MessageCreatedUpdate:
			out := "bot прочитал текст: " + u.GetText()

			switch u.GetCommand() {
			case "/start":
				out = "Команда: " + u.GetCommand()
				message := maxbot.NewMessage().
					SetUser(u.Message.Sender.UserId).
					SetText(out)
				err := api.Messages.Send(ctx, message)
				if err != nil {
				}
			default:
				userMessage := u.Message
				b, _ := json.MarshalIndent(userMessage, "", "  ")
				fmt.Println()
				fmt.Println(fmt.Sprintf("Get message: %v", string(b)))

				message := maxbot.NewMessage().
					SetUser(userMessage.Sender.UserId).
					Reply("Hello from bot", userMessage)

				err := api.Messages.Send(ctx, message)
				if err != nil {
				}
			}

		}
	}
}
