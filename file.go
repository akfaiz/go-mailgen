package mailer

import (
	"embed"
	"io"
	"io/fs"

	"github.com/wneessen/go-mail"
)

type file struct {
	name string
	cfg  *fileConfig
}

type fileEmbedFS struct {
	file
	fs *embed.FS
}

type fileIOFS struct {
	file
	fs.FS
}

type fileReader struct {
	file
	io.Reader
}

type fileReadSeeker struct {
	file
	io.ReadSeeker
}

type fileConfig struct {
	contentID   string
	name        string
	desc        string
	enc         Encoding
	contentType ContentType
}

func newFileConfig(opts ...FileOption) *fileConfig {
	fc := &fileConfig{}
	for _, opt := range opts {
		opt(fc)
	}
	return fc
}

func (fc *fileConfig) toMailFileOption() []mail.FileOption {
	var opts []mail.FileOption
	if fc.contentID != "" {
		opts = append(opts, mail.WithFileContentID(fc.contentID))
	}
	if fc.name != "" {
		opts = append(opts, mail.WithFileName(fc.name))
	}
	if fc.desc != "" {
		opts = append(opts, mail.WithFileDescription(fc.desc))
	}
	if fc.enc != "" {
		opts = append(opts, mail.WithFileEncoding(mail.Encoding(fc.enc)))
	}
	if fc.contentType != "" {
		opts = append(opts, mail.WithFileContentType(mail.ContentType(fc.contentType)))
	}
	return opts
}

// FileOption is a function that modifies the file configuration for a file attachment.
type FileOption func(*fileConfig)

// WithFileContentID sets the "Content-ID" header in the file configuration.
func WithFileContentID(id string) FileOption {
	return func(f *fileConfig) {
		f.contentID = id
	}
}

// WithFileName sets the name of the file attachment.
func WithFileName(name string) FileOption {
	return func(f *fileConfig) {
		f.name = name
	}
}

// WithFileDescription sets an optional description for the file, which is used in the Content-Description header.
func WithFileDescription(description string) FileOption {
	return func(f *fileConfig) {
		f.desc = description
	}
}

// WithFileEncoding sets the encoding type for a file attachment.
func WithFileEncoding(encoding Encoding) FileOption {
	return func(f *fileConfig) {
		f.enc = encoding
	}
}

// WithFileContentType sets the MIME type for the file attachment.
func WithFileContentType(contentType ContentType) FileOption {
	return func(f *fileConfig) {
		f.contentType = contentType
	}
}
