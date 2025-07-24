package gomailer

type config struct {
	fromName    string
	fromAddress string
	product     Product
}

// Option defines a function type that can be used to configure the Mailer.
type Option func(*config)

// WithFrom sets the sender's name and address for the email messages sent by the Mailer.
func WithFrom(name, address string) Option {
	return func(c *config) {
		c.fromName = name
		c.fromAddress = address
	}
}

// WithProduct sets the product information for the email messages sent by the Mailer.
func WithProduct(product Product) Option {
	return func(c *config) {
		c.product = product
	}
}

func newConfig(opts ...Option) *config {
	cfg := &config{
		fromName:    "GoMailer",
		fromAddress: "noreply@example.com",
		product: Product{
			Name: "GoMailer",
			URL:  "https://github.com/ahmadfaizk/go-mailer",
		},
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}
