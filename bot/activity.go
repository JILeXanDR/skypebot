package bot

import (
	"github.com/JILeXanDR/skypebot/skypeapi"
	"strings"
)

type Activity struct {
	activity *skypeapi.Activity
}

func (message *Activity) RecipientID() string {
	return message.Full().Conversation.ID
}

func (message *Activity) Full() *skypeapi.Activity {
	return message.activity
}

func (message *Activity) Text() string {
	if message.IsGroup() {
		// remove mention text "@botname " in the beginning of message
		return strings.Replace(message.activity.Text, message.activity.Recipient.Name+" ", "", 1)
	}
	return message.activity.Text
}

func (message *Activity) IsGroup() bool {
	return message.activity.Conversation.IsGroup
}

func (message *Activity) SomeoneWroteToMe() bool {
	return message.activity.Type == "message"
}

func (message *Activity) AddedToContacts() bool {
	return message.activity.Type == "contactRelationUpdate" && message.activity.Action == "add"
}

func (message *Activity) RemovedFromContacts() bool {
	return message.activity.Type == "contactRelationUpdate" && message.activity.Action == "remove"
}

func (message *Activity) AddedToConversation() bool {
	return message.activity.Type == "conversationUpdate" && len(message.activity.MembersAdded) > 0
}

func (message *Activity) RemovedFromConversation() bool {
	return message.activity.Type == "conversationUpdate" && len(message.activity.MembersRemoved) > 0
}

type Sender struct {
	account skypeapi.ChannelAccount
}

func (s *Sender) RecipientID() string {
	return s.account.ID
}

func (message *Activity) Sender() *Sender {
	return &Sender{
		account: message.activity.From,
	}
}

func (message *Activity) Command() Command {
	return ""
}
