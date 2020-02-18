package bot

import (
	"encoding/json"
	"fmt"
	"github.com/JILeXanDR/skypebot/skypeapi"
	"log"
	"net/http"
)

var defaultMessageHandler MessageHandler = func(message *IncomingActivity, replyWithText MessageReply) error {
	log.Printf(`received message "%v", but message handler is not set and bot cant process it`, message.Text())
	return nil
}

func newIncomingActivity(activity *skypeapi.Activity) *IncomingActivity {
	return &IncomingActivity{activity: activity}
}

type Config struct {
	AppID     string
	AppSecret string
	OnMessage MessageHandler
}

type Bot struct {
	config *Config
	token  *skypeapi.TokenResponse
	api    *API
}

// Handle processes incoming request and passes it to handler
func (bot *Bot) handleActivity(activity *skypeapi.Activity) error {
	message := newIncomingActivity(activity)
	log.Printf("received incoming message: %v, group chat=%v", message.Text(), message.IsGroup())
	return bot.config.OnMessage(message, func(message Sendable) error {
		switch msg := message.(type) {
		case *TextMessage:
			return skypeapi.SendReplyMessage(activity, msg.Text, bot.api.token.AccessToken)
		case *ImageMessage:
			activity.Attachments = []skypeapi.Attachment{
				{
					Content: skypeapi.AttachmentContent{
						Type: "xxx",
					},
				},
			}
			return skypeapi.SendActivityRequest(activity, activity.ServiceURL, bot.api.token.AccessToken)
			//return skypeapi.SendReplyMessage(activity, msg.Text, bot.api.token.AccessToken)
		}
		return nil
	})
}

func (bot *Bot) SendMessageToConversation(conversationID string, text string) error {
	return bot.api.SendToConversation(conversationID, text)
}

func (bot *Bot) Run() error {
	return bot.api.Authenticate()
}

func (bot *Bot) WebHookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var activity skypeapi.Activity

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&activity); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad activity: %+v", err)
			return
		}

		if err := bot.handleActivity(&activity); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "can't handle inconing activity: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func New(config Config) *Bot {
	if config.OnMessage == nil {
		config.OnMessage = defaultMessageHandler
	}
	return &Bot{
		config: &config,
		api:    newAPI(config.AppID, config.AppSecret),
	}
}
