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

// Address represents an email address with an optional name.
type Address struct {
	Name    string
	Address string
}

// Action represents an action button in the email message.
type Action struct {
	Instruction string
	Text        string
	URL         string
	Color       string
}

// Product represents the product information used in the email.
// It includes the product name, logo URL, product URL, and copyright information.
type Product struct {
	Name      string
	URL       string
	Copyright string
}

// Table represents a simple table structure for the email message.
type Table struct {
	Instruction string
	Headers     []TableHeader
	Rows        [][]string
}

// TableHeader represents a header in the table.
type TableHeader struct {
	Text  string
	Width string
	Align string
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
	introLines []string
	outroLines []string
	action     *Action
	product    Product
	table      *Table

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
func (m *Message) Action(text, url string, act ...Action) *Message {
	action := Action{
		Text:  text,
		URL:   url,
		Color: "primary",
	}
	if len(act) > 0 {
		action.Instruction = act[0].Instruction
		if act[0].Color != "" {
			action.Color = act[0].Color
		}
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

	m.table = &table
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
	Theme       string
	Preheader   string
	Greeting    string
	Salutation  string
	IntroLines  []string
	OutroLines  []string
	Action      *Action
	Product     Product
	Table       *Table
	TableString string // For plain text representation of the table
}

//go:embed templates/default/*
var defaultTmplFS embed.FS

//go:embed templates/plain/*
var plainTmplFS embed.FS

var defaultHtmlTmpl = htmltemplate.Must(htmltemplate.New("message.html").
	ParseFS(defaultTmplFS, "templates/default/*.html"),
)
var defaultPlainTextTmpl = texttemplate.Must(texttemplate.New("message.txt").
	ParseFS(defaultTmplFS, "templates/default/*.txt"),
)
var plainHtmlTmpl = htmltemplate.Must(htmltemplate.New("message.html").
	ParseFS(plainTmplFS, "templates/plain/*.html"),
)
var plainPlainTextTmpl = texttemplate.Must(texttemplate.New("message.txt").
	ParseFS(plainTmplFS, "templates/plain/*.txt"),
)

// GenerateHTML generates the HTML content for the email message using the provided template.
func (m *Message) GenerateHTML() (string, error) {
	data := templateData{
		Theme:      m.theme,
		Preheader:  m.preheader,
		Greeting:   m.greeting,
		Salutation: m.salutation,
		IntroLines: m.introLines,
		OutroLines: m.outroLines,
		Action:     m.action,
		Product:    m.product,
		Table:      m.table,
	}
	var buf bytes.Buffer
	var tmpl *htmltemplate.Template
	if m.theme == "plain" {
		tmpl = plainHtmlTmpl
	} else {
		tmpl = defaultHtmlTmpl
	}
	if err := tmpl.ExecuteTemplate(&buf, "message.html", data); err != nil {
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
		Preheader:  m.preheader,
		Salutation: m.salutation,
		IntroLines: m.introLines,
		OutroLines: m.outroLines,
		Action:     m.action,
		Product:    m.product,
	}
	if m.table != nil {
		data.TableString = m.table.String()
	}
	var buf bytes.Buffer
	var tmpl *texttemplate.Template
	if m.theme == "plain" {
		tmpl = plainPlainTextTmpl
	} else {
		tmpl = defaultPlainTextTmpl
	}
	if err := tmpl.ExecuteTemplate(&buf, "message.txt", data); err != nil {
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

// String returns a string representation of the email message.
func (a Address) String() string {
	if a.Name == "" {
		return a.Address
	}
	return fmt.Sprintf("%s <%s>", a.Name, a.Address)
}

// String returns a string representation of the table.
func (t Table) String() string {
	var sb strings.Builder

	// Add optional instruction
	if t.Instruction != "" {
		sb.WriteString(t.Instruction + "\n\n")
	}

	// Determine column widths
	colWidths := make([]int, len(t.Headers))
	for i, header := range t.Headers {
		colWidths[i] = len(header.Text)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Write headers
	for i, header := range t.Headers {
		sb.WriteString(t.padCell(header.Text, colWidths[i], header.Align))
		if i < len(t.Headers)-1 {
			sb.WriteString(" | ")
		}
	}
	sb.WriteString("\n")

	// Write separator line
	for i, width := range colWidths {
		sb.WriteString(strings.Repeat("-", width))
		if i < len(colWidths)-1 {
			sb.WriteString("-+-")
		}
	}
	sb.WriteString("\n")

	// Write rows
	for _, row := range t.Rows {
		for i, cell := range row {
			align := t.Headers[i].Align
			sb.WriteString(t.padCell(cell, colWidths[i], align))
			if i < len(row)-1 {
				sb.WriteString(" | ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (t Table) padCell(text string, width int, align string) string {
	switch strings.ToLower(align) {
	case "right":
		return fmt.Sprintf("%*s", width, text)
	case "center":
		padding := width - len(text)
		left := padding / 2
		right := padding - left
		return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
	default: // left
		return fmt.Sprintf("%-*s", width, text)
	}
}
