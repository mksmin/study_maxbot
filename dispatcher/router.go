package dispatcher

import (
	"context"
	"reflect"
	"runtime"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/rs/zerolog/log"
)

type Router struct {
	handlers []handlerRecord
}

func NewRouter() *Router {
	return &Router{
		handlers: make([]handlerRecord, 0),
	}
}

func (r *Router) Message(
	handler Handler,
	filters ...Filter,
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
	for _, record := range r.handlers {
		if r.match(u, record.filters) {
			return r.execute(ctx, api, u, record.handler)
		}
	}
	return false, nil
}

func (r *Router) match(u *schemes.MessageCreatedUpdate, filters []Filter) bool {
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
	handler Handler,
) (bool, error) {
	pc := reflect.ValueOf(handler).Pointer()
	handlerName := runtime.FuncForPC(pc).Name()

	log.Debug().
		Str("handler", handlerName).
		Int64("user_id", u.Message.Sender.UserId).
		Msg("Executing handler")

	err := handler(ctx, api, u)
	return true, err
}
