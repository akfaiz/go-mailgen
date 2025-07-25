package mailgen

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"regexp"
	"strings"
	"sync/atomic"
	texttemplate "text/template"
	"time"

	"github.com/ahmadfaizk/go-mailgen/templates"
	"github.com/vanng822/go-premailer/premailer"
)

// Product represents the product information used in the email.
// It includes the product name, logo URL, product URL, and copyright information.
type Product struct {
	Name      string
	URL       string
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

	theme      string
	preheader  string
	greeting   string
	salutation string
	actions    []Action
	components []Component
	product    Product
}

var defaultBuilder atomic.Pointer[Builder]

func init() {
	defaultBuilder.Store(newdefaultBuilder())
}

func newdefaultBuilder() *Builder {
	builder := &Builder{}
	return builder.Theme("default").
		Greeting("Hello").
		Salutation("Best regards").
		Product(Product{
			Name: "GoMailgen",
			URL:  "https://github.com/ahmadfaizk/go-mailgen",
		})
}

func (b *Builder) clone() *Builder {
	return &Builder{
		subject:    b.subject,
		from:       b.from,
		to:         append([]string{}, b.to...),
		cc:         append([]string{}, b.cc...),
		bcc:        append([]string{}, b.bcc...),
		theme:      b.theme,
		preheader:  b.preheader,
		greeting:   b.greeting,
		salutation: b.salutation,
		actions:    append([]Action{}, b.actions...),
		components: append([]Component{}, b.components...),
		product:    b.product,
	}
}

// SetDefault sets the default Builder instance.
//
// It can be useful for set global defaults or configurations for the email messages.
func SetDefault(b *Builder) {
	if b == nil {
		return
	}
	defaultBuilder.Store(b)
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

// To adds a recipient's email address to the email message.
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

func (b *Builder) CC(cc string, others ...string) *Builder {
	values := b.filterRecipients(cc, others...)
	if len(values) == 0 {
		return b
	}
	b.cc = append(b.cc, values...)
	return b
}

func (b *Builder) BCC(bcc string, others ...string) *Builder {
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
// Default is "Hello".
func (b *Builder) Greeting(greeting string) *Builder {
	b.greeting = greeting
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

// Action sets the action text and URL for the email message.
// It creates a button in the email that links to the specified URL.
// The action can also include an optional instruction and color for the button.
func (b *Builder) Action(text, url string, color ...string) *Builder {
	action := Action{
		Text:  text,
		URL:   url,
		Color: "#3869D4",
	}
	if len(color) > 0 && color[0] != "" {
		action.Color = color[0]
	}
	b.actions = append(b.actions, action)
	b.components = append(b.components, action)
	return b
}

// Product sets the product information for the email message.
func (b *Builder) Product(product Product) *Builder {
	b.product = product
	if b.product.Name == "" {
		b.product.Name = "GoMailgen"
	}
	if b.product.URL == "" {
		b.product.URL = "#"
	}
	if b.product.Copyright == "" {
		b.product.Copyright = fmt.Sprintf("© %d %s. All rights reserved.", time.Now().Year(), b.product.Name)
	}
	return b
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
func (b *Builder) Table(table Table) *Builder {
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
		return b // No table to add
	}

	b.components = append(b.components, table)
	return b
}

// Build generates the final Message object with the HTML and plaintext content.
// It processes all the components, actions, and other fields set in the Builder.
// Returns an error if there is an issue generating the HTML or plaintext content.
func (b *Builder) Build() (Message, error) {
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

type templateData struct {
	Theme          string
	Preheader      string
	Greeting       string
	Salutation     string
	ComponentsHTML []htmltemplate.HTML
	ComponentsText []string
	Actions        []Action // Used for sub-copy in HTML
	Product        Product
}

func (b *Builder) htmlTemplate() *htmltemplate.Template {
	if b.theme == "plain" {
		return templates.PlainHtmlTmpl
	}
	return templates.DefaultHtmlTmpl
}

func (b *Builder) plainTextTemplate() *texttemplate.Template {
	if b.theme == "plain" {
		return templates.PlainPlainTextTmpl
	}
	return templates.DefaultPlainTextTmpl
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
		Theme:          b.theme,
		Preheader:      b.preheader,
		Greeting:       b.greeting,
		Salutation:     b.salutation,
		Product:        b.product,
		ComponentsHTML: componentsHTML,
		Actions:        b.actions,
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

func (b *Builder) generatePlaintext() (string, error) {
	tmpl := b.plainTextTemplate()
	var componentsText []string
	for _, comp := range b.components {
		text := comp.PlainText()
		componentsText = append(componentsText, text)
	}

	data := templateData{
		Greeting:       boxString(fmt.Sprintf("%s,", b.greeting)),
		Preheader:      b.preheader,
		Salutation:     b.salutation,
		Product:        b.product,
		ComponentsText: componentsText,
	}
	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "index.txt", data); err != nil {
		return "", err
	}
	text := buf.String()

	text = strings.TrimSpace(text)
	re := regexp.MustCompile(`\n{3,}`)
	text = re.ReplaceAllString(text, "\n\n")

	return text, nil
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
