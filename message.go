package mailer

import (
	"bytes"
	"embed"
	"fmt"
	htmltemplate "html/template"
	"io"
	"io/fs"
	texttemplate "text/template"
	"time"

	"github.com/vanng822/go-premailer/premailer"
)

// Product represents the product information used in the email.
// It includes the product name, logo URL, product URL, and copyright information.
type Product struct {
	Name      string
	URL       string
	Copyright string
}

// Action represents an action button in the email message.
type Action struct {
	Text  string
	URL   string
	Color string
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
	return fmt.Sprintf("%s <%s>", a.Name, a.Address)
}

// Message represents an email message with various fields such as subject, recipients, and content.
// It provides methods to set these fields and generate the HTML content for the email.
type Message struct {
	subject string
	from    *Address
	replyTo string
	to      []string
	cc      []string
	bcc     []string

	greeting   string
	salutation string
	introLines []string
	outroLines []string
	action     *Action
	product    Product

	files           []file
	filesEmbedFS    []fileEmbedFS
	filesIOFS       []fileIOFS
	filesReader     []fileReader
	filesReadSeeker []fileReadSeeker
}

// NewMessage creates a new Message instance with default values for greeting, salutation, and product.
//
// Example usage:
//
//	messsage := mailer.NewMessage().
//		Subject("Reset Password").
//		To("recipient@example.com").
//		Line("Click the button below to reset your password").
//		Action("Reset Password", "https://example.com/reset-password").
//		Line("If you did not request this, please ignore this email")
func NewMessage() *Message {
	m := &Message{}

	return m.Greeting("Hello").
		Salutation("Best regards").
		Product(Product{
			Name: "GoMailer",
			URL:  "https://github.com/ahmadfaizk/go-mailer",
		})
}

// Subject sets the subject of the email message.
func (m *Message) Subject(subject string) *Message {
	m.subject = subject
	return m
}

// From sets the sender's address of the email message.
func (m *Message) From(address string, name ...string) *Message {
	if m.from == nil {
		m.from = &Address{}
	}
	m.from.Address = address
	if len(name) > 0 {
		m.from.Name = name[0]
	}
	return m
}

// ReplyTo sets the reply-to address for the email message.
func (m *Message) ReplyTo(replyTo string) *Message {
	m.replyTo = replyTo
	return m
}

// To sets the recipient(s) of the email message.
// It can accept multiple email addresses.
func (m *Message) To(to ...string) *Message {
	m.to = append(m.to, to...)
	return m
}

// Cc adds carbon copy (CC) recipients to the email message.
func (m *Message) Cc(cc ...string) *Message {
	m.cc = append(m.cc, cc...)
	return m
}

// Bcc adds blind carbon copy (BCC) recipients to the email message.
func (m *Message) Bcc(bcc ...string) *Message {
	m.bcc = append(m.bcc, bcc...)
	return m
}

// Greeting sets the greeting line of the email message.
// Default is "Hello".
func (m *Message) Greeting(greeting string) *Message {
	m.greeting = greeting
	return m
}

// Salutation sets the closing salutation of the email message.
// Default is "Best regards".
func (m *Message) Salutation(salutation string) *Message {
	m.salutation = salutation
	return m
}

// Line adds a line of text to the email message.
// If an action is set, it will be added to the outro lines; otherwise, it will be added to the intro lines.
func (m *Message) Line(text string) *Message {
	if m.action == nil {
		m.introLines = append(m.introLines, text)
	} else {
		m.outroLines = append(m.outroLines, text)
	}
	return m
}

// Linef adds a formatted line of text to the email message.
// If an action is set, it will be added to the outro lines; otherwise, it will be added to the intro lines.
func (m *Message) Linef(format string, args ...interface{}) *Message {
	text := fmt.Sprintf(format, args...)
	return m.Line(text)
}

// Action sets the action text and URL for the email message.
// It can also accept a color for the button, defaulting to "primary" if not provided.
//
// Supporting colors are: primary, green, and red.
func (m *Message) Action(text, url string, color ...string) *Message {
	action := Action{
		Text:  text,
		URL:   url,
		Color: "primary",
	}
	if len(color) > 0 {
		action.Color = color[0]
	}
	m.action = &action
	return m
}

// Product sets the product information for the email message.
func (m *Message) Product(product Product) *Message {
	m.product = product
	if m.product.Name == "" {
		m.product.Name = "GoMailer"
	}
	if m.product.URL == "" {
		m.product.URL = "#"
	}
	if m.product.Copyright == "" {
		m.product.Copyright = fmt.Sprintf("© %d %s. All rights reserved.", time.Now().Year(), m.product.Name)
	}
	return m
}

