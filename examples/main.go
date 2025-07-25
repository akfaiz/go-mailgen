package main

import (
	"fmt"
	"os"

	"github.com/ahmadfaizk/go-mailgen"
)

func main() {
	defaultProduct := mailgen.Product{
		Name: "GoMailer",
		Link: "https://github.com/ahmadfaizk/go-mailer",
	}
	messageBuilders := map[string]*mailgen.Builder{
		"reset":   resetMessage(),
		"welcome": welcomeMessage(),
		"receipt": receiptMessage(),
	}
	themes := []string{"default", "plain"}

	for _, theme := range themes {
		for name, builder := range messageBuilders {
			builder.Product(defaultProduct)
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

func resetMessage() *mailgen.Builder {
	return mailgen.New().
		Preheader("Use this link to reset your password. The link is only valid for 24 hours.").
		Line("Click the button below to reset your password").
		Action("Reset your password", "https://example.com/reset-password").
		Line("If you did not request this, please ignore this email")
}

func welcomeMessage() *mailgen.Builder {
	return mailgen.New().
		Line("Thank you for signing up for our service!").
		Line("We're glad to have you on board.").
		Line("If you have any questions, feel free to reach out to our support team.")
}

func receiptMessage() *mailgen.Builder {
	return mailgen.New().
		Line("Thank you for your order!").
		Line("Here are the details of your purchase:").
		Table(mailgen.Table{
			Data: [][]mailgen.Entry{
				{{Key: "Item", Value: "Widget A"}, {Key: "Count", Value: "2"}, {Key: "Price", Value: "$20.00"}},
				{{Key: "Item", Value: "Widget B"}, {Key: "Count", Value: "1"}, {Key: "Price", Value: "$15.00"}},
			},
			Columns: mailgen.Columns{
				CustomWidth: map[string]string{
					"Item":  "50%",
					"Count": "25%",
					"Price": "25%",
				},
				CustomAlign: map[string]string{
					"Item":  "left",
					"Count": "right",
					"Price": "right",
				},
			},
		}).
		Line("Click the button below to view your order details.").
		Action("View Order", "https://example.com/order").
		Line("If you have any questions, feel free to contact us.")
}
