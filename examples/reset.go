package main

import "github.com/ahmadfaizk/go-mailgen"

func resetMessage() *mailgen.Builder {
	return mailgen.New().
		Name("John Doe").
		Line("You have received this email because a password reset request for your account was received.").
		Line("Click the button below to reset your password:").
		Action("Reset your password", "https://example.com/reset-password").
		Line("If you did not request a password reset, no further action is required on your part.")
}
