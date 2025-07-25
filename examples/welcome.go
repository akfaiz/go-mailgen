package main

import "github.com/ahmadfaizk/go-mailgen"

func welcomeMessage() *mailgen.Builder {
	return mailgen.New().
		Name("John Doe").
		Line("Welcome to Go-Mailgen! We're very excited to have you on board.").
		Line("To get started with Mailgen, please click here:").
		Action("Get Started", "https://example.com/get-started").
		Line("We're glad to have you on board.").
		Line("Need help, or have questions? Just reply to this email, we'd love to help.")
}
