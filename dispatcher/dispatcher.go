package dispatcher

import (
	"context"
	"reflect"
	"runtime"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/rs/zerolog/log"
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
		matched := true
		for _, filter := range record.filters {
			if !filter(u) {
				matched = false
				break
			}
		}

		if matched {
			pc := reflect.ValueOf(record.handler).Pointer()
			handlerName := runtime.FuncForPC(pc).Name()
			log.Debug().
				Str("Handler name", handlerName).
				Int64("User", u.Message.Sender.UserId).
				Msg("Executing handler.")

			err := record.handler(ctx, api, u)
			return true, err
		}
	}
	return false, nil
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		routers: make([]*Router, 0),
	}
}

type Dispatcher struct {
	routers []*Router
}

func (d *Dispatcher) IncludeRouter(
	router *Router,
) {
	d.routers = append(
		d.routers,
		router,
	)
}

func (d *Dispatcher) Handle(
	ctx context.Context,
	api *maxbot.Api,
	u *schemes.MessageCreatedUpdate,
) error {
	log.Debug().Str("Text", u.GetText()).Msg("Dispatcher received message")
	for _, router := range d.routers {
		handled, err := router.Handle(ctx, api, u)
		if err != nil {
			return err
		}
		if handled {
			return nil
		}
	}
	return nil
}
