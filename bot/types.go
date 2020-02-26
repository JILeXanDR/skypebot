package bot

import (
	"fmt"
	"strings"
)

type event string

type Event interface {
	EventID() string
}

const (
	EventMessage                 event = "event:message"
	EventAddedToContacts         event = "event:added_to_contacts"
	EventRemovedFromContacts     event = "event:removed_from_contacts"
	EventAddedToConversation     event = "event:added_to_conversation"
	EventRemovedFromConversation event = "event:removed_from_conversation"
	EventAll                     event = "event:all"
)

func (e event) EventID() string {
	return string(e)
}

type Recipienter interface {
	RecipientID() string
}

type ConversationID string

func (id ConversationID) RecipientID() string {
	return string(id)
}

// get_updates :klan
type Command string

func (c Command) Name() string {
	values := strings.Split(string(c), " ")
	return values[0]
}

func (c Command) Args(activity *Activity) map[string]interface{} {
	return map[string]interface{}{}
}

func (c Command) EventID() string {
	return fmt.Sprintf("command:%s", c.Name())
}

func (c Command) Match(message string) bool {
	values := strings.Split(message, " ")
	return values[0] == c.Name()
}

func (c Command) Parse(activity *Activity) Cmd {
	return Cmd{}
}

type Cmd struct {
	Text string
	Name string
	Args map[string]interface{}
}
