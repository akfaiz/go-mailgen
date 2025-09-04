package mailgen

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/akfaiz/go-mailgen/templates"
	"github.com/vanng822/go-premailer/premailer"
)

// Product represents the product information used in the email.
type Product struct {
	Name      string
	Link      string
	Logo      string // Optional logo URL
	Copyright string
}

// Builder represents an email message with various fields such as subject, recipients, and content.
// It provides methods to set these fields and generate the HTML content for the email.
type Builder struct {
	subject string
	from    Address
	to      []string
	cc      []string
	bcc     []string

	textDirection  string
	theme          string
	preheader      string
	greeting       string
	name           string
	salutation     string
	components     []Component
	fallbacks      []*Action
	fallbackFormat string
	product        Product
}

var defaultBuilder atomic.Pointer[Builder]

func init() {
	defaultBuilder.Store(newDefaultBuilder())
}

func newDefaultBuilder() *Builder {
	return &Builder{
		textDirection: "ltr",
		theme:         "default",
		greeting:      "Hi",
		salutation:    "Best regards",
		product: Product{
			Name:      "Go-Mailgen",
			Link:      "https://github.com/akfaiz/go-mailgen",
			Copyright: fmt.Sprintf("© %d Go-Mailgen. All rights reserved.", time.Now().Year()),
		},
		fallbackFormat: "If you're having trouble clicking the \"[ACTION]\" button, copy and paste the URL below into your web browser:",
	}
}

func (b *Builder) clone() *Builder {
	return &Builder{
		textDirection:  b.textDirection,
		subject:        b.subject,
		from:           b.from,
		to:             append([]string{}, b.to...),
		cc:             append([]string{}, b.cc...),
		bcc:            append([]string{}, b.bcc...),
		theme:          b.theme,
		fallbackFormat: b.fallbackFormat,
		preheader:      b.preheader,
		greeting:       b.greeting,
		name:           b.name,
		salutation:     b.salutation,
		fallbacks:      append([]*Action{}, b.fallbacks...),
		components:     append([]Component{}, b.components...),
		product:        b.product,
	}
}

// SetDefault sets the default Builder instance.
//
// It can be useful for set global defaults or configurations for the email messages.
//
// Example usage:
//
//	mailgen.SetDefault(mailgen.New().
//		Product(mailgen.Product{
//			Name: "Go-Mailgen",
//			Link: "https://github.com/akfaiz/go-mailgen",
//			Logo: "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
//		}).
//		Theme("default"))
func SetDefault(b *Builder) {
	if b == nil {
		return
	}
	defaultBuilder.Store(b)
}

// New creates a new Message instance with default values for greeting, salutation, and product.
//
// Example usage:
//
//	message := mailgen.New().
//		Subject("Reset Password").
//		To("recipient@example.com").
//		Line("Click the button below to reset your password").
//		Action("Reset Password", "https://example.com/reset-password").
//		Line("If you did not request this, please ignore this email")
func New() *Builder {
	return defaultBuilder.Load().clone()
}

// Subject sets the subject of the email message.
func (b *Builder) Subject(subject string) *Builder {
	b.subject = subject
	return b
}

// From sets the sender's email address for the email message.
// It can include a name for the sender.
func (b *Builder) From(address string, name ...string) *Builder {
	addr := Address{
		Address: address,
	}
	if len(name) > 0 {
		addr.Name = name[0]
	}
	b.from = addr
	return b
}

// To add a recipient's email address to the email message.
func (b *Builder) To(to string, others ...string) *Builder {
	values := b.filterRecipients(to, others...)
	if len(values) == 0 {
		return b
	}
	b.to = append(b.to, values...)
	return b
}

func (b *Builder) filterRecipients(first string, others ...string) []string {
	if first == "" && len(others) == 0 {
		return nil
	}
	values := make([]string, 0, len(others)+1)
	values = append(values, first)
	values = append(values, others...)
	var filtered []string
	for _, recipient := range values {
		if recipient != "" {
			filtered = append(filtered, recipient)
		}
	}
	return filtered
}

// Cc adds a carbon copy (CC) recipient's email address to the email message.
func (b *Builder) Cc(cc string, others ...string) *Builder {
	values := b.filterRecipients(cc, others...)
	if len(values) == 0 {
		return b
	}
	b.cc = append(b.cc, values...)
	return b
}

