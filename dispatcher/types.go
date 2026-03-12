package dispatcher

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Handler func(
	ctx context.Context,
	api *maxbot.Api,
	u *schemes.MessageCreatedUpdate,
) error

type Filter func(
	u *schemes.MessageCreatedUpdate,
) bool

type handlerRecord struct {
	handler Handler
	filters []Filter
}
