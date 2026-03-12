package dispatcher

import (
	"context"
	"fmt"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
	"github.com/rs/zerolog"
)

type Dispatcher struct {
	routers []*Router
	logger  zerolog.Logger
}

func NewDispatcher(logger zerolog.Logger) *Dispatcher {
	return &Dispatcher{
		routers: make([]*Router, 0),
		logger:  logger,
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
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("dispatcher context error: %w", err)
	}

	d.logger.Debug().
		Str("text", u.GetText()).
		Msg("Dispatcher received message")

	for _, router := range d.routers {
		handled, err := router.Handle(ctx, api, u)
		if err != nil {
			return fmt.Errorf("router handle failed: %w", err)
		}
		if handled {
			return nil
		}
	}
	return nil
}
