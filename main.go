package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"study_bot_go/dispatcher"
	"study_bot_go/handlers"
)

func main() {
	logger := log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: "15:04:05"},
	)
	log.Logger = logger

	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM,
		os.Interrupt,
	)
	defer stop()

	if err := godotenv.Load(".env"); err != nil {
		log.Warn().Msg("Could not load .env file")
	}

	api, err := maxbot.New(
		os.Getenv("MAX_BOT_TOKEN"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("API initialization failed")
	}

	disp := dispatcher.NewDispatcher(logger)

	disp.IncludeRouter(handlers.StartRouter(logger))
	disp.IncludeRouter(handlers.CommandsRouter(logger))

	log.Info().Msg("Бот успешно запущен!")

	for upd := range api.GetUpdates(ctx) {
		switch u := upd.(type) {
		case *schemes.BotStartedUpdate:
			log.Info().
				Int64("chat_id", u.ChatId).
				Msg("Bot started in chat")

			msg := maxbot.NewMessage().
				SetChat(u.ChatId).
				SetText("Добро пожаловать!")
			_ = api.Messages.Send(ctx, msg)

		case *schemes.MessageCreatedUpdate:
			if err := disp.Handle(ctx, api, u); err != nil {
				log.Error().Err(err).Msg("Ошибка при обработке сообщения")
			}
		}
	}
	log.Info().Msg("Завершение работы по сигналу...")
}
