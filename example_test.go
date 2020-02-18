package main

import (
	"fmt"
	"github.com/JILeXanDR/skypebot/bot/message"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/JILeXanDR/skypebot/bot"
)

func TestExample(t *testing.T) {
	b := bot.New(bot.Config{
		AppID:     os.Getenv("SKYPE_APP_ID"),
		AppSecret: os.Getenv("SKYPE_APP_SECRET"),
		Logger:    log.New(os.Stdout, "", 0),
	})

	b.On(bot.EventMessage, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage(fmt.Sprintf(`Ваше сообщение "%s."`, activity.Text())))
	})

	b.On(bot.Command("ping"), func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("pong"))
	})

	b.On(bot.Command("cmd :id"), func(activity *bot.Activity) {
		cmd := activity.Command()
		name := cmd.Name()
		args := cmd.Args()
		b.Send(activity, message.TextMessage(fmt.Sprintf("Получена команда <%s> с параметрами <%v>", name, args)))
	})

	b.On(bot.EventAddedToContacts, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("привет! спасибо что добавил!"))
	})

	// TODO: it should not send message because bot was removed from contacts?
	b.On(bot.EventRemovedFromContacts, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("как жаль..."))
	})

	b.On(bot.EventAddedToConversation, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("привет! спасибо что добавили в диалог!"))
	})

	// TODO: it should not send message because bot was removed from conversation?
	b.On(bot.EventRemovedFromConversation, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("жаль что удалили..."))
	})

	if err := b.Run(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hook", b.WebHookHandler())

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
