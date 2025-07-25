package mailgen_test

import (
	"testing"

	"github.com/ahmadfaizk/go-mailgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCase struct {
	name        string
	builderFunc func() *mailgen.Builder
	expectError bool
	expectFunc  func(msg mailgen.Message)
}

func (tc testCase) run(t *testing.T, modifyBuilder ...func(*mailgen.Builder)) {
	t.Run(tc.name, func(t *testing.T) {
		builder := tc.builderFunc()
		for _, cb := range modifyBuilder {
			cb(builder)
		}
		msg, err := builder.Build()
		if tc.expectError {
			assert.Error(t, err, "Build should return an error")
			return
		}
		require.NoError(t, err, "Build should not return an error")
		assert.NotNil(t, msg, "Build should return a non-nil Message")
		tc.expectFunc(msg)
	})
}

func TestBuilder_Subject(t *testing.T) {
	testCases := []testCase{
		{
			name: "set subject",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Subject("Test Subject")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Equal(t, "Test Subject", msg.Subject(), "Subject should match the set value")
			},
		},
		{
			name: "not set subject",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New()
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Equal(t, "", msg.Subject(), "Subject should be empty")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_From(t *testing.T) {
	testCases := []testCase{
		{
			name: "set from",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().From("sender@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Equal(t, "sender@example.com", msg.From().String(), "From should match the set value")
			},
		},
		{
			name: "set from with name",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().From("sender@example.com", "Sender Name")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Equal(t, "Sender Name <sender@example.com>", msg.From().String(), "From should match the set value")
			},
		},
		{
			name: "not set from",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New()
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Equal(t, "", msg.From().String(), "From should be empty")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_To(t *testing.T) {
	testCases := []testCase{
		{
			name: "set single recipient",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().To("user1@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.To(), 1, "To should contain one recipient")
				assert.Contains(t, msg.To(), "user1@example.com", "To should contain the added recipient")
			},
		},
		{
			name: "set multiple recipients",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().To("user2@example.com", "user3@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.To(), 2, "To should contain two recipients")
				assert.Contains(t, msg.To(), "user2@example.com", "To should contain the added recipient")
				assert.Contains(t, msg.To(), "user3@example.com", "To should contain the added recipient")
			},
		},
		{
			name: "set no recipients",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New()
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.To(), "To should be empty when no recipients are set")
			},
		},
		{
			name: "set empty recipient",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().To("")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.To(), "To should be empty when an empty recipient is set")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_CC(t *testing.T) {
	testCases := []testCase{
		{
			name: "set single CC",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().CC("cc1@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.CC(), 1, "CC should contain one recipient")
				assert.Contains(t, msg.CC(), "cc1@example.com", "CC should contain the added recipient")
			},
		},
		{
			name: "set multiple CCs",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().CC("cc2@example.com", "cc3@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.CC(), 2, "CC should contain two recipients")
				assert.Contains(t, msg.CC(), "cc2@example.com", "CC should contain the added recipient")
				assert.Contains(t, msg.CC(), "cc3@example.com", "CC should contain the added recipient")
			},
		},
		{
			name: "set no CCs",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New()
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.CC(), "CC should be empty when no recipients are set")
			},
		},
		{
			name: "set empty CC",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().CC("")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.CC(), "CC should be empty when an empty recipient is set")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_BCC(t *testing.T) {
	testCases := []testCase{
		{
			name: "set single BCC",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().BCC("bcc1@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.BCC(), 1, "BCC should contain one recipient")
				assert.Contains(t, msg.BCC(), "bcc1@example.com", "BCC should contain the added recipient")
			},
		},
		{
			name: "set multiple BCCs",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().BCC("bcc2@example.com", "bcc3@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.BCC(), 2, "BCC should contain two recipients")
				assert.Contains(t, msg.BCC(), "bcc2@example.com", "BCC should contain the added recipient")
				assert.Contains(t, msg.BCC(), "bcc3@example.com", "BCC should contain the added recipient")
			},
		},
		{
			name: "set no BCCs",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New()
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.BCC(), "BCC should be empty when no recipients are set")
			},
		},
		{
			name: "set empty BCC",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().BCC("")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.BCC(), "BCC should be empty when an empty recipient is set")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Preheader(t *testing.T) {
	testCases := []testCase{
		{
			name: "set preheader",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Preheader("This is a preheader text")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "This is a preheader text", "HTML should contain the preheader text")
				assert.Contains(t, msg.PlainText(), "This is a preheader text", "PlainText should contain the preheader text")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Greeting(t *testing.T) {
	testCases := []testCase{
		{
			name: "set greeting",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Greeting("Hi there")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "Hi there", "HTML should contain the greeting text")
				assert.Contains(t, msg.PlainText(), "Hi there", "PlainText should contain the greeting text")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Salutation(t *testing.T) {
	testCases := []testCase{
		{
			name: "set salutation",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Salutation("Kind regards")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "Kind regards", "HTML should contain the salutation text")
				assert.Contains(t, msg.PlainText(), "Kind regards", "PlainText should contain the salutation text")
			},
		},
		{
			name: "not set salutation should use default",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New()
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "Best regards", "HTML should contain the default salutation text")
				assert.Contains(t, msg.PlainText(), "Best regards", "PlainText should contain the default salutation text")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Line(t *testing.T) {
	testCases := []testCase{
		{
			name: "add line",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Line("This is a line of text")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "This is a line of text", "HTML should contain the line text")
				assert.Contains(t, msg.PlainText(), "This is a line of text", "PlainText should contain the line text")
			},
		},
		{
			name: "add multiple lines",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Line("First line").Line("Second line")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "First line", "HTML should contain the first line text")
				assert.Contains(t, msg.HTML(), "Second line", "HTML should contain the second line text")
				assert.Contains(t, msg.PlainText(), "First line", "PlainText should contain the first line text")
				assert.Contains(t, msg.PlainText(), "Second line", "PlainText should contain the second line text")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Linef(t *testing.T) {
	testCases := []testCase{
		{
			name: "add formatted line",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Linef("Hello %s, your order #%d is ready", "John", 12345)
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				expectedText := "Hello John, your order #12345 is ready"
				assert.Contains(t, msg.HTML(), expectedText, "HTML should contain the formatted line text")
				assert.Contains(t, msg.PlainText(), expectedText, "PlainText should contain the formatted line text")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Action(t *testing.T) {
	testCases := []testCase{
		{
			name: "add action with default color",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Action("Click Here", "https://example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "Click Here", "HTML should contain the action text")
				assert.Contains(t, msg.HTML(), "https://example.com", "HTML should contain the action URL")
				assert.Contains(t, msg.HTML(), "background-color:#3869D4", "HTML should contain the default action color")
				assert.Contains(t, msg.PlainText(), "Click Here", "PlainText should contain the action text")
				assert.Contains(t, msg.PlainText(), "https://example.com", "PlainText should contain the action URL")
			},
		},
		{
			name: "add action with custom color",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Action("Custom Button", "https://custom.com", "#FF0000")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "Custom Button", "HTML should contain the custom action text")
				assert.Contains(t, msg.HTML(), "https://custom.com", "HTML should contain the custom action URL")
				assert.Contains(t, msg.HTML(), "background-color:#FF0000", "HTML should contain the custom action color")
				assert.Contains(t, msg.PlainText(), "Custom Button", "PlainText should contain the custom action text")
				assert.Contains(t, msg.PlainText(), "https://custom.com", "PlainText should contain the custom action URL")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Product(t *testing.T) {
	testCases := []testCase{
		{
			name: "set product with complete info",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Product(mailgen.Product{
					Name:      "Test Product",
					URL:       "https://test.com",
					Copyright: "© 2023 Test Product",
				})
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "Test Product", "HTML should contain the product name")
				assert.Contains(t, msg.HTML(), "https://test.com", "HTML should contain the product URL")
				assert.Contains(t, msg.HTML(), "© 2023 Test Product", "HTML should contain the product copyright")
				assert.Contains(t, msg.PlainText(), "Test Product", "PlainText should contain the product name")
				assert.Contains(t, msg.PlainText(), "https://test.com", "PlainText should contain the product URL")
				assert.Contains(t, msg.PlainText(), "© 2023 Test Product", "PlainText should contain the product copyright")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Table(t *testing.T) {
	testCases := []testCase{
		{
			name: "simple table with headers and rows",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Table(mailgen.Table{
					Headers: []mailgen.TableHeader{
						{Text: "Item", Align: "left", Width: "70%"},
						{Text: "Price"},
					},
					Rows: [][]string{
						{"Widget A", "$10.00"},
						{"Widget B", "$15.00"},
					},
				})
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "<table", "HTML should contain a table")
				assert.Contains(t, msg.HTML(), ">Item</", "HTML should contain the Item header")
				assert.Contains(t, msg.HTML(), ">Widget A</", "HTML should contain the first row item")
				assert.Contains(t, msg.HTML(), ">$10.00</", "HTML should contain the first row price")
				assert.Contains(t, msg.PlainText(), "Item", "PlainText should contain the Item header")
				assert.Contains(t, msg.PlainText(), "Widget A", "PlainText should contain the first row item")
				assert.Contains(t, msg.PlainText(), "$10.00", "PlainText should contain the first row price")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Build(t *testing.T) {
	themes := []string{"default", "plain"}
	testCases := []testCase{
		{
			name: "reset password message",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().
					Subject("Reset your password").
					From("no-reply@example.com").
					To("user@example.com").
					Line("Click the link below to reset your password:").
					Action("Reset Password", "https://example.com/reset").
					Line("If you did not request this, please ignore this email.")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Equal(t, "Reset your password", msg.Subject(), "Subject should match the set value")
				assert.Equal(t, "no-reply@example.com", msg.From().String(), "From should match the set value")
				assert.Equal(t, []string{"user@example.com"}, msg.To(), "To should match the set value")

				assert.Contains(t, msg.HTML(), "Click the link below to reset your password:", "HTML should contain the line text")
				assert.Contains(t, msg.HTML(), "Reset Password", "HTML should contain the action text")
				assert.Contains(t, msg.HTML(), "https://example.com/reset", "HTML should contain the action URL")

				assert.Contains(t, msg.PlainText(), "Click the link below to reset your password:", "PlainText should match the set value")
				assert.Contains(t, msg.PlainText(), "Reset Password", "PlainText should contain the action text")
				assert.Contains(t, msg.PlainText(), "https://example.com/reset", "PlainText should contain the action URL")
			},
		},
	}
	for _, theme := range themes {
		for _, tc := range testCases {
			tc.name = theme + " " + tc.name
			tc.run(t, func(b *mailgen.Builder) {
				b.Theme(theme)
			})
		}
	}
}
