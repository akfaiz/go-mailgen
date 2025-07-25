package main

import (
	"fmt"
	"os"

	"github.com/afkdevs/go-mailgen"
)

func main() {
	mailgen.SetDefault(mailgen.New().
		Product(
			mailgen.Product{
				Name: "Go-Mailgen",
				Link: "https://github.com/afkdevs/go-mailgen",
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
