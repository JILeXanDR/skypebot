package message

import (
	"github.com/pkg/errors"
	"net/http"
)

type Builder struct {
	Text          *string
	AttachmentURL *string
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) WithText(text string) *Builder {
	b.Text = &text
	return b
}

func (b *Builder) WithAttachmentFromURL(url string) *Builder {
	b.AttachmentURL = &url
	return b
}

func (b *Builder) Build() (Sendable, error) {
	if b.AttachmentURL != nil {
		resp, err := http.Get(*b.AttachmentURL)
		if err != nil {
			return nil, errors.Wrapf(err, "can't download image via URL %s", *b.AttachmentURL)
		}
		return NewImageMessage(resp.Body), nil
	} else if b.Text != nil {
		return TextMessage(*b.Text), nil
	} else {
		return nil, errors.New("can't build message, any options was not passed")
	}
}
