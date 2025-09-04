package mailgen

// Message represents an email message with its components.
type Message interface {
	// Subject returns the subject of the email.
	Subject() string
	// From returns the sender's address.
	From() Address
	// FromString returns the sender's address as a formatted string.
	FromString() string
	// ReplyTo returns the Reply-To address.
	ReplyTo() *Address
	// ReplyToString returns the Reply-To address as a formatted string.
	ReplyToString() string
	// To returns the list of recipient addresses.
	To() []string
	// Cc returns the list of CC addresses.
	Cc() []string
	// Bcc returns the list of BCC addresses.
	Bcc() []string
	// HTML returns the HTML content of the email.
	HTML() string
	// PlainText returns the plain text content of the email.
	PlainText() string
}

// Address represents an email address with an optional name.
type Address struct {
	Name    string
	Address string
}

func (a Address) String() string {
	if a.Name == "" {
		return a.Address
	}
	return a.Name + " <" + a.Address + ">"
}

var _ Message = (*message)(nil)

type message struct {
	subject   string
	from      Address
	replyTo   *Address
	to        []string
	cc        []string
	bcc       []string
	html      string
	plainText string
}

func (m *message) Subject() string {
	return m.subject
}

func (m *message) From() Address {
	return m.from
}

func (m *message) FromString() string {
	return m.from.String()
}

func (m *message) ReplyTo() *Address {
	return m.replyTo
}

func (m *message) ReplyToString() string {
	if m.replyTo == nil {
		return ""
	}
	return m.replyTo.String()
}

func (m *message) To() []string {
	return m.to
}

func (m *message) Cc() []string {
	return m.cc
}

func (m *message) Bcc() []string {
	return m.bcc
}

func (m *message) HTML() string {
	return m.html
}

func (m *message) PlainText() string {
	return m.plainText
}
