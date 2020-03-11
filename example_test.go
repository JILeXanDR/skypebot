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

	b.Handle(bot.OnTextMessage, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage(fmt.Sprintf(`Ваше сообщение "%s."`, activity.Text())))
	})

	b.Handle(bot.OnAttachment, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage(fmt.Sprintf(`Ваше сообщение "%s."`, activity.Text())))
	})

	b.Handle(bot.Command("ping"), func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("pong"))
	})

	var cdmTestArgs bot.Command = "test :name"
	b.Handle(cdmTestArgs, func(activity *bot.Activity) {
		cmd := cdmTestArgs.Parse(activity)
		b.Send(activity, message.TextMessage(fmt.Sprintf("Получена команда <%s> с параметрами <%v>", cmd.Name, cmd.Args)))
	})

	b.Handle(bot.OnAddedToContacts, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("привет! спасибо что добавил!"))
	})

	// TODO: it should not send message because bot was removed from contacts?
	b.Handle(bot.OnRemovedFromContacts, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("как жаль..."))
	})

	b.Handle(bot.OnAddedToConversation, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("привет! спасибо что добавили в диалог!"))
	})

	// TODO: it should not send message because bot was removed from conversation?
	b.Handle(bot.OnRemovedFromConversation, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("жаль что удалили..."))
	})

	if err := b.Run(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hook", b.WebHookHandler())

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
