package mailer

import (
	"bytes"
	"embed"
	"fmt"
	htmltemplate "html/template"
	"io"
	"io/fs"
	"regexp"
	"strings"
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

// Message represents an email message with various fields such as subject, recipients, and content.
// It provides methods to set these fields and generate the HTML content for the email.
type Message struct {
	subject string
	from    *Address
	replyTo string
	to      []string
	cc      []string
	bcc     []string

	theme      string
	preheader  string
	greeting   string
	salutation string
	actions    []*Action
	components []Component
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
	m := &Message{
		theme:      "default",
		greeting:   "Hello",
		salutation: "Best regards",
	}

	return m.Product(Product{
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

// Theme sets the theme for the email message.
// Supported themes are "default" and "plain".
func (m *Message) Theme(theme string) *Message {
	m.theme = theme
	return m
}

// Preheader sets the preheader text for the email message.
// The preheader is a short summary text that follows the subject line when an email is viewed in the inbox.
// It is often used to provide additional context or a preview of the email content.
//
// Preheader is not displayed in the email body but is included in the email headers.
func (m *Message) Preheader(preheader string) *Message {
	m.preheader = preheader
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
	m.components = append(m.components, &Line{Text: text})
	return m
}

// Linef adds a formatted line of text to the email message.
// If an action is set, it will be added to the outro lines; otherwise, it will be added to the intro lines.
func (m *Message) Linef(format string, args ...interface{}) *Message {
	text := fmt.Sprintf(format, args...)
	return m.Line(text)
}

// Action sets the action text and URL for the email message.
// It creates a button in the email that links to the specified URL.
// The action can also include an optional instruction and color for the button.
func (m *Message) Action(text, url string, act ...Action) *Message {
	action := Action{
		Text:  text,
		URL:   url,
		Color: "#3869D4",
	}
	if len(act) > 0 {
		if act[0].Color != "" {
			action.Color = act[0].Color
		}
	}
	m.actions = append(m.actions, &action)
	m.components = append(m.components, &action)
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

// Table sets a table to be included in the email message.
//
// Example usage:
//
//	message.Table(mailer.Table{
//		Headers: []mailer.TableHeader{
//			{Text: "Item", Align: "left", Width: "70%"},
//			{Text: "Price", Align: "right", Width: "30%"},
//		},
//		Rows: [][]string{
//			{"Widget A", "$10.00"},
//			{"Widget B", "$15.00"},
//		},
//	})
func (m *Message) Table(table Table) *Message {
	// Ensure headers have default values for width and alignment
	for i, header := range table.Headers {
		if header.Width == "" {
			table.Headers[i].Width = "auto"
		}
		if header.Align == "" {
			table.Headers[i].Align = "left"
		}
	}
	if table.Rows == nil {
		table.Rows = [][]string{}
	}
	if len(table.Rows) == 0 && len(table.Headers) == 0 {
		return m // No table to add
	}

	m.components = append(m.components, &table)
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
	m.filesEmbedFS = append(m.filesEmbedFS, fileEmbedFS{
		file: file{name: name, cfg: newFileConfig(opts...)},
		fs:   fs,
	})
	return m
}

// AttachFromIOFS attaches a file from an IOFS filesystem to the email message.
// The file is specified by its name and the IOFS filesystem, along with additional options for configuration.
func (m *Message) AttachFromIOFS(name string, fs fs.FS, opts ...FileOption) *Message {
	m.filesIOFS = append(m.filesIOFS, fileIOFS{
		file: file{name: name, cfg: newFileConfig(opts...)},
		FS:   fs,
	})
	return m
}

// AttachReader attaches a file from an io.Reader to the email message.
// The file is specified by its name and the reader, along with additional options for configuration.
func (m *Message) AttachReader(name string, reader io.Reader, opts ...FileOption) *Message {
	m.filesReader = append(m.filesReader, fileReader{
		file:   file{name: name, cfg: newFileConfig(opts...)},
		Reader: reader,
	})
	return m
}

// AttachReadSeeker attaches a file from an io.ReadSeeker to the email message.
// The file is specified by its name and the read seeker, along with additional options for configuration.
func (m *Message) AttachReadSeeker(name string, readSeeker io.ReadSeeker, opts ...FileOption) *Message {
	m.filesReadSeeker = append(m.filesReadSeeker, fileReadSeeker{
		file:       file{name: name, cfg: newFileConfig(opts...)},
		ReadSeeker: readSeeker,
	})
	return m
}

type templateData struct {
	Theme          string
	Preheader      string
	Greeting       string
	Salutation     string
	ComponentsHTML []htmltemplate.HTML
	ComponentsText []string
	Actions        []*Action // Used for sub-copy in HTML
	Product        Product
}

//go:embed templates/default/*
var defaultTmplFS embed.FS

//go:embed templates/plain/*
var plainTmplFS embed.FS

var defaultHtmlTmpl = htmltemplate.Must(htmltemplate.New("index.html").
	ParseFS(defaultTmplFS, "templates/default/*.html"),
)
var defaultPlainTextTmpl = texttemplate.Must(texttemplate.New("index.txt").
	ParseFS(defaultTmplFS, "templates/default/*.txt"),
)
var plainHtmlTmpl = htmltemplate.Must(htmltemplate.New("index.html").
	ParseFS(plainTmplFS, "templates/plain/*.html"),
)
var plainPlainTextTmpl = texttemplate.Must(texttemplate.New("index.txt").
	ParseFS(plainTmplFS, "templates/plain/*.txt"),
)

func (m *Message) htmlTemplate() *htmltemplate.Template {
	if m.theme == "plain" {
		return plainHtmlTmpl
	}
	return defaultHtmlTmpl
}

func (m *Message) plainTextTemplate() *texttemplate.Template {
	if m.theme == "plain" {
		return plainPlainTextTmpl
	}
	return defaultPlainTextTmpl
}

// GenerateHTML generates the HTML content for the email message using the provided template.
func (m *Message) GenerateHTML() (string, error) {
	tmpl := m.htmlTemplate()

	var componentsHTML []htmltemplate.HTML
	for _, comp := range m.components {
		html, err := comp.HTML(tmpl)
		if err != nil {
			return "", err
		}
		componentsHTML = append(componentsHTML, htmltemplate.HTML(html))
	}

	data := templateData{
		Theme:          m.theme,
		Preheader:      m.preheader,
		Greeting:       m.greeting,
		Salutation:     m.salutation,
		Product:        m.product,
		ComponentsHTML: componentsHTML,
		Actions:        m.actions,
	}
	var buf bytes.Buffer

	if err := tmpl.ExecuteTemplate(&buf, "index.html", data); err != nil {
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
	tmpl := m.plainTextTemplate()
	var componentsText []string
	for _, comp := range m.components {
		text := comp.PlainText()
		componentsText = append(componentsText, text)
	}

	data := templateData{
		Greeting:       boxString(fmt.Sprintf("%s,", m.greeting)),
		Preheader:      m.preheader,
		Salutation:     m.salutation,
		Product:        m.product,
		ComponentsText: componentsText,
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "index.txt", data); err != nil {
		return "", err
	}
	text := buf.String()
	return m.cleanPlainText(text), nil
}

func (m *Message) cleanPlainText(text string) string {
	text = strings.TrimSpace(text)
	re := regexp.MustCompile(`\n{3,}`)
	text = re.ReplaceAllString(text, "\n\n")

	return text
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

func boxString(s string) string {
	// Find the max line length (in case of multi-line input)
	lines := strings.Split(s, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// The border should match the longest line length
	border := strings.Repeat("*", maxLen)

	// Combine
	var b strings.Builder
	b.WriteString(border + "\n")
	b.WriteString(s + "\n")
	b.WriteString(border)

	return b.String()
}
