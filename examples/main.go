package main

import (
	"fmt"
	"os"

	"github.com/akfaiz/go-mailgen"
)

func main() {
	mailgen.SetDefault(mailgen.New().
		Product(
			mailgen.Product{
				Name: "Go-Mailgen",
				Link: "https://github.com/akfaiz/go-mailgen",
				Logo: "https://upload.wikimedia.org/wikipedia/commons/0/05/Go_Logo_Blue.svg",
			},
		))
	messageBuilders := map[string]*mailgen.Builder{
		"reset":   resetMessage(),
		"welcome": welcomeMessage(),
		"receipt": receiptMessage(),
	}
	themes := []string{"default", "plain"}

	for _, theme := range themes {
		for name, builder := range messageBuilders {
			builder.Theme(theme)

			msg, err := builder.Build()
			if err != nil {
				panic(fmt.Sprintf("failed to build message %s: %v", name, err))
			}

			htmlFileName := fmt.Sprintf("examples/%s/%s.html", theme, name)
			plainTextFileName := fmt.Sprintf("examples/%s/%s.txt", theme, name)

			if err := os.WriteFile(htmlFileName, []byte(msg.HTML()), 0644); err != nil {
				panic(fmt.Sprintf("failed to write HTML file %s: %v", htmlFileName, err))
			}
			if err := os.WriteFile(plainTextFileName, []byte(msg.PlainText()), 0644); err != nil {
				panic(fmt.Sprintf("failed to write plaintext file %s: %v", plainTextFileName, err))
			}
		}
	}
}

func receiptMessage() *mailgen.Builder {
	return mailgen.New().
		Name("John Doe").
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
		Line("You can check the status of your order and more in your dashboard:").
		Action("Go to Dashboard", "https://example.com/dashboard").
		Line("We thank you for your purchase.")
}

func resetMessage() *mailgen.Builder {
	return mailgen.New().
		Name("John Doe").
		Line("You have received this email because a password reset request for your account was received.").
		Line("Click the button below to reset your password:").
		Action("Reset your password", "https://example.com/reset-password").
		Line("If you did not request a password reset, no further action is required on your part.")
}

func welcomeMessage() *mailgen.Builder {
	return mailgen.New().
		Name("John Doe").
		Line("Welcome to Go-Mailgen! We're very excited to have you on board.").
		Line("To get started with Mailgen, please click here:").
		Action("Get Started", "https://example.com/get-started").
		Line("We're glad to have you on board.").
		Line("Need help, or have questions? Just reply to this email, we'd love to help.")
}
