package mailer

import (
	"context"
	"errors"

	"github.com/wneessen/go-mail"
)

// Mailer is responsible for sending email messages using the provided mail client and configuration.
type Mailer struct {
	client *mail.Client
	cfg    *config
}

// New creates a new Mailer instance with the provided mail client and configuration options.
func New(client *mail.Client, opts ...Option) *Mailer {
	cfg := newConfig(opts...)
	return &Mailer{client: client, cfg: cfg}
}

// Send sends the email message using the Mailer instance.
func (m *Mailer) Send(message *Message) error {
	return m.SendContext(context.Background(), message)
}

// SendContext sends the email message using the Mailer instance with a provided context.
func (m *Mailer) SendContext(ctx context.Context, message *Message) error {
	msg, err := m.toMailMsg(message)
	if err != nil {
		return err
	}
	return m.client.DialAndSendWithContext(ctx, msg)
}

func (m *Mailer) toMailMsg(message *Message) (*mail.Msg, error) {
	if message == nil {
		return nil, errors.New("message cannot be nil")
	}
	if len(message.GetTo()) == 0 {
		return nil, errors.New("message must have at least one recipient")
	}
	if message.GetSubject() == "" {
		return nil, errors.New("message must have a subject")
	}

	message.Product(m.cfg.product)
	html, err := message.GenerateHTML()
	if err != nil {
		return nil, err
	}

	msg := mail.NewMsg()

	from := m.cfg.from
	if message.GetFrom() != nil {
		from = *message.GetFrom()
	}
	if err := msg.From(from.String()); err != nil {
		return nil, err
	}

	if err := msg.To(message.GetTo()...); err != nil {
		return nil, err
	}
	if len(message.GetCc()) > 0 {
		if err := msg.Cc(message.GetCc()...); err != nil {
			return nil, err
		}
	}
	if len(message.GetBcc()) > 0 {
		if err := msg.Bcc(message.GetBcc()...); err != nil {
			return nil, err
		}
	}
	replyTo := message.GetReplyTo()
	if replyTo == "" {
		replyTo = m.cfg.replyTo
	}
	if replyTo != "" {
		if err := msg.ReplyTo(replyTo); err != nil {
			return nil, err
		}
	}
	msg.Subject(message.GetSubject())
	msg.SetBodyString(mail.TypeTextHTML, html)

	if err := m.setAttachments(msg, message); err != nil {
		return nil, err
	}

	return msg, nil
}

func (m *Mailer) setAttachments(msg *mail.Msg, message *Message) error {
	for _, f := range message.files {
		msg.AttachFile(f.name, f.cfg.toMailFileOption()...)
	}

	for _, f := range message.filesEmbedFS {
		if err := msg.AttachFromEmbedFS(f.name, f.fs, f.cfg.toMailFileOption()...); err != nil {
			return err
		}
	}

	for _, f := range message.filesIOFS {
		if err := msg.AttachFromIOFS(f.name, f.FS, f.cfg.toMailFileOption()...); err != nil {
			return err
		}
	}

	for _, f := range message.filesReader {
		if err := msg.AttachReader(f.name, f.Reader, f.cfg.toMailFileOption()...); err != nil {
			return err
		}
	}

	for _, f := range message.filesReadSeeker {
		msg.AttachReadSeeker(f.name, f.ReadSeeker, f.cfg.toMailFileOption()...)
	}

	return nil
}
