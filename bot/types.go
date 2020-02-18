package bot

type Event uint8

const (
	EventMessage Event = iota
	EventAddedToContacts
	EventRemovedFromContacts
	EventAddedToConversation
	EventRemovedFromConversation
	EventAll
)

func (e Event) String() string {
	m := map[Event]string{
		EventMessage:                 "message",
		EventAddedToContacts:         "added_to_contacts",
		EventRemovedFromContacts:     "removed_from_contacts",
		EventAddedToConversation:     "added_to_conversation",
		EventRemovedFromConversation: "removed_from_conversation",
		EventAll:                     "*",
	}
	return m[e]
}

type Recipient interface {
	ConversationID() string
}

type ConversationID string

func (id ConversationID) ConversationID() string {
	return string(id)
}
