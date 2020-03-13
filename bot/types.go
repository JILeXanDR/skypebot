package bot

import (
	"fmt"
	"strconv"
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

func NewCommand(name string, args []string) *Command {
	trimmed := strings.ReplaceAll(name, " ", "")
	if len(trimmed) != len(name) || trimmed == "" {
		panic(fmt.Sprintf(`command name "%s" is wrong, can't contain spaces or be an empty string`, name))
	}
	return &Command{
		name:       name,
		args:       args,
		parsedArgs: make(map[string]interface{}, 0),
	}
}

type Command struct {
	name       string
	args       []string
	parsedArgs map[string]interface{}
}

func (c *Command) ID() string {
	return fmt.Sprintf("command:%s", c.Name())
}

func (c *Command) Name() string {
	return c.name
}

func (c *Command) Args() map[string]interface{} {
	return c.parsedArgs
}

// Match checks does "message" is command.
func (c *Command) Match(message string) bool {
	values := strings.Split(message, " ")
	return values[0] == c.Name()
}

func (c *Command) Parse(text string) {
	values := strings.Split(text, " ")
	args := make(map[string]interface{}, 0)
	if len(values) > 1 {
		messageValues := values[1:]
		for i, arg := range c.args {
			val := messageValues[i]
			if int, err := strconv.Atoi(val); err == nil {
				args[arg] = int
			} else {
				args[arg] = val
			}
		}
	}
	c.parsedArgs = args
}