// Bcc adds a blind carbon copy (BCC) recipient's email address to the email message.
func (b *Builder) Bcc(bcc string, others ...string) *Builder {
	values := b.filterRecipients(bcc, others...)
	if len(values) == 0 {
		return b
	}
	b.bcc = append(b.bcc, values...)
	return b
}

// Theme sets the theme for the email message.
// Supported themes are "default" and "plain".
func (b *Builder) Theme(theme string) *Builder {
	b.theme = theme
	return b
}

// TextDirection sets the text direction for the email message.
// It can be "ltr" (left-to-right) or "rtl" (right-to-left).
func (b *Builder) TextDirection(direction string) *Builder {
	if direction != "ltr" && direction != "rtl" {
		return b // Invalid direction, do nothing
	}
	b.textDirection = direction
	return b
}

// FallbackFormat sets the fallback format for action buttons in the email message.
// This format is used when the email client does not support HTML buttons.
//
// Example usage:
//
//	email := mailgen.New().
//		FallbackFormat("If you're having trouble clicking the \"[ACTION]\" button, copy and paste the URL below into your web browser:")
func (b *Builder) FallbackFormat(format string) *Builder {
	if format == "" {
		return b // No format provided, do nothing
	}
	b.fallbackFormat = format
	return b
}

// Preheader sets the preheader text for the email message.
// The preheader is a short summary text that follows the subject line when an email is viewed in the inbox.
// It is often used to provide additional context or a preview of the email content.
//
// Preheader is not displayed in the email body but is included in the email headers.
func (b *Builder) Preheader(preheader string) *Builder {
	b.preheader = preheader
	return b
}

// Greeting sets the greeting line of the email message.
// The default is "Hi".
func (b *Builder) Greeting(greeting string) *Builder {
	b.greeting = greeting
	return b
}

// Name sets the name of the greeting line in the email message.
// This is typically used to personalize the greeting with the recipient's name.
//
// If not set, the greeting will be a generic "Hi".
func (b *Builder) Name(name string) *Builder {
	b.name = name
	return b
}

// Salutation sets the closing salutation of the email message.
// Default is "Best regards".
func (b *Builder) Salutation(salutation string) *Builder {
	b.salutation = salutation
	return b
}

// Line adds a line of text to the email message.
// If an action is set, it will be added to the outro lines; otherwise, it will be added to the intro lines.
func (b *Builder) Line(text string) *Builder {
	b.components = append(b.components, Line{Text: text})
	return b
}

// Linef adds a formatted line of text to the email message.
// If an action is set, it will be added to the outro lines; otherwise, it will be added to the intro lines.
func (b *Builder) Linef(format string, args ...interface{}) *Builder {
	text := fmt.Sprintf(format, args...)
	return b.Line(text)
}

// Action sets the action text and link for the email message.
// It creates a button that the recipient can click to perform an action.
//
// Example usage:
//
//	email := mailgen.New().
//		Line("Click the button below to get started").
//		Action("Get Started", "https://example.com/get-started")
func (b *Builder) Action(text, link string, cfg ...Action) *Builder {
	action := &Action{
		Text:  text,
		Link:  link,
		Color: "#3869D4",
	}
	noFallback := false
	if len(cfg) > 0 {
		if cfg[0].Color != "" {
			action.Color = cfg[0].Color
		}
		noFallback = cfg[0].NoFallback
	}
	b.components = append(b.components, action)
	if !noFallback {
		b.fallbacks = append(b.fallbacks, action)
	}
	return b
}

// Product sets the product information for the email message.
func (b *Builder) Product(product Product) *Builder {
	defaultProduct := defaultBuilder.Load().product

	b.product = product
	if b.product.Name == "" {
		b.product.Name = defaultProduct.Name
	}
	if b.product.Copyright == "" {
		b.product.Copyright = fmt.Sprintf("© %d %s. All rights reserved.", time.Now().Year(), b.product.Name)
	}
	b.product.Link = product.Link
	return b
}

