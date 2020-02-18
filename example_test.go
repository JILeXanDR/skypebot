package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/JILeXanDR/skypebot/bot"

	"github.com/pkg/errors"
)

func TestExample(t *testing.T) {
	b := bot.New(bot.Config{
		AppID:     os.Getenv("SKYPE_APP_ID"),
		AppSecret: os.Getenv("SKYPE_APP_SECRET"),
		OnMessage: func(activity *bot.IncomingActivity, reply bot.MessageReply) error {
			if activity.SomeoneWroteToMe() {
				return reply(bot.NewTextMessage(fmt.Sprintf(`Не понимаю, что такое "%s"?`, activity.Text())))
			} else if activity.AddedToContacts() {
				return reply(bot.NewTextMessage("привет! спасибо что добавил!"))
			} else if activity.RemovedFromContacts() {
				return reply(bot.NewTextMessage("как жаль..."))
			} else if activity.AddedToConversation() {
				return reply(bot.NewTextMessage("привет! спасибо что добавили в диалог!"))
			} else if activity.RemovedFromConversation() {
				return reply(bot.NewTextMessage("жаль что удалили..."))
			}
			return errors.New("the data was received via hook can not be processed")
		},
	})

	if err := b.Run(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hook", b.WebHookHandler())

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
