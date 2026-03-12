package dispatcher

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"study_bot_go/internal/bot"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/rs/zerolog"
)

type Router struct {
	handlers []handlerRecord
	logger   zerolog.Logger
}

func NewRouter(logger zerolog.Logger) *Router {
	return &Router{
		handlers: make([]handlerRecord, 0),
		logger:   logger,
	}
}

func (r *Router) Message(
	handler bot.Handler,
	filters ...bot.Filter,
) {
	r.handlers = append(
		r.handlers,
		handlerRecord{
			handler: handler,
			filters: filters,
		})
}

func (r *Router) Handle(
	ctx context.Context,
	api *maxbot.Api,
	u *schemes.MessageCreatedUpdate,
) (bool, error) {
	if err := ctx.Err(); err != nil {
		return false, fmt.Errorf("context cancelled before processing: %w", err)
	}

	for _, record := range r.handlers {
		if r.match(u, record.filters) {
			return r.execute(ctx, api, u, record.handler)
		}
	}
	return false, nil
}

func (r *Router) match(
	u *schemes.MessageCreatedUpdate,
	filters []bot.Filter,
) bool {
	for _, filter := range filters {
		if !filter(u) {
			return false
		}
	}
	return true
}

func (r *Router) execute(
	ctx context.Context,
	api *maxbot.Api,
	u *schemes.MessageCreatedUpdate,
	handler bot.Handler,
) (bool, error) {
	pc := reflect.ValueOf(handler).Pointer()
	handlerName := runtime.FuncForPC(pc).Name()

	r.logger.Debug().
		Str("handler", handlerName).
		Int64("user_id", u.Message.Sender.UserId).
		Msg("Executing handler")

	if err := handler(ctx, api, u); err != nil {
		return true, fmt.Errorf("handler '%s' failed: %w", handlerName, err)
	}
	return true, nil
}
