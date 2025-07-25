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
- **Dynamic Components Ordering**: Add components like buttons, tables, and lines in any order you want.

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

## Supported Themes

The following open-source themes are bundled with this package:

- `default` by [Postmark Transactional Email Templates](https://github.com/ActiveCampaign/postmark-templates)

<img src="examples/default/welcome.png" height="200" /> <img src="examples/default/reset.png" height="200" /> <img src="examples/default/receipt.png" height="200" />

- `plain` by [Postmark Transactional Email Templates](https://github.com/ActiveCampaign/postmark-templates)

<img src="examples/plain/welcome.png" height="200" /> <img src="examples/plain/reset.png" height="200" /> <img src="examples/plain/receipt.png" height="200" />

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
