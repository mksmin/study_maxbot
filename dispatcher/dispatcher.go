package dispatcher

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/rs/zerolog/log"
)

type Dispatcher struct {
	routers []*Router
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		routers: make([]*Router, 0),
	}
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
	log.Debug().
		Str("text", u.GetText()).
		Msg("Dispatcher received message")

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
