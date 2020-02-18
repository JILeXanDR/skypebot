package bot

import (
	"github.com/JILeXanDR/skypebot/skypeapi"
	"strings"
)

type IncomingActivity struct {
	activity *skypeapi.Activity
}

func (message *IncomingActivity) Full() *skypeapi.Activity {
	return message.activity
}

func (message *IncomingActivity) Text() string {
	if message.IsGroup() {
		return strings.Replace(message.activity.Text, message.activity.Recipient.Name+" ", "", 1)
	}
	return message.activity.Text
}

func (message *IncomingActivity) IsGroup() bool {
	return message.activity.Conversation.IsGroup
}

func (message *IncomingActivity) SomeoneWroteToMe() bool {
	return message.activity.Type == "message"
}

func (message *IncomingActivity) AddedToContacts() bool {
	return message.activity.Type == "contactRelationUpdate" && message.activity.Action == "add"
}

func (message *IncomingActivity) RemovedFromContacts() bool {
	return message.activity.Type == "contactRelationUpdate" && message.activity.Action == "remove"
}

func (message *IncomingActivity) AddedToConversation() bool {
	return message.activity.Type == "conversationUpdate" && len(message.activity.MembersAdded) > 0
}

func (message *IncomingActivity) RemovedFromConversation() bool {
	return message.activity.Type == "conversationUpdate" && len(message.activity.MembersRemoved) > 0
}

func (message *IncomingActivity) FromUser() skypeapi.ChannelAccount {
	return message.activity.From
}
