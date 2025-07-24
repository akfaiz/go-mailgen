package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"time"
)

// Product represents the product information used in the email.
// It includes the product name, logo URL, product URL, and copyright information.
type Product struct {
	Name      string
	LogoURL   string
	URL       string
	Copyright string
}

// Message represents an email message with various fields such as subject, recipients, and content.
// It provides methods to set these fields and generate the HTML content for the email.
type Message struct {
	subject string
	to      []string
	cc      []string
	bcc     []string

	greeting   string
	salutation string
	introLines []string
	outroLines []string
	actionText string
	actionURL  string
	product    Product
}

// NewMessage creates a new Message instance with default values for greeting, salutation, and product.
//
// Example usage:
//
//	html, err := gomailer.NewMessage().
//		Subject("Reset Password").
//		To("recipient@example.com").
//		Line("Click the button below to reset your password").
//		Action("Reset Password", "https://example.com/reset-password").
//		Line("If you did not request this, please ignore this email").
//		GenerateHTML()
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
	if m.actionText == "" {
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
func (m *Message) Action(text, url string) *Message {
	m.actionText = text
	m.actionURL = url
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

type templateData struct {
	Greeting         string
	Salutation       string
	IntroLines       []string
	OutroLines       []string
	ActionText       string
	ActionURL        string
	ProductName      string
	ProductURL       string
	ProductLogo      string
	ProductCopyright string
}

//go:embed templates/default/*
var defaultTmplFS embed.FS

var defaultTmpl = template.Must(template.New("message.html").ParseFS(defaultTmplFS, "templates/default/*.html"))

// GenerateHTML generates the HTML content for the email message using the provided template.
func (m *Message) GenerateHTML() (string, error) {
	data := templateData{
		Greeting:         m.greeting,
		Salutation:       m.salutation,
		IntroLines:       m.introLines,
		OutroLines:       m.outroLines,
		ActionText:       m.actionText,
		ActionURL:        m.actionURL,
		ProductName:      m.product.Name,
		ProductURL:       m.product.URL,
		ProductLogo:      m.product.LogoURL,
		ProductCopyright: m.product.Copyright,
	}
	var buf bytes.Buffer
	if err := defaultTmpl.ExecuteTemplate(&buf, "message.html", data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetSubject returns the subject of the email message.
func (m *Message) GetSubject() string {
	return m.subject
}

// GetTo returns the recipient(s) of the email message.
func (m *Message) GetTo() []string {
	return m.to
}

// GetCc returns the carbon copy (CC) recipients of the email message.
func (m *Message) GetCc() []string {
	return m.cc
}

// GetBcc returns the blind carbon copy (BCC) recipients of the email message.
func (m *Message) GetBcc() []string {
	return m.bcc
}