// AttachFile attaches a file to the email message.
// The file is specified by its name and can include additional options for configuration.
func (m *Message) AttachFile(name string, opts ...FileOption) *Message {
	m.files = append(m.files, file{name: name, cfg: newFileConfig(opts...)})
	return m
}

// AttachFromEmbedFS attaches a file from an embedded filesystem to the email message.
// The file is specified by its name and the embedded filesystem, along with additional options for configuration
func (m *Message) AttachFromEmbedFS(name string, fs *embed.FS, opts ...FileOption) *Message {
	m.filesEmbedFS = append(m.filesEmbedFS, fileEmbedFS{file{name: name, cfg: newFileConfig(opts...)}, fs})
	return m
}

// AttachFromIOFS attaches a file from an IOFS filesystem to the email message.
// The file is specified by its name and the IOFS filesystem, along with additional options for configuration.
func (m *Message) AttachFromIOFS(name string, fs fs.FS, opts ...FileOption) *Message {
	m.filesIOFS = append(m.filesIOFS, fileIOFS{file: file{name: name, cfg: newFileConfig(opts...)}, FS: fs})
	return m
}

// AttachReader attaches a file from an io.Reader to the email message.
// The file is specified by its name and the reader, along with additional options for configuration.
func (m *Message) AttachReader(name string, reader io.Reader, opts ...FileOption) *Message {
	m.filesReader = append(m.filesReader, fileReader{file: file{name: name, cfg: newFileConfig(opts...)}, Reader: reader})
	return m
}

// AttachReadSeeker attaches a file from an io.ReadSeeker to the email message.
// The file is specified by its name and the read seeker, along with additional options for configuration.
func (m *Message) AttachReadSeeker(name string, readSeeker io.ReadSeeker, opts ...FileOption) *Message {
	m.filesReadSeeker = append(m.filesReadSeeker, fileReadSeeker{file: file{name: name, cfg: newFileConfig(opts...)}, ReadSeeker: readSeeker})
	return m
}

type templateData struct {
	Greeting   string
	Salutation string
	IntroLines []string
	OutroLines []string
	Action     *Action
	Product    Product
}

//go:embed templates/default/*
var defaultTemplateFS embed.FS

//go:embed templates/plaintext/*
var plaintextTemplateFS embed.FS

var defaultTmpl = htmltemplate.Must(htmltemplate.New("message.html").ParseFS(defaultTemplateFS, "templates/default/*.html"))
var plaintextTmpl = texttemplate.Must(texttemplate.New("message.txt").ParseFS(plaintextTemplateFS, "templates/plaintext/*.txt"))

// GenerateHTML generates the HTML content for the email message using the provided template.
func (m *Message) GenerateHTML() (string, error) {
	data := templateData{
		Greeting:   m.greeting,
		Salutation: m.salutation,
		IntroLines: m.introLines,
		OutroLines: m.outroLines,
		Action:     m.action,
		Product:    m.product,
	}
	var buf bytes.Buffer
	if err := defaultTmpl.ExecuteTemplate(&buf, "message.html", data); err != nil {
		return "", err
	}
	prem, err := premailer.NewPremailerFromBytes(buf.Bytes(), premailer.NewOptions())
	if err != nil {
		return "", err
	}
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}
	return html, nil
}

// GeneratePlaintext generates the plaintext content for the email message using the provided template.
func (m *Message) GeneratePlaintext() (string, error) {
	data := templateData{
		Greeting:   m.greeting,
		Salutation: m.salutation,
		IntroLines: m.introLines,
		OutroLines: m.outroLines,
		Action:     m.action,
		Product:    m.product,
	}
	var buf bytes.Buffer
	if err := plaintextTmpl.ExecuteTemplate(&buf, "message.txt", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetSubject returns the subject of the email message.
func (m *Message) GetSubject() string {
	return m.subject
}

// GetFrom returns the sender's address of the email message.
func (m *Message) GetFrom() *Address {
	return m.from
}

// GetTo returns the recipient(s) of the email message.
func (m *Message) GetTo() []string {
	return m.to
}

// GetReplyTo returns the reply-to address of the email message.
func (m *Message) GetReplyTo() string {
	return m.replyTo
}

// GetCc returns the carbon copy (CC) recipients of the email message.
func (m *Message) GetCc() []string {
	return m.cc
}

// GetBcc returns the blind carbon copy (BCC) recipients of the email message.
func (m *Message) GetBcc() []string {
	return m.bcc
}
