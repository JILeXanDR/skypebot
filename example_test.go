package main

import (
	"fmt"
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
	})

	b.On(bot.EventMessage, func(activity *bot.IncomingActivity) {
		b.Reply(activity, bot.NewTextMessage(fmt.Sprintf(`Ваше сообщение "%s."`, activity.Text())))
	})

	b.On(bot.EventAddedToContacts, func(activity *bot.IncomingActivity) {
		b.Reply(activity, bot.NewTextMessage("привет! спасибо что добавил!"))
	})

	// TODO: it should not send message because bot was removed from contacts?
	b.On(bot.EventRemovedFromContacts, func(activity *bot.IncomingActivity) {
		b.Reply(activity, bot.NewTextMessage("как жаль..."))
	})

	b.On(bot.EventAddedToConversation, func(activity *bot.IncomingActivity) {
		b.Reply(activity, bot.NewTextMessage("привет! спасибо что добавили в диалог!"))
	})

	// TODO: it should not send message because bot was removed from conversation?
	b.On(bot.EventRemovedFromConversation, func(activity *bot.IncomingActivity) {
		b.Reply(activity, bot.NewTextMessage("жаль что удалили..."))
	})

	if err := b.Run(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hook", b.WebHookHandler())

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
