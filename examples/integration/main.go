package main

import (
	"github.com/afkdevs/go-mailgen"
	"github.com/wneessen/go-mail"
)

func main() {
	// Create a new mailer instance
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
			Link: "https://github.com/afkdevs/go-mailgen",
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
