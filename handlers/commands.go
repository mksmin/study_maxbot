package handlers

import (
	"context"
	"study_bot_go/dispatcher"
	"study_bot_go/filters"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/rs/zerolog"
)

func CommandsRouter(logger zerolog.Logger) *dispatcher.Router {
	r := dispatcher.NewRouter(logger)
	r.Message(handleHelp, filters.Command("help"))
	r.Message(handlePing, filters.Command("ping"))
	return r
}

func handleHelp(
	ctx context.Context,
	api *maxbot.Api,
	u *schemes.MessageCreatedUpdate,
) error {
	text := "Доступные команды:\n/start - Начать\n/help - Помощь\n/ping - Понг\n/echo - Повтор\n/me - Данные о пользователе"
	msg := maxbot.NewMessage().SetUser(u.Message.Sender.UserId).SetText(text)
	return api.Messages.Send(ctx, msg)
}

func handlePing(
	ctx context.Context,
	api *maxbot.Api,
	u *schemes.MessageCreatedUpdate,
) error {
	msg := maxbot.NewMessage().SetUser(u.Message.Sender.UserId).SetText("Понг! 🏓")
	return api.Messages.Send(ctx, msg)
}
