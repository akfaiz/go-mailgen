package main

import "github.com/ahmadfaizk/go-mailgen"

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
