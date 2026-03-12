package filters

import (
	"strings"
	"study_bot_go/dispatcher"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func Command(name string) dispatcher.Filter {
	return func(u *schemes.MessageCreatedUpdate) bool {
		command := u.GetCommand()
		if len(command) > 0 && command[0] == '/' {
			return command[1:] == name

		}
		return false
	}
}

func Text(target string) dispatcher.Filter {
	return func(u *schemes.MessageCreatedUpdate) bool {
		return u.GetText() == target
	}
}

func Contains(substring string) dispatcher.Filter {
	return func(u *schemes.MessageCreatedUpdate) bool {
		return strings.Contains(u.GetText(), substring)
	}
}
