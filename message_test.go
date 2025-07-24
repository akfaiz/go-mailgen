package mailer_test

import (
	"testing"
	"time"

	"github.com/ahmadfaizk/go-mailer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMessage(t *testing.T) {
	msg := mailer.NewMessage()

	assert.Empty(t, msg.GetSubject())
	assert.Empty(t, msg.GetTo())
	assert.Empty(t, msg.GetCc())
	assert.Empty(t, msg.GetBcc())
}

func TestMessage_Subject(t *testing.T) {
	msg := mailer.NewMessage()
	subject := "Test Subject"

	msg.Subject(subject)

	assert.Equal(t, subject, msg.GetSubject())
}

func TestMessage_To(t *testing.T) {
	msg := mailer.NewMessage()

	msg.To("test1@example.com", "test2@example.com")

	expected := []string{"test1@example.com", "test2@example.com"}
	assert.Equal(t, expected, msg.GetTo())
}

func TestMessage_Cc(t *testing.T) {
	msg := mailer.NewMessage()

	msg.Cc("cc1@example.com", "cc2@example.com")

	expected := []string{"cc1@example.com", "cc2@example.com"}
	assert.Equal(t, expected, msg.GetCc())
}

func TestMessage_Bcc(t *testing.T) {
	msg := mailer.NewMessage()

	msg.Bcc("bcc1@example.com", "bcc2@example.com")

	expected := []string{"bcc1@example.com", "bcc2@example.com"}
	assert.Equal(t, expected, msg.GetBcc())
}

func TestMessage_Line(t *testing.T) {
	msg := mailer.NewMessage()

	msg.Line("First line")
	msg.Line("Second line")

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "First line")
	assert.Contains(t, html, "Second line")
}

func TestMessage_Linef(t *testing.T) {
	msg := mailer.NewMessage()

	msg.Linef("Hello %s, you have %d messages", "John", 5)

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "Hello John, you have 5 messages")
}

func TestMessage_Action(t *testing.T) {
	msg := mailer.NewMessage()

	msg.Action("Click Here", "https://example.com")

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "Click Here")
	assert.Contains(t, html, "https://example.com")
}

func TestMessage_Product(t *testing.T) {
	msg := mailer.NewMessage()

	product := mailer.Product{
		Name:      "Test Product",
		LogoURL:   "https://example.com/logo.png",
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
	msg := mailer.NewMessage()

	msg.Product(mailer.Product{})

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "GoMailer")

	currentYear := time.Now().Year()
	expectedCopyright := "© " + string(rune(currentYear+'0'))
	assert.Contains(t, html, expectedCopyright[:3])
}

func TestMessage_ChainedMethods(t *testing.T) {
	msg := mailer.NewMessage().
		Subject("Chained Test").
		To("test@example.com").
		Cc("cc@example.com").
		Bcc("bcc@example.com").
		Greeting("Hi there").
		Line("This is a test").
		Action("Test Action", "https://test.com").
		Line("After action line")

	assert.Equal(t, "Chained Test", msg.GetSubject())
	assert.Equal(t, []string{"test@example.com"}, msg.GetTo())

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.Contains(t, html, "Hi there")
}

func TestMessage_GenerateHTML(t *testing.T) {
	msg := mailer.NewMessage().
		Subject("Test Email").
		Greeting("Hello").
		Line("This is a test email").
		Action("Visit Site", "https://example.com").
		Line("Thank you")

	html, err := msg.GenerateHTML()
	require.NoError(t, err)

	assert.NotEmpty(t, html)
	assert.Contains(t, html, "<html")
}
