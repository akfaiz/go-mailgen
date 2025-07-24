# Go-Mailer
[![Go](https://github.com/ahmadfaizk/go-mailer/actions/workflows/ci.yml/badge.svg)](https://github.com/ahmadfaizk/go-mailer/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ahmadfaizk/go-mailer)](https://goreportcard.com/report/github.com/ahmadfaizk/go-mailer)
[![codecov](https://codecov.io/gh/ahmadfaizk/go-mailer/graph/badge.svg?token=7tbSVRaD4b)](https://codecov.io/gh/ahmadfaizk/go-mailer)
[![GoDoc](https://pkg.go.dev/badge/github.com/ahmadfaizk/go-mailer)](https://pkg.go.dev/github.com/ahmadfaizk/go-mailer)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ahmadfaizk/go-mailer)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

Go-Mailer is a Go library that wraps the [wneessen/go-mail](https://github.com/wneessen/go-mail) library to provide a simplified interface for sending emails using an existing `mail.Client` instance. This library enhances `go-mail` with a fluent HTML message builder, making it easy to create and send rich, dynamic email content.

## Features
- **Uses Existing `mail.Client`**: Integrates with an existing `wneessen/go-mail` client for flexible configuration and reuse.
- **Fluent HTML Message Builder**: Provides a chainable API for constructing HTML emails with methods like `Subject`, `To`, `Line`, `Action`, etc.
- Attachment Support: Easily attach files from local disk, embedded filesystems, IOFS filesystems, or `io.Reader`/`io.ReadSeeker`.
- Responsive HTML Template: Automatically formats emails with a clean, responsive design compatible with most email clients.
- CSS-inlined HTML: Automatically inlines CSS styles for better compatibility with email clients that strip out `<style>` tags.

## Installation
To install Go-Mailer, run the following command:

```bash
go get github.com/ahmadfaizk/go-mailer
```

This will also install the `wneessen/go-mail` dependency.

## Usage
Go-Mailer requires an existing `mail.Client` instance from `wneessen/go-mail` to send emails. Below is an example demonstrating how to initialize the mailer and send an HTML email using the fluent message builder:

```go
package main

import (
    "context"
    "github.com/ahmadfaizk/go-mailer"
    "github.com/wneessen/go-mail"
)

func main() {
    // Create a new go-mail client
    client, err := mail.NewClient("smtp.example.com",
        mail.WithPort(587),
        mail.WithSMTPAuth(mail.SMTPAuthPlain),
        mail.WithUsername("your-username"),
        mail.WithPassword("your-password"),
    )
    if err != nil {
        panic(err)
    }

    // Initialize the Go-Mailer instance with the existing client
    m := mailer.New(client,
        mailer.WithFrom("noreply@example.com"),
        mailer.WithProduct(mailer.Product{
            Name: "My Application",
            URL:  "https://example.com",
        })
    )

    // Build an HTML email using the fluent message builder
    message := mailer.NewMessage().
        Subject("Reset Password").
        To("recipient@example.com").
        Line("Click the button below to reset your password").
        Action("Reset Password", "https://example.com/reset-password").
        Line("If you did not request this, please ignore this email")

    // Send the email
    err = m.Send(message)
    if err != nil {
        panic(err)
    }
}
```

## More Examples

You can find more examples in the [examples](examples) directory.

## Supported Themes

The following open-source themes are bundled with this package:

* `default` by [Postmark Transactional Email Templates](https://github.com/ActiveCampaign/postmark-templates)

<img src="examples/default/welcome.png" height="200" /> <img src="examples/default/reset.png" height="200" /> <img src="examples/default/receipt.png" height="200" />

* `plain` by [Postmark Transactional Email Templates](https://github.com/ActiveCampaign/postmark-templates)

<img src="examples/plain/welcome.png" height="200" /> <img src="examples/plain/reset.png" height="200" /> <img src="examples/plain/receipt.png" height="200" />

## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.