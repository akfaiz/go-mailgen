package mailgen_test

import (
	"testing"
	"time"

	"github.com/ahmadfaizk/go-mailgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMessage_Line(t *testing.T) {
	msg := mailgen.NewMessage()

	msg.Line("First line")
	msg.Line("Second line")

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "First line")
	assert.Contains(t, html, "Second line")
}

func TestMessage_Linef(t *testing.T) {
	msg := mailgen.NewMessage()

	msg.Linef("Hello %s, you have %d messages", "John", 5)

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "Hello John, you have 5 messages")
}

func TestMessage_Action(t *testing.T) {
	msg := mailgen.NewMessage()

	msg.Action("Click Here", "https://example.com")

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "Click Here")
	assert.Contains(t, html, "https://example.com")
}

func TestMessage_Product(t *testing.T) {
	msg := mailgen.NewMessage()

	product := mailgen.Product{
		Name:      "Test Product",
		URL:       "https://example.com",
		Copyright: "© 2023 Test Product",
	}

	msg.Product(product)

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "Test Product")
	assert.Contains(t, html, "© 2023 Test Product")
}

func TestMessage_ProductDefaults(t *testing.T) {
	msg := mailgen.NewMessage()

	msg.Product(mailgen.Product{})

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "GoMailgen")

	currentYear := time.Now().Year()
	expectedCopyright := "© " + string(rune(currentYear+'0'))
	assert.Contains(t, html, expectedCopyright[:3])
}

func TestMessage_ChainedMethods(t *testing.T) {
	msg := mailgen.NewMessage().
		Greeting("Hi there").
		Line("This is a test").
		Action("Test Action", "https://test.com").
		Line("After action line")

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "Hi there")
}

func TestMessage_GenerateHTML(t *testing.T) {
	msg := mailgen.NewMessage().
		Line("You are receiving this email because we received a password reset request for your account.").
		Action("Reset Password", "https://example.com/reset-password").
		Linef("This password reset link will expire in %d minutes.", 60).
		Line("If you did not request a password reset, no further action is required.")

	html, err := msg.GenerateHTML()
	require.NoError(t, err)
	plainHtml, err := msg.Theme("plain").GenerateHTML()
	require.NoError(t, err)

	assert.NotEmpty(t, html)
	assert.Contains(t, html, "<html")
	assert.Contains(t, html, "You are receiving this email because we received a password reset request for your account.")
	assert.NotEmpty(t, plainHtml)
	assert.Contains(t, plainHtml, "<html")
	assert.Contains(t, plainHtml, "You are receiving this email because we received a password reset request for your account.")
}

func TestMessage_GenerateHTMLWithTable(t *testing.T) {
	msg := mailgen.NewMessage().
		Line("Thank you for your order!").
		Table(mailgen.Table{
			Headers: []mailgen.TableHeader{
				{Text: "Item", Align: "left", Width: "70%"},
				{Text: "Price", Align: "right", Width: "30%"},
			},
			Rows: [][]string{
				{"Widget A", "$10.00"},
				{"Widget B", "$15.00"},
			},
		}).
		Action("View Order", "https://example.com/order").
		Line("If you have any questions, feel free to contact us.")

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "<table")
	assert.Contains(t, html, "Widget A")
	assert.Contains(t, html, "$10.00")
}

func TestMessage_GeneratePlaintext(t *testing.T) {
	msg := mailgen.NewMessage().
		Line("You are receiving this email because we received a password reset request for your account.").
		Action("Reset Password", "https://example.com/reset-password").
		Linef("This password reset link will expire in %d minutes.", 60).
		Line("If you did not request a password reset, no further action is required.")

	plaintext, err := msg.GeneratePlaintext()
	require.NoError(t, err)

	assert.NotEmpty(t, plaintext)
	assert.Contains(t, plaintext, "You are receiving this email because we received a password reset request for your account.")
}

func TestMessage_GeneratePlaintextWithTable(t *testing.T) {
	msg := mailgen.NewMessage().
		Line("Thank you for your order!").
		Table(mailgen.Table{
			Headers: []mailgen.TableHeader{
				{Text: "Item", Align: "left", Width: "70%"},
				{Text: "Price", Align: "right", Width: "30%"},
			},
			Rows: [][]string{
				{"Widget A", "$10.00"},
				{"Widget B", "$15.00"},
			},
		}).
		Action("View Order", "https://example.com/order").
		Line("If you have any questions, feel free to contact us.")

	plaintext, err := msg.GeneratePlaintext()
	require.NoError(t, err)

	assert.Contains(t, plaintext, "Widget A")
	assert.Contains(t, plaintext, "$10.00")
}
