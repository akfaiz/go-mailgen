package mailgen

type Message interface {
	Subject() string
	From() Address
	To() []string
	CC() []string
	BCC() []string
	HTML() string
	PlainText() string
}

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

func (m *message) To() []string {
	return m.to
}

func (m *message) CC() []string {
	return m.cc
}

func (m *message) BCC() []string {
	return m.bcc
}

func (m *message) HTML() string {
	return m.html
}

func (m *message) PlainText() string {
	return m.plainText
}
