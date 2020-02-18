package bot

import (
	"encoding/json"
	"fmt"
	"github.com/JILeXanDR/skypebot/skypeapi"
	"log"
	"net/http"
)

func newIncomingActivity(activity *skypeapi.Activity) *IncomingActivity {
	return &IncomingActivity{activity: activity}
}

type Config struct {
	AppID     string
	AppSecret string
}

type Bot struct {
	config        *Config
	token         *skypeapi.TokenResponse
	api           *API
	eventHandlers map[Event]func(*IncomingActivity)
}

// Handle processes incoming request and passes it to handler
func (bot *Bot) handleActivity(activity *skypeapi.Activity) error {
	b, err := json.MarshalIndent(activity, "", "  ")
	if err != nil {
		log.Printf("can't convert incoming activity to json: %+v", err)
	} else {
		log.Printf("received incoming activity: %v", string(b))
	}

	message := newIncomingActivity(activity)
	log.Printf("incoming message details: from=%s (%s), text=%s, group chat=%v", message.FromUser().Name, message.FromUser().ID, message.Text(), message.IsGroup())

	if message.SomeoneWroteToMe() {
		bot.callEventHandlerIfExists(EventMessage, message)
	} else if message.AddedToContacts() {
		bot.callEventHandlerIfExists(EventAddedToContacts, message)
	} else if message.RemovedFromContacts() {
		bot.callEventHandlerIfExists(EventRemovedFromContacts, message)
	} else if message.AddedToConversation() {
		bot.callEventHandlerIfExists(EventAddedToConversation, message)
	} else if message.RemovedFromConversation() {
		bot.callEventHandlerIfExists(EventRemovedFromConversation, message)
	} else {
		log.Printf("activity has unknown type and we can't find supported event for it")
	}

	bot.callEventHandlerIfExists(EventAll, message)

	return nil
}

func (bot *Bot) SendMessageToConversation(conversationID string, text string) error {
	return bot.api.SendToConversation(conversationID, text)
}

func (bot *Bot) callEventHandlerIfExists(event Event, activity *IncomingActivity) {
	handler, ok := bot.eventHandlers[event]
	if ok {
		log.Printf(`calling handler for event "%s"`, event)
		handler(activity)
	}
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

func (bot *Bot) On(event Event, handler func(*IncomingActivity)) {
	bot.eventHandlers[event] = handler
}

func (bot *Bot) Reply(activity *IncomingActivity, message Sendable) {
	original := activity.Full()
	switch msg := message.(type) {
	case *TextMessage:
		skypeapi.SendReplyMessage(original, msg.Text, bot.api.token.AccessToken)
	case *ImageMessage:
		// TODO
		original.Attachments = []skypeapi.Attachment{
			{
				Content: skypeapi.AttachmentContent{
					Type: "xxx",
				},
			},
		}
		skypeapi.SendActivityRequest(original, original.ServiceURL, bot.api.token.AccessToken)
	}
}

func (bot *Bot) Send(recipient Recipient, message Sendable) error {
	switch msg := message.(type) {
	case *TextMessage:
		return bot.SendMessageToConversation(recipient.ConversationID(), msg.Text)
	}
	return nil
}

func New(config Config) *Bot {
	return &Bot{
		config:        &config,
		api:           newAPI(config.AppID, config.AppSecret),
		eventHandlers: make(map[Event]func(*IncomingActivity), 0),
	}
}
