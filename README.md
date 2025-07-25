# Go-Mailgen

[![Go](https://github.com/ahmadfaizk/go-mailgen/actions/workflows/ci.yml/badge.svg)](https://github.com/ahmadfaizk/go-mailgen/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ahmadfaizk/go-mailgen)](https://goreportcard.com/report/github.com/ahmadfaizk/go-mailgen)
[![codecov](https://codecov.io/gh/ahmadfaizk/go-mailgen/graph/badge.svg?token=7tbSVRaD4b)](https://codecov.io/gh/ahmadfaizk/go-mailgen)
[![GoDoc](https://pkg.go.dev/badge/github.com/ahmadfaizk/go-mailgen)](https://pkg.go.dev/github.com/ahmadfaizk/go-mailgen)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ahmadfaizk/go-mailgen)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

**Go-Mailgen** is a Go library for generating professional HTML emails using a fluent, intuitive API. Simplify email creation with customizable templates and seamless integration into your Go applications. This project is inspired by the [mailgen](https://github.com/eladnava/mailgen) Node.js package, bringing its elegant email generation approach to the Go ecosystem.

## Features

- **Fluent API**: Build emails with a clean, chainable interface.
- **Inline CSS**: Ensures compatibility across major email clients.
- **Template-Based**: Use pre-built or custom templates for rapid development.
- **Easy Integration**: Works effortlessly with popular Go mail libraries like go-mail.
- **Chainable Methods**: Add content, actions, and tables in a straightforward manner.
- **Global Configuration**: Set default sender, product information, and theme for all emails.

## Installation

To install Go-Mailgen, run the following command:

```bash
go get github.com/ahmadfaizk/go-mailgen
```

## Usage

Here's a simple example of how to use Go-Mailgen to create an email:

```go
package main

import (
	"github.com/ahmadfaizk/go-mailgen"
	"github.com/wneessen/go-mail"
)

func main() {
	// Initialize SMTP client
	mailer, err := mail.NewClient("smtp.example.com",
		mail.WithPort(587),
		mail.WithUsername("user"),
		mail.WithPassword("pass"),
	)
	if err != nil {
		panic(err)
	}

	// Set global configuration (optional)
	mailgen.SetDefault(mailgen.New().
		From("no-reply@example.com", "Go-Mailgen").
		Product(mailgen.Product{
			Name: "Go-Mailgen",
			URL:  "https://github.com/ahmadfaizk/go-mailgen",
		}).
		Theme("default"),
	)

	// Build the email
	email := mailgen.New().
		Subject("Reset Your Password").
		To("johndoe@mail.com").
		Line("Click the button below to reset your password").
		Action("Reset your password", "https://example.com/reset-password").
		Line("If you did not request this, please ignore this email")
	message, err := email.Build()
	if err != nil {
		panic(err)
	}

	// Send the email
	msg := mail.NewMsg()
	msg.Subject(message.Subject())
	msg.From(message.From().String())
	msg.To(message.To()...)
	msg.SetBodyString(mail.TypeTextPlain, message.PlainText())
	msg.SetBodyString(mail.TypeTextHTML, message.HTML())
	if err := mailer.DialAndSend(msg); err != nil {
		panic(err)
	}
}
```

## More Examples

You can find more examples in the [examples](examples) directory.

## Documentation
For detailed documentation, please visit the [Go-Mailgen documentation](https://pkg.go.dev/github.com/ahmadfaizk/go-mailgen).

## Supported Themes

The following open-source themes are bundled with this package:

- `default` by [Postmark Transactional Email Templates](https://github.com/ActiveCampaign/postmark-templates)

<img src="examples/default/welcome.png" height="200" /> <img src="examples/default/reset.png" height="200" /> <img src="examples/default/receipt.png" height="200" />

- `plain` by [Postmark Transactional Email Templates](https://github.com/ActiveCampaign/postmark-templates)

<img src="examples/plain/welcome.png" height="200" /> <img src="examples/plain/reset.png" height="200" /> <img src="examples/plain/receipt.png" height="200" />

## Elements

Go-Mailgen provides several methods to add content to your emails. Here are some of the most commonly used methods:

### Action

To add an action button to your email, use the `Action` method:

```go
mailgen.New().
	Line("To confirm your email address, please click the button below:").
	Action("Confirm Email", "https://example.com/confirm")
```

To add multiple actions to your email, you can chain the `Action` method:

```go
mailgen.New().
	Line("To confirm your email address, please click the buttons below:").
	Action("Confirm Email", "https://example.com/confirm").
	Line("Or you can visit our website:").
	Action("Visit Website", "https://example.com")
```

### Table

To add a table to your email, use the `Table` method:

```go
mailgen.New().
	Line("Your order has been processed successfully.").
	Table(mailgen.Table{
		Data: [][]mailgen.Entry{
			{
				{Key: "Item", Value: "Golang"},
				{Key: "Description", Value: "An open-source programming language supported by Google."},
				{Key: "Price", Value: "$10.99"},
			},
			{
				{Key: "Item", Value: "Mailgen"},
				{Key: "Description", Value: "Programmatically create beautiful e-mails using Golang"},
				{Key: "Price", Value: "$1.99"},
			},
		},
		Columns: mailgen.Columns{
			CustomWidth: map[string]string{
				"Item":  "20%",
				"Price": "15%",
			},
			CustomAlign: map[string]string{
				"Price": "right",
			},
		},
	})
```

You can chain the `Table` with other methods to add more content to your email:

```go
mailgen.New().
	Line("Your order has been processed successfully.").
	Table(mailgen.Table{
		Data: [][]mailgen.Entry{
			{
				{Key: "Item", Value: "Golang"},
				{Key: "Description", Value: "An open-source programming language supported by Google."},
				{Key: "Price", Value: "$10.99"},
			},
			{
				{Key: "Item", Value: "Mailgen"},
				{Key: "Description", Value: "Programmatically create beautiful e-mails using Golang"},
				{Key: "Price", Value: "$1.99"},
			},
		},
		Columns: mailgen.Columns{
			CustomWidth: map[string]string{
				"Item":  "20%",
				"Price": "15%",
			},
			CustomAlign: map[string]string{
				"Price": "right",
			},
		},
	}).
	Line("If you have any questions, feel free to contact us.").
	Action("Contact Support", "https://example.com/contact").
	Line("Thank you for your order!")
```

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
