package message

import "io"

type ImageMessage struct {
	Reader io.Reader
}

func (m *ImageMessage) Send() error {
	return nil
}

func NewImageMessage(image io.Reader) Sendable {
	return &ImageMessage{
		Reader: image,
	}
}
