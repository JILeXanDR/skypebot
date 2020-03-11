package bot

import (
	"fmt"
	"strings"
)

type event string

type Action interface {
	ID() string
}

const (
	OnTextMessage             event = "text_message"
	OnAttachment              event = "attachment"
	OnAddedToContacts         event = "added_to_contacts"
	OnRemovedFromContacts     event = "removed_from_contacts"
	OnAddedToConversation     event = "added_to_conversation"
	OnRemovedFromConversation event = "removed_from_conversation"
	OnAll                     event = "all"
)

func (e event) ID() string {
	return fmt.Sprintf("event:%s", e)
}

type Recipienter interface {
	RecipientID() string
}

type ConversationID string

func (id ConversationID) RecipientID() string {
	return string(id)
}

type Command string

func (c Command) Name() string {
	values := strings.Split(string(c), " ")
	return values[0]
}

func (c Command) Args(activity *Activity) map[string]interface{} {
	return map[string]interface{}{}
}

func (c Command) ID() string {
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
