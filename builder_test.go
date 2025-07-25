package mailgen_test

import (
	"fmt"
	"testing"
	"time"

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

func TestSetDefault(t *testing.T) {
	originalDefault := mailgen.New()
	defer mailgen.SetDefault(originalDefault) // Restore original default after tests

	t.Run("set valid builder as default", func(t *testing.T) {
		customBuilder := mailgen.New()
		customBuilder.Subject("Custom Default Subject").
			Greeting("Custom Greeting").
			Salutation("Custom Salutation").
			Product(mailgen.Product{
				Name:      "Custom Product",
				URL:       "https://custom.com",
				Copyright: "© 2023 Custom",
			})

		mailgen.SetDefault(customBuilder)

		// Create a new message and verify it uses the custom defaults
		msg, err := mailgen.New().Build()
		require.NoError(t, err)

		assert.Equal(t, "Custom Default Subject", msg.Subject())
		assert.Contains(t, msg.HTML(), "Custom Greeting")
		assert.Contains(t, msg.HTML(), "Custom Salutation")
		assert.Contains(t, msg.HTML(), "Custom Product")
		assert.Contains(t, msg.HTML(), "https://custom.com")
		assert.Contains(t, msg.HTML(), "© 2023 Custom")
	})

	t.Run("set nil builder should not change default", func(t *testing.T) {
		// First set a custom default
		customBuilder := mailgen.New()
		customBuilder.Subject("Before Nil Test")
		mailgen.SetDefault(customBuilder)

		// Verify the custom default is set
		msg1, err := mailgen.New().Build()
		require.NoError(t, err)
		assert.Equal(t, "Before Nil Test", msg1.Subject())

		// Try to set nil
		mailgen.SetDefault(nil)

		// Verify the default hasn't changed
		msg2, err := mailgen.New().Build()
		require.NoError(t, err)
		assert.Equal(t, "Before Nil Test", msg2.Subject())
	})

	t.Run("new instances are independent after setting default", func(t *testing.T) {
		customBuilder := mailgen.New()
		customBuilder.Subject("Base Subject")
		mailgen.SetDefault(customBuilder)

		// Create two new instances
		builder1 := mailgen.New().Subject("Modified Subject 1")
		builder2 := mailgen.New().Subject("Modified Subject 2")

		msg1, err := builder1.Build()
		require.NoError(t, err)
		msg2, err := builder2.Build()
		require.NoError(t, err)

		// Verify they have different subjects
		assert.Equal(t, "Modified Subject 1", msg1.Subject())
		assert.Equal(t, "Modified Subject 2", msg2.Subject())

		// Verify a new unmodified instance still has the default
		msg3, err := mailgen.New().Build()
		require.NoError(t, err)
		assert.Equal(t, "Base Subject", msg3.Subject())
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

func TestBuilder_Cc(t *testing.T) {
	testCases := []testCase{
		{
			name: "set single CC",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Cc("cc1@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.Cc(), 1, "CC should contain one recipient")
				assert.Contains(t, msg.Cc(), "cc1@example.com", "CC should contain the added recipient")
			},
		},
		{
			name: "set multiple CCs",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Cc("cc2@example.com", "cc3@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.Cc(), 2, "CC should contain two recipients")
				assert.Contains(t, msg.Cc(), "cc2@example.com", "CC should contain the added recipient")
				assert.Contains(t, msg.Cc(), "cc3@example.com", "CC should contain the added recipient")
			},
		},
		{
			name: "set no CCs",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New()
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.Cc(), "CC should be empty when no recipients are set")
			},
		},
		{
			name: "set empty CC",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Cc("")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.Cc(), "CC should be empty when an empty recipient is set")
			},
		},
	}
	for _, tc := range testCases {
		tc.run(t)
	}
}

