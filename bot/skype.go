package bot

import (
	"encoding/base64"
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
	eventHandlers map[Event]func(*Activity)
	logger        *log.Logger
}

// Handle processes incoming request and passes it to handler
func (bot *Bot) handleActivity(activity *skypeapi.Activity) error {
	msg := &Activity{activity: activity}

	sender := msg.Sender()
	bot.log(fmt.Sprintf("handling message, details: from=%s (%s), text=%s, group chat=%v", sender.account.Name, sender.account.ID, msg.Text(), msg.IsGroup()))

	if bot.lookupCommand(msg) {
		return nil
	} else if msg.SomeoneWroteToMe() {
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
	handler, ok := bot.eventHandlers[event]
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
	bot.eventHandlers[event] = handler
}

func (bot *Bot) Send(recipient Recipienter, msg message.Sendable) error {
	switch m := msg.(type) {
	case message.TextMessage:
		bot.log(fmt.Sprintf("sending text message '%s' to <%s>", m, recipient.RecipientID()))
		if err := bot.api.SendToConversation(recipient.RecipientID(), string(m)); err != nil {
			bot.log(fmt.Sprintf("can't sent message: %+v", err))
		}
		return nil
	case *message.ImageMessage:
		// https://docs.microsoft.com/en-us/azure/bot-service/rest-api/bot-framework-rest-connector-add-media-attachments?view=azure-bot-service-4.0#add-a-media-attachment

		buf, err := ioutil.ReadAll(m.Reader)
		if err != nil {
			return err
		}

		base64Image := `data:image/png;base64,` + base64.StdEncoding.EncodeToString(buf)

		activity := skypeapi.Activity{
			Type: "message",
			Recipient: skypeapi.ChannelAccount{
				ID: recipient.RecipientID(),
			},
			Conversation: skypeapi.ConversationAccount{
				ID: recipient.RecipientID(),
			},
			Text: "image...",
			Attachments: []skypeapi.Attachment{
				{
					ContentType: "image/png",
					ContentUrl:  base64Image,
					Name:        "image.png",
				},
			},
		}

		if err := bot.api.SendActivity(&activity); err != nil {
			bot.log(fmt.Sprintf("can't sent message with image: %+v", err))
		}
	default:
		return errors.New(fmt.Sprintf("message of type <%s> is not supported", m))
	}
	return nil
}

func (bot *Bot) sendActivity(activity *skypeapi.Activity) error {
	bot.log(fmt.Sprintf("sending activity: %+v", *activity))
	return bot.api.SendActivity(activity)
}

func (bot *Bot) SendActions(recipient Recipienter, text string, actions []skypeapi.CardAction) error {
	activity := skypeapi.Activity{
		Type: "message",
		Conversation: skypeapi.ConversationAccount{
			ID: recipient.RecipientID(),
		},
		Text: text,
		SuggestedActions: &skypeapi.SuggestedActions{
			Actions: actions,
		},
	}
	return bot.sendActivity(&activity)
}

func (bot *Bot) MyConversations() (*skypeapi.ConversationsResult, error) {
	resp, err := bot.api.PlainRequest(http.MethodGet, "/v3/conversations?continuationToken=76579e23-9d24-4a8e-8530-04cd07a104f2", nil)
	if err != nil {
		return nil, err
	}

	var list skypeapi.ConversationsResult
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		return nil, err
	}

	return &list, nil
}

func (bot *Bot) lookupCommand(msg *Activity) bool {
	for event, handler := range bot.eventHandlers {
		switch cmd := event.(type) {
		case Command:
			if cmd.Match(msg.Text()) {
				handler(msg)
				return true
			}
		}
	}
	return false
}

func New(config Config) *Bot {
	logger := log.New(ioutil.Discard, "", 0)

	if config.Logger != nil {
		logger = config.Logger
	}

	return &Bot{
		api:           newAPI(config.AppID, config.AppSecret),
		eventHandlers: make(map[Event]func(*Activity), 0),
		logger:        logger,
	}
}
