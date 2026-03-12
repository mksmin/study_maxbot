package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	dp "study_bot_go/dispatcher"
	"study_bot_go/filters"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func StartRouter() *dp.Router {
	r := dp.NewRouter()
	r.Message(handleStart, filters.Command("start"))
	//r.Message(anyMessage)

	return r
}

func handleStart(
	ctx context.Context,
	api *maxbot.Api,
	u *schemes.MessageCreatedUpdate,
) error {
	msg := maxbot.NewMessage().
		SetUser(u.Message.Sender.UserId).
		SetText("Привет! Я твой первый бот на Go. Используй /help.")
	return api.Messages.Send(ctx, msg)
}

func HandleMe(ctx context.Context, api *maxbot.Api, u *schemes.MessageCreatedUpdate) error {
	sender := u.Message.Sender

	userDataJSON, _ := json.MarshalIndent(sender, "", "  ")

	fullDataMap := structToMap(sender)
	fullDataJSON, _ := json.MarshalIndent(fullDataMap, "", "  ")

	responseText := fmt.Sprintf(
		"📦 *Стандартный JSON (скрывает пустые поля):*\n```json\n%s\n```\n\n🛠 *Полный дамп (все поля структуры):*\n```json\n%s\n```\n\n_Всего полей в структуре: %d_",
		string(userDataJSON),
		string(fullDataJSON),
		len(fullDataMap),
	)

	msg := maxbot.NewMessage().
		SetUser(sender.UserId).
		SetText(responseText)
	return api.Messages.Send(ctx, msg)
}

func structToMap(item interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	v := reflect.ValueOf(item)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return out
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		if field.PkgPath != "" {
			continue
		}

		out[field.Name] = value.Interface()
	}
	return out
}