func TestBuilder_Bcc(t *testing.T) {
	testCases := []testCase{
		{
			name: "set single BCC",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Bcc("bcc1@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.Bcc(), 1, "BCC should contain one recipient")
				assert.Contains(t, msg.Bcc(), "bcc1@example.com", "BCC should contain the added recipient")
			},
		},
		{
			name: "set multiple BCCs",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Bcc("bcc2@example.com", "bcc3@example.com")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Len(t, msg.Bcc(), 2, "BCC should contain two recipients")
				assert.Contains(t, msg.Bcc(), "bcc2@example.com", "BCC should contain the added recipient")
				assert.Contains(t, msg.Bcc(), "bcc3@example.com", "BCC should contain the added recipient")
			},
		},
		{
			name: "set no BCCs",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New()
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.Bcc(), "BCC should be empty when no recipients are set")
			},
		},
		{
			name: "set empty BCC",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Bcc("")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Empty(t, msg.Bcc(), "BCC should be empty when an empty recipient is set")
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
		{
			name: "set product with only name",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Product(mailgen.Product{
					Name: "Test Product",
				})
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				defaultCopyright := fmt.Sprintf("© %d Test Product. All rights reserved.", time.Now().Year())
				assert.Contains(t, msg.HTML(), "Test Product", "HTML should contain the product name")
				assert.Contains(t, msg.HTML(), defaultCopyright, "HTML should contain the default product copyright")
				assert.Contains(t, msg.PlainText(), "Test Product", "PlainText should contain the product name")
				assert.Contains(t, msg.PlainText(), defaultCopyright, "PlainText should contain the default product copyright")
			},
		},
		{
			name: "set product with only copyright",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().Product(mailgen.Product{
					Copyright: "© 2023 Test Product. All rights reserved.",
				})
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				defaultProductName := "Go-Mailgen"
				assert.Contains(t, msg.HTML(), defaultProductName, "HTML should contain the default product name")
				assert.Contains(t, msg.HTML(), "© 2023 Test Product. All rights reserved.", "HTML should contain the product copyright")
				assert.Contains(t, msg.PlainText(), "© 2023 Test Product. All rights reserved.", "PlainText should contain the product copyright")
				assert.Contains(t, msg.PlainText(), defaultProductName, "PlainText should contain the default product name")
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
					Data: [][]mailgen.Entry{
						{{Key: "Item", Value: "Widget A"}, {Key: "Price", Value: "$10.00"}, {Key: "Count", Value: "2"}, {Key: "Total", Value: "$20.00"}},
						{{Key: "Item", Value: "Widget B"}, {Key: "Price", Value: "$150.00"}, {Key: "Count", Value: "1"}, {Key: "Total", Value: "$150.00"}},
					},
					Columns: mailgen.Columns{
						CustomAlign: map[string]string{
							"Item":  "left",
							"Price": "center",
							"Count": "center",
							"Total": "right",
						},
						CustomWidth: map[string]string{
							"Item":  "40%",
							"Price": "20%",
							"Count": "20%",
							"Total": "20%",
						},
					},
				})
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				htmlContains := []string{
					"<table", ">Item</", ">Widget A</", ">$10.00</", ">2</", ">$20.00</", ">Widget B</", ">$150.00</",
					">1</",
				}
				for _, str := range htmlContains {
					assert.Contains(t, msg.HTML(), str, fmt.Sprintf("HTML should contain '%s'", str))
				}
				textContains := []string{"Item", "Widget A", "$10.00", "2", "$20.00", "Widget B", "$150.00", "1"}
				for _, str := range textContains {
					assert.Contains(t, msg.PlainText(), str, fmt.Sprintf("PlainText should contain '%s'", str))
				}
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
					Line("Click the link below to reset your password:").
					Action("Reset Password", "https://example.com/reset").
					Line("If you did not request this, please ignore this email.")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "Click the link below to reset your password:", "HTML should contain the line text")
				assert.Contains(t, msg.HTML(), "Reset Password", "HTML should contain the action text")
				assert.Contains(t, msg.HTML(), "https://example.com/reset", "HTML should contain the action URL")

				assert.Contains(t, msg.PlainText(), "Click the link below to reset your password:", "PlainText should match the set value")
				assert.Contains(t, msg.PlainText(), "Reset Password", "PlainText should contain the action text")
				assert.Contains(t, msg.PlainText(), "https://example.com/reset", "PlainText should contain the action URL")
			},
		},
		{
			name: "welcome message",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().
					Line("Welcome to our service!").
					Line("We're glad to have you on board.").
					Line("If you have any questions, feel free to reach out to our support team.")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "Welcome to our service!", "HTML should contain the welcome message")
				assert.Contains(t, msg.HTML(), "We&#39;re glad to have you on board.", "HTML should contain the second line")
				assert.Contains(t, msg.HTML(), "If you have any questions, feel free to reach out to our support team.", "HTML should contain the third line")

				assert.Contains(t, msg.PlainText(), "Welcome to our service!", "PlainText should match the set value")
				assert.Contains(t, msg.PlainText(), "We're glad to have you on board.", "PlainText should contain the second line")
				assert.Contains(t, msg.PlainText(), "If you have any questions, feel free to reach out to our support team.", "PlainText should contain the third line")
			},
		},
		{
			name: "invoice message",
			builderFunc: func() *mailgen.Builder {
				return mailgen.New().
					Line("Thank you for your purchase!").
					Line("Below are the details of your order:").
					Table(mailgen.Table{
						Data: [][]mailgen.Entry{
							{{Key: "Item", Value: "Widget A"}, {Key: "Price", Value: "$10.00"}},
							{{Key: "Item", Value: "Widget B"}, {Key: "Price", Value: "$15.00"}},
							{{Key: "Item", Value: "Widget C"}, {Key: "Price", Value: "$20.00"}},
							{{Key: "Item", Value: "Total"}, {Key: "Price", Value: "$45.00"}},
						},
						Columns: mailgen.Columns{
							CustomAlign: map[string]string{
								"Item":  "left",
								"Price": "right",
							},
							CustomWidth: map[string]string{
								"Item":  "70%",
								"Price": "30%",
							},
						},
					}).
					Line("Click the button below to track your order.").
					Action("Track Order", "https://example.com/track-order").
					Line("If you have any questions, please contact our support team.")
			},
			expectError: false,
			expectFunc: func(msg mailgen.Message) {
				assert.Contains(t, msg.HTML(), "Thank you for your purchase!", "HTML should contain the thank you message")
				assert.Contains(t, msg.HTML(), "Below are the details of your order:", "HTML should contain the order details message")
				assert.Contains(t, msg.HTML(), "<table", "HTML should contain a table")
				assert.Contains(t, msg.HTML(), ">Widget A</", "HTML should contain the first row item")
				assert.Contains(t, msg.HTML(), ">$10.00</", "HTML should contain the first row price")
				assert.Contains(t, msg.HTML(), "Track Order", "HTML should contain the action text")
				assert.Contains(t, msg.HTML(), "https://example.com/track-order", "HTML should contain the action URL")

				assert.Contains(t, msg.PlainText(), "Thank you for your purchase!", "PlainText should match the set value")
				assert.Contains(t, msg.PlainText(), "Below are the details of your order:", "PlainText should contain the order details message")
				assert.Contains(t, msg.PlainText(), "Widget A", "PlainText should contain the first row item")
				assert.Contains(t, msg.PlainText(), "$10.00", "PlainText should contain the first row price")
				assert.Contains(t, msg.PlainText(), "Track Order", "PlainText should contain the action text")
				assert.Contains(t, msg.PlainText(), "https://example.com/track-order", "PlainText should contain the action URL")
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
