package mailer

import "fmt"

// Address represents an email address with an optional name.
type Address struct {
	Name    string
	Address string
}

// String returns a string representation of the email message.
func (a Address) String() string {
	if a.Name == "" {
		return a.Address
	}
	return fmt.Sprintf("%s <%s>", a.Name, a.Address)
}

type config struct {
	theme   string
	from    Address
	replyTo string
	product Product
}

// Option defines a function type that can be used to configure the Mailer.
type Option func(*config)

// WithFrom sets the sender's name and address for the email messages sent by the Mailer.
func WithFrom(address string, name ...string) Option {
	return func(c *config) {
		c.from.Address = address
		if len(name) > 0 {
			c.from.Name = name[0]
		}
	}
}

// WithProduct sets the product information for the email messages sent by the Mailer.
func WithProduct(product Product) Option {
	return func(c *config) {
		c.product = product
	}
}

// WithReplyTo sets the reply-to address for the email messages sent by the Mailer.
func WithReplyTo(replyTo string) Option {
	return func(c *config) {
		c.replyTo = replyTo
	}
}

// WithTheme sets the theme for the email messages sent by the Mailer.
// Supported themes are "default" and "plain".
func WithTheme(theme string) Option {
	return func(c *config) {
		c.theme = theme
	}
}

func newConfig(opts ...Option) *config {
	cfg := &config{
		product: Product{
			Name: "GoMailer",
			URL:  "https://github.com/ahmadfaizk/go-mailer",
		},
		theme: "default",
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
