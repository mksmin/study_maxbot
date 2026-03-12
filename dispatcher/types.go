package dispatcher

import (
	"study_bot_go/internal/bot"
)

type handlerRecord struct {
	handler bot.Handler
	filters []bot.Filter
}
