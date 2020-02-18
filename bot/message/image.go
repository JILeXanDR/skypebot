package message

import "io"

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

