package mailer_test

import (
	"testing"
	"time"

	"github.com/ahmadfaizk/go-mailer"
	smtpmock "github.com/mocktools/go-smtp-mock/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wneessen/go-mail"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		client  *mail.Client
		opts    []mailer.Option
		wantNil bool
	}{
		{
			name:    "creates mailer with client and no options",
			client:  &mail.Client{},
			opts:    nil,
			wantNil: false,
		},
		{
			name:   "creates mailer with client and options",
			client: &mail.Client{},
			opts: []mailer.Option{
				mailer.WithFrom("Test Sender", "test@example.com"),
				mailer.WithProduct(mailer.Product{
					Name:    "Test Product",
					LogoURL: "https://example.com/logo.png",
					URL:     "https://example.com",
				}),
			},
			wantNil: false,
		},
		{
			name:    "creates mailer with nil client",
			client:  nil,
			opts:    nil,
			wantNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mailer.New(tt.client, tt.opts...)

			if tt.wantNil {
				assert.Nil(t, got)
				return
			}

			assert.NotNil(t, got)
		})
	}
}

func TestMailer_Send(t *testing.T) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		// LogToStdout:       true,
		// LogServerActivity: true,
	})
	if err := server.Start(); err != nil {
		t.Fatalf("failed to start mock SMTP server: %v", err)
	}
	defer server.Stop()
	client, err := mail.NewClient("127.0.0.1",
		mail.WithPort(server.PortNumber()),
		mail.WithSMTPAuth(mail.SMTPAuthNoAuth),
		mail.WithTLSPortPolicy(mail.NoTLS),
	)
	require.NoError(t, err)

	mailMailer := mailer.New(client)

	tests := []struct {
		name        string
		message     *mailer.Message
		wantErr     bool
		errContains string
		expectFunc  func(*testing.T)
	}{
		{
			name: "sends message successfully",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					Subject("Test Subject").
					Line("This is a test email.").
					Action("Click here", "https://example.com").
					Line("Thank you for using our service.")
				return msg
			}(),
			wantErr: false,
			expectFunc: func(t *testing.T) {
				messages, err := server.WaitForMessagesAndPurge(1, 1*time.Second)
				require.NoError(t, err)
				require.Len(t, messages, 1)
			},
		},
		{
			name: "sends message with CC and BCC",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					Cc("cc@example.com").
					Bcc("bcc@example.com").
					Subject("Test Subject").
					Line("This is a test email.").
					Action("Click here", "https://example.com").
					Line("Thank you for using our service.")
				return msg
			}(),
			wantErr: false,
			expectFunc: func(t *testing.T) {
				messages, err := server.WaitForMessagesAndPurge(1, 1*time.Second)
				require.NoError(t, err)
				require.Len(t, messages, 1)
			},
		},
		{
			name:        "returns error when message is nil",
			message:     nil,
			wantErr:     true,
			errContains: "message cannot be nil",
		},
		{
			name: "returns error when message has no recipients",
			message: func() *mailer.Message {
				msg := mailer.NewMessage()
				msg.Subject("Test Subject")
				return msg
			}(),
			wantErr:     true,
			errContains: "message must have at least one recipient",
		},
		{
			name: "returns error when message has no subject",
			message: func() *mailer.Message {
				msg := mailer.NewMessage()
				msg.To("recipient@example.com")
				return msg
			}(),
			wantErr:     true,
			errContains: "message must have a subject",
		},
		{
			name: "returns error when recipient email is invalid",
			message: func() *mailer.Message {
				msg := mailer.NewMessage()
				msg.To("invalid-email")
				msg.Subject("Test Subject")
				return msg
			}(),
			wantErr:     true,
			errContains: "failed to parse mail address",
		},
		{
			name: "returns error when cc email is invalid",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					Cc("invalid-email").
					Subject("Test Subject")
				return msg
			}(),
			wantErr:     true,
			errContains: "failed to parse mail address",
		},
		{
			name: "returns error when bcc email is invalid",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					Bcc("invalid-email").
					Subject("Test Subject")
				return msg
			}(),
			wantErr:     true,
			errContains: "failed to parse mail address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mailMailer.Send(tt.message)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				assert.NoError(t, err)
				if tt.expectFunc != nil {
					tt.expectFunc(t)
				}
			}
		})
	}
}
