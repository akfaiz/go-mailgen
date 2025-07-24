package mailer_test

import (
	"bytes"
	"embed"
	"os"
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
					Name: "Test Product",
					URL:  "https://example.com",
				}),
				mailer.WithReplyTo("reply@example.com"),
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

//go:embed testdata/*
var testDataFS embed.FS

func TestMailer_Send(t *testing.T) {
	server := smtpmock.New(smtpmock.ConfigurationAttr{
		HostAddress: "localhost",
		// LogToStdout:       true,
		// LogServerActivity: true,
	})
	if err := server.Start(); err != nil {
		t.Fatalf("failed to start mock SMTP server: %v", err)
	}
	defer func() {
		_ = server.Stop()
	}()

	client, err := mail.NewClient("localhost",
		mail.WithPort(server.PortNumber()),
		mail.WithSMTPAuth(mail.SMTPAuthNoAuth),
		mail.WithTLSPortPolicy(mail.NoTLS),
		mail.WithHELO("localhost"),
	)
	require.NoError(t, err)

	mailMailer := mailer.New(client, mailer.WithFrom("noreply@example.com", "No Reply"))

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
			name: "send message with attachment from file",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					Subject("Test Subject with Attachment").
					Line("This is a test email with an attachment.").
					AttachFile("testdata/golang.png",
						mailer.WithFileContentType(mailer.TypeAppOctetStream),
						mailer.WithFileEncoding(mailer.EncodingB64),
						mailer.WithFileName("golang.png"),
						mailer.WithFileDescription("Go Programming Language Logo"),
						mailer.WithFileContentID("golang-logo"),
					)
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
			name: "send message with attachment from embed FS",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					Subject("Test Subject with Attachment").
					Line("This is a test email with an attachment.").
					AttachFromEmbedFS("testdata/golang.png", &testDataFS)
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
			name: "send message with attachment from IOFS",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					Subject("Test Subject with Attachment").
					Line("This is a test email with an attachment.").
					AttachFromIOFS("golang.png", os.DirFS("testdata"))
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
			name: "send message with attachment from reader",
			message: func() *mailer.Message {
				var buffer bytes.Buffer
				_, err := buffer.ReadFrom(bytes.NewReader([]byte("This is a test attachment.")))
				require.NoError(t, err)
				reader := bytes.NewReader(buffer.Bytes())

				msg := mailer.NewMessage().
					To("recipient@example.com").
					Subject("Test Subject with Attachment").
					Line("This is a test email with an attachment.").
					AttachReader("test.txt", reader)
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
			name: "send message with attachment from read seeker",
			message: func() *mailer.Message {
				file, err := os.Open("testdata/golang.png")
				require.NoError(t, err)

				msg := mailer.NewMessage().
					To("recipient@example.com").
					Subject("Test Subject with Attachment").
					Line("This is a test email with an attachment.").
					AttachReadSeeker("golang.png", file)
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
			name: "return errors when from email is invalid",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					From("invalid-email").
					Subject("Test Subject")
				return msg
			}(),
			wantErr:     true,
			errContains: "failed to parse mail address",
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
		{
			name: "return errors when reply-to email is invalid",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					ReplyTo("invalid-email").
					Subject("Test Subject")
				return msg
			}(),
			wantErr:     true,
			errContains: "failed to parse reply-to address",
		},
		{
			name: "return errors when attachment embedded file is not found",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					Subject("Test Subject").
					Line("This is a test email with an attachment.").
					AttachFromEmbedFS("testdata/golang.jpg", &testDataFS)
				return msg
			}(),
			wantErr:     true,
			errContains: "file does not exist",
		},
		{
			name: "return errors when attachment IOFS file is not found",
			message: func() *mailer.Message {
				msg := mailer.NewMessage().
					To("recipient@example.com").
					Subject("Test Subject").
					Line("This is a test email with an attachment.").
					AttachFromIOFS("testdata/golang.jpg", os.DirFS("testdata"))
				return msg
			}(),
			wantErr:     true,
			errContains: "no such file or directory",
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