// Table sets a table to be included in the email message.
//
// Example usage:
//
//	email := mailgen.New().
//		Table(mailgen.Table{
//			Data: [][]mailgen.Entry{
//			{
//				{Key: "Name", Value: "John Doe"},
//				{Key: "Email", Value: "john.doe@example.com"},
//			},
//			{
//				{Key: "Name", Value: "Jane Smith"},
//				{Key: "Email", Value: "jane.smith@example.com"},
//			},
//		},
//		Columns: mailgen.Columns{
//			CustomWidth: map[string]string{
//				"Name":  "50%",
//				"Email": "50%",
//			},
//			CustomAlign: map[string]string{
//				"Name":  "left",
//				"Email": "right",
//			},
//		},
//	})
func (b *Builder) Table(table Table) *Builder {
	if len(table.Data) == 0 {
		return b // No table to add
	}

	b.components = append(b.components, &table)
	return b
}

// Build generates the final Message object with the HTML and plaintext content.
//
// It processes all the components, actions, and other fields set in the Builder.
//
// Returns an error if there is an issue generating the HTML or plaintext content.
func (b *Builder) Build() (Message, error) {
	b.beforeBuild()
	html, err := b.generateHTML()
	if err != nil {
		return nil, err
	}
	plainText, err := b.generatePlaintext()
	if err != nil {
		return nil, err
	}
	return &message{
		subject:   b.subject,
		from:      b.from,
		to:        b.to,
		cc:        b.cc,
		bcc:       b.bcc,
		html:      html,
		plainText: plainText,
	}, nil
}

func (b *Builder) beforeBuild() {
	for _, fallback := range b.fallbacks {
		fallback.FallbackText = strings.ReplaceAll(b.fallbackFormat, "[ACTION]", fallback.Text)
	}
}

type templateData struct {
	TextDirection  string
	Preheader      string
	Greeting       string
	Salutation     string
	ComponentsHTML []htmltemplate.HTML
	ComponentsText []string
	Fallbacks      []*Action
	Product        Product
}

func (b *Builder) htmlTemplate() *htmltemplate.Template {
	if b.theme == "plain" {
		return templates.PlainHtmlTmpl
	}
	return templates.DefaultHtmlTmpl
}

func (b *Builder) generateHTML() (string, error) {
	tmpl := b.htmlTemplate()

	var componentsHTML []htmltemplate.HTML
	for _, comp := range b.components {
		html, err := comp.HTML(tmpl)
		if err != nil {
			return "", err
		}
		componentsHTML = append(componentsHTML, htmltemplate.HTML(html))
	}

	data := templateData{
		TextDirection:  b.textDirection,
		Preheader:      b.preheader,
		Greeting:       b.greetingLine(),
		Salutation:     b.salutation,
		Product:        b.product,
		ComponentsHTML: componentsHTML,
		Fallbacks:      b.fallbacks,
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
	return cleanEmailHTML(html), nil
}

func cleanEmailHTML(input string) string {
	// Remove spaces and newlines between HTML tags
	reBetweenTags := regexp.MustCompile(`>\s+<`)
	clean := reBetweenTags.ReplaceAllString(input, "><")

	// Remove leading/trailing spaces on each line
	lines := strings.Split(clean, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	clean = strings.Join(lines, "\n")

	// Remove multiple empty lines
	reEmptyLines := regexp.MustCompile(`\n{2,}`)
	clean = reEmptyLines.ReplaceAllString(clean, "\n")

	// Final trim
	clean = strings.TrimSpace(clean)

	return clean
}

func (b *Builder) generatePlaintext() (string, error) {
	var componentsText []string
	for _, comp := range b.components {
		text, err := comp.PlainText()
		if err != nil {
			return "", err
		}
		componentsText = append(componentsText, text)
	}

	data := templateData{
		Greeting:       b.greetingLine(),
		Preheader:      b.preheader,
		Salutation:     b.salutation,
		Product:        b.product,
		ComponentsText: componentsText,
	}
	var buf bytes.Buffer
	if err := templates.DefaultPlainTextTmpl.ExecuteTemplate(&buf, "index.txt", data); err != nil {
		return "", err
	}
	text := buf.String()

	return cleanEmailText(text), nil
}

func cleanEmailText(input string) string {
	clean := strings.TrimSpace(input)
	re := regexp.MustCompile(`\n{3,}`)
	clean = re.ReplaceAllString(clean, "\n\n")
	return clean
}

func (b *Builder) greetingLine() string {
	if b.name != "" {
		if b.textDirection == "rtl" {
			return fmt.Sprintf("%s %s", b.name, b.greeting)
		}
		return fmt.Sprintf("%s %s", b.greeting, b.name)
	}
	if b.greeting == "" {
		return defaultBuilder.Load().greeting
	}
	return b.greeting
}
