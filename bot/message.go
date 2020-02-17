package bot

import "io"

type Replier interface {
	Reply(Sendable) error
}

type Sendable interface {
	Send() error
}

type TextMessage struct {
	Text string
}

func (m *TextMessage) Send() error {
	return nil
}

func NewTextMessage(text string) Sendable {
	return &TextMessage{Text: text}
}

// TextReplier sends text message
type MessageReply func(message Sendable) error

// MessageHandler todo
type MessageHandler func(message *IncomingActivity, reply MessageReply) error

type ImageMessage struct {
	image io.Reader
}

func (m *ImageMessage) Send() error {
	return nil
}

func NewImageMessage(image io.Reader) Sendable {
	return &ImageMessage{
		image: image,
	}
}
