package message

type TextMessage string

func (m TextMessage) Send() error {
	return nil
}
