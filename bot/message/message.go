package message

type Sendable interface {
	Send() error
}
