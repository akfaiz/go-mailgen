package mailgen

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"maps"
	"regexp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/ahmadfaizk/go-mailgen/component"
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

type Table struct {
	Data    [][]Entry
	Columns Columns
}

type Entry struct {
	Key   string
	Value string
}

type Columns struct {
	CustomWidth map[string]string
	CustomAlign map[string]string
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
	actions    []*component.Action
	components []component.Component
	product    Product
}

var defaultBuilder atomic.Pointer[Builder]

func init() {
	defaultBuilder.Store(newDefaultBuilder())
}

func newDefaultBuilder() *Builder {
	return &Builder{
		theme:      "default",
		greeting:   "Hello",
		salutation: "Best regards",
		product: Product{
			Name:      "Go-Mailgen",
			URL:       "https://github.com/ahmadfaizk/go-mailgen",
			Copyright: fmt.Sprintf("© %d Go-Mailgen. All rights reserved.", time.Now().Year()),
		},
	}
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
		actions:    append([]*component.Action{}, b.actions...),
		components: append([]component.Component{}, b.components...),
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

// New creates a new Message instance with default values for greeting, salutation, and product.
//
// Example usage:
//
//	message := mailer.NewMessage().
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
// The default is "Hello".
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
	b.components = append(b.components, component.Line{Text: text})
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
	action := &component.Action{
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
	defaultProduct := defaultBuilder.Load().product

	b.product = product
	if b.product.Name == "" {
		b.product.Name = defaultProduct.Name
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
	if len(table.Data) == 0 {
		return b // No table to add
	}
	// headers := make([]string, 0, len(table.Data[0]))
	// for _, entry := range table.Data[0] {
	// headers = append(headers, entry.Key)

	// width, ok := table.Columns.CustomWidth[entry.Key]
	// if !ok || width == "" {
	// 	width = "auto"
	// }
	// table.Columns.CustomWidth[entry.Key] = width

	// align, ok := table.Columns.CustomAlign[entry.Key]
	// if !ok || align == "" {
	// 	align = "left"
	// }
	// table.Columns.CustomAlign[entry.Key] = align
	// }

	b.components = append(b.components, table.component())
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
	Actions        []*component.Action // Used for sub-copy in HTML
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
	tmpl := templates.PlainTextTmpl
	var componentsText []string
	for _, comp := range b.components {
		text, err := comp.PlainText(tmpl)
		if err != nil {
			return "", err
		}
		componentsText = append(componentsText, text)
	}

	data := templateData{
		Greeting:       b.greeting,
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

func (t Table) component() component.Component {
	if len(t.Data) == 0 {
		return nil
	}
	tableComponent := component.Table{
		Data: make([][]component.Entry, len(t.Data)),
		Columns: component.Columns{
			CustomWidth: make(map[string]string),
			CustomAlign: make(map[string]string),
		},
	}
	for i, row := range t.Data {
		tableComponent.Data[i] = make([]component.Entry, len(row))
		for j, entry := range row {
			tableComponent.Data[i][j] = component.Entry{
				Key:   entry.Key,
				Value: entry.Value,
			}
		}
	}
	maps.Copy(tableComponent.Columns.CustomWidth, t.Columns.CustomWidth)
	maps.Copy(tableComponent.Columns.CustomAlign, t.Columns.CustomAlign)

	return &tableComponent
}
