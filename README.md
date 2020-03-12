# SKYPE BOT

## Example
See working example in [example_test.go](example_test.go).
Just set before the following env variables:
- SKYPE_APP_ID `// bot id`
- SKYPE_APP_SECRET `// bot token`
- PORT `// web server port (used for hook)`

## How to find bot app credentials?
- open Azure Portal https://portal.azure.com
- go "App Registrations" https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps/ApplicationsListBlade
- find your bot in the list (my bot app was located inside "Applications from personal account") and open it
- see your SKYPE_APP_ID in "Application (client) ID"
- see your SKYPE_APP_SECRET in the section "Client secrets" of "Manage -> Certificates & secrets"

## Features
- send text messages
- receive text messages (setting of webhook is required)
- receive conversation events (setting of webhook is required)

## TODO
- tests
- add command handler
- send/receive attachments

## Example
```go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/JILeXanDR/skypebot/bot"
	"github.com/JILeXanDR/skypebot/bot/message"
)

func main() {
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

	b.On(bot.EventAddedToContacts, func(activity *bot.Activity) {
		b.Send(activity, message.TextMessage("привет! спасибо что добавил!"))
	})

	if err := b.Run(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/hook", b.WebHookHandler())

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
```
