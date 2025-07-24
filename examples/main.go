package main

import (
	"fmt"
	"os"

	"github.com/ahmadfaizk/go-mailer"
)

func main() {
	defaultProduct := mailer.Product{
		Name: "GoMailer",
		URL:  "https://github.com/ahmadfaizk/go-mailer",
	}
	messages := map[string]*mailer.Message{
		"reset_password": resetPassword(),
		"welcome":        welcomeEmail(),
		"receipt":        receiptEmail(),
	}
	themes := []string{"default", "plain"}

	for _, theme := range themes {
		for name, msg := range messages {
			msg.Product(defaultProduct)
			msg.Theme(theme)
			html, err := msg.GenerateHTML()
			if err != nil {
				panic(err)
			}
			plainText, err := msg.GeneratePlaintext()
			if err != nil {
				panic(err)
			}

			htmlFileName := fmt.Sprintf("examples/%s/%s.html", theme, name)
			plainTextFileName := fmt.Sprintf("examples/%s/%s.txt", theme, name)

			if err := os.WriteFile(htmlFileName, []byte(html), 0644); err != nil {
				panic(fmt.Sprintf("failed to write HTML file %s: %v", htmlFileName, err))
			}
			if err := os.WriteFile(plainTextFileName, []byte(plainText), 0644); err != nil {
				panic(fmt.Sprintf("failed to write plaintext file %s: %v", plainTextFileName, err))
			}
		}
	}
}

func resetPassword() *mailer.Message {
	return mailer.NewMessage().
		Subject("Reset your password").
		To("recipient@example.com").
		Preheader("Use this link to reset your password. The link is only valid for 24 hours.").
		Line("Click the button below to reset your password").
		Action("Reset your password", "https://example.com/reset-password").
		Line("If you did not request this, please ignore this email")
}

func welcomeEmail() *mailer.Message {
	return mailer.NewMessage().
		Subject("Welcome to our service!").
		To("recipient@example.com").
		Line("Thank you for signing up for our service!").
		Line("We're glad to have you on board.").
		Line("If you have any questions, feel free to reach out to our support team.")
}

func receiptEmail() *mailer.Message {
	return mailer.NewMessage().
		Subject("Your order receipt").
		To("recipient@example.com").
		Line("Thank you for your order!").
		Table(mailer.Table{
			Instruction: "Here are the details of your order:",
			Headers: []mailer.TableHeader{
				{Text: "Item", Align: "left", Width: "50%"},
				{Text: "Count", Align: "right", Width: "25%"},
				{Text: "Price", Align: "right", Width: "25%"},
			},
			Rows: [][]string{
				{"Widget A", "2", "$20.00"},
				{"Widget B", "1", "$15.00"},
			},
		}).
		Action("View Order", "https://example.com/order", mailer.Action{Instruction: "Click here to view your order"}).
		Line("If you have any questions, feel free to contact us.")
}
