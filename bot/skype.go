package bot

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JILeXanDR/skypebot/bot/message"
	"github.com/JILeXanDR/skypebot/skypeapi"
	"io/ioutil"
	"log"
	"net/http"
)

type Config struct {
	AppID     string
	AppSecret string
	Logger    *log.Logger
}

type Bot struct {
	api           *API
	eventHandlers map[string]func(*Activity)
	logger        *log.Logger
}

// Handle processes incoming request and passes it to handler
func (bot *Bot) handleActivity(activity *skypeapi.Activity) error {
	msg := &Activity{activity: activity}

	sender := msg.Sender()
	bot.log(fmt.Sprintf("handling message, details: from=%s (%s), text=%s, group chat=%v", sender.account.Name, sender.account.ID, msg.Text(), msg.IsGroup()))

	if msg.SomeoneWroteToMe() {
		bot.callEventHandlerIfExists(EventMessage, msg)
	} else if msg.AddedToContacts() {
		bot.callEventHandlerIfExists(EventAddedToContacts, msg)
	} else if msg.RemovedFromContacts() {
		bot.callEventHandlerIfExists(EventRemovedFromContacts, msg)
	} else if msg.AddedToConversation() {
		bot.callEventHandlerIfExists(EventAddedToConversation, msg)
	} else if msg.RemovedFromConversation() {
		bot.callEventHandlerIfExists(EventRemovedFromConversation, msg)
	} else {
		bot.log("activity has unknown type and we can't find supported event for it")
	}

	bot.callEventHandlerIfExists(EventAll, msg)

	return nil
}

func (bot *Bot) callEventHandlerIfExists(event Event, activity *Activity) {
	handler, ok := bot.eventHandlers[event.EventID()]
	if ok {
		bot.log(fmt.Sprintf(`calling handler for event "%s"`, event))
		handler(activity)
	}
}

func (bot *Bot) Run() error {
	bot.log(fmt.Sprintf("authenticating..."))
	if err := bot.api.Authenticate(); err != nil {
		bot.log(fmt.Sprintf("authenticating failed: %+v", err))
		return err
	}

	bot.log(fmt.Sprintf("authenticating succeed"))
	return nil
}

func (bot *Bot) WebHookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var activity skypeapi.Activity

		bot.log("hook is called")

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&activity); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "bad activity: %+v", err)
			return
		}

		b, err := json.MarshalIndent(activity, "", "  ")
		if err != nil {
			bot.log(fmt.Sprintf("can't convert incoming activity to json: %+v", err))
		} else {
			bot.log(fmt.Sprintf("received activity: %v", string(b)))
		}

		if err := bot.handleActivity(&activity); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "can't handle inconing activity: %s", err.Error())
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func (bot *Bot) log(text string) {
	bot.logger.Printf("[SKYPE_BOT] %s", text)
}

func (bot *Bot) On(event Event, handler func(*Activity)) {
	bot.log(fmt.Sprintf("setting an event handler for '%s'", event.EventID()))
	bot.eventHandlers[event.EventID()] = handler
}

func (bot *Bot) Send(recipient Recipienter, msg message.Sendable) error {
	switch m := msg.(type) {
	case message.TextMessage:
		bot.log(fmt.Sprintf("sending text message '%s' to <%s>", m, recipient.RecipientID()))
		if err := bot.api.SendToConversation(recipient.RecipientID(), string(m)); err != nil {
			bot.log(fmt.Sprintf("can't sent message: %+v", err))
		}
		return nil
	default:
		return errors.New(fmt.Sprintf("message of type <%s> is not supported", m))
		//case *ImageMessage:
		//	// TODO
		//	original.Attachments = []skypeapi.Attachment{
		//		{
		//			Content: skypeapi.AttachmentContent{
		//				Type: "xxx",
		//			},
		//		},
		//	}
	}
	return nil
}

func New(config Config) *Bot {
	logger := log.New(ioutil.Discard, "", 0)

	if config.Logger != nil {
		logger = config.Logger
	}

	return &Bot{
		api:           newAPI(config.AppID, config.AppSecret),
		eventHandlers: make(map[string]func(*Activity), 0),
		logger:        logger,
	}
}
