# Go-Mailer

## Overview
Go-Mailer is a Go library that wraps the [wneessen/go-mail](https://github.com/wneessen/go-mail) library to provide a simplified interface for sending emails using an existing `mail.Client` instance. This library enhances `go-mail` with a fluent HTML message builder, making it easy to create and send rich, dynamic email content.

## Features
- **Uses Existing `mail.Client`**: Integrates with an existing `wneessen/go-mail` client for flexible configuration and reuse.
- **Fluent HTML Message Builder**: Provides a chainable API for constructing HTML emails with methods like `Subject`, `To`, `Line`, `Action`, etc.
- Attachment Support: Easily attach files from local disk, embedded filesystems, IOFS filesystems, or `io.Reader`/`io.ReadSeeker`.
- Customizable Email Content: Allows setting greetings, salutations, and product information for a personalized experience.
- Responsive HTML Template: Automatically formats emails with a clean, responsive design compatible with most email clients.

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
        mailer.WithFrom("noreply@example.com")
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

## Fluent HTML Message Builder
The fluent HTML message builder provides a chainable API to construct HTML emails intuitively. Key methods include:

- `Subject(string)`: Sets the email subject.
- `To(string)`: Adds a recipient email address.
- `Cc(string)`: Adds a CC recipient email address.
- `Bcc(string)`: Adds a BCC recipient email address.
- `Line(string)`: Adds a plain text line to the email body.
- `Linef(string, ...interface{})`: Adds a formatted line to the email body.
- `Action(string, string)`: Adds a clickable button with a label and URL.
- `Greeting(string)`: Sets a greeting line at the top of the email.
- `Salutation(string)`: Sets a salutation line at the bottom of the email.
- `Product(mailer.Product)`: Sets product information for the email, which can be used in the footer.
- `AttachFile(string, ...mailer.FileOption)`: Attaches a file from the local disk.
- `AttachFromEmbedFS(string, *embed.FS, ...mailer.FileOption)`: Attaches a file from an embedded filesystem.
- `AttachFromIOFS(string, fs.FS, ...mailer.FileOption)`: Attaches a file from an IOFS filesystem.
- `AttachReader(string, io.Reader, ...mailer.FileOption)`: Attaches a file from an `io.Reader`.
- `AttachReadSeeker(string, io.ReadSeeker, ...mailer.FileOption)`: Attaches a file from an `io.ReadSeeker`.

The builder automatically formats the email with a clean, responsive HTML template, ensuring compatibility with most email clients.

Example of a more complex email:

```go
message := mailer.NewMessage().
    Subject("Welcome to Our Platform").
    To("user@example.com").
    Line("Thank you for signing up!").
    Action("Activate Account", "https://example.com/activate").
    Line("We're excited to have you on board.")
```

## Configuration
Go-Mailer uses an existing `mail.Client` from `wneessen/go-mail`, which must be configured separately. The `mailer.New` function accepts optional configuration options for additional customization (e.g., default sender, reply-to address, etc.). Example configuration for the underlying `go-mail` client:

```go
client, err := mail.NewClient("smtp.sendgrid.net",
    mail.WithPort(587),
    mail.WithSMTPAuth(mail.SMTPAuthPlain),
    mail.WithUsername("apikey"),
    mail.WithPassword("your-sendgrid-api-key"),
)
if err != nil {
    panic(err)
}

m := mailer.New(client,
    mailer.WithFrom("noreply@example.com"),
    mailer.WithReplyTo("support@example.com"),
    mailer.WithProduct(mailer.Product{
        Name: "Go Mailer",
        URL: "https://github.com/ahmadfaizk/go-mailer",
    })
)
```

## Dependencies
Go-Mailer relies on the following key dependency:
- [wneessen/go-mail](https://github.com/wneessen/go-mail): Provides the core email sending functionality with support for SMTP, STARTTLS, and advanced features like DKIM and S/MIME.

## Contributing
Contributions are welcome! To contribute to Go-Mailer:

1. Fork the repository
2. Create a new branch (`git checkout -b feature/your-feature`)
3. Make your changes
4. Commit your changes (`git commit -m 'Add your feature'`)
5. Push to the branch (`git push origin feature/your-feature`)
6. Open a pull request

Please ensure your code adheres to the project's coding standards and includes tests.

<!-- ## License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details. -->

## Acknowledgements
- [wneessen/go-mail](https://github.com/wneessen/go-mail) for providing the robust foundation for this library.
- All contributors to the `go-mail` project for their excellent work.

## Contact
For questions or feedback, reach out to the maintainer at [ahmadfaizk's GitHub profile](https://github.com/ahmadfaizk).