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

// any person writes message to bot (private chat or in group chat)
func (message *Activity) SomeoneWroteToMe() bool {
	return message.activity.Type == "message"
}

// bot added to contacts
func (message *Activity) AddedToContacts() bool {
	return message.activity.Type == "contactRelationUpdate" && message.activity.Action == "add"
}

// bot removed from contacts
func (message *Activity) RemovedFromContacts() bool {
	return message.activity.Type == "contactRelationUpdate" && message.activity.Action == "remove"
}

// any person added to the conversation (NOT A BOT)
func (message *Activity) AddedToConversation() bool {
	return message.activity.Type == "conversationUpdate" && len(message.activity.MembersAdded) > 0
}

// any person removed from the conversation (NOT A BOT)
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
