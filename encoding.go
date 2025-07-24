package mailer

// ContentType is a type wrapper for a string and represents the MIME type of the content being handled.
type ContentType string

// Encoding is a type wrapper for a string and represents the type of encoding used for email messages
// and/or parts.
type Encoding string

const (
	// EncodingB64 represents the Base64 encoding as specified in RFC 2045.
	EncodingB64 Encoding = "base64"

	// EncodingQP represents the "quoted-printable" encoding as specified in RFC 2045.
	EncodingQP Encoding = "quoted-printable"

	// EncodingUSASCII represents encoding with only US-ASCII characters (aka 7Bit)
	EncodingUSASCII Encoding = "7bit"

	// NoEncoding represents 8-bit encoding for email messages as specified in RFC 6152.
	NoEncoding Encoding = "8bit"
)

const (
	// TypeAppOctetStream represents the MIME type for arbitrary binary data.
	TypeAppOctetStream ContentType = "application/octet-stream"

	// TypeMultipartAlternative represents the MIME type for a message body that can contain multiple alternative
	// formats.
	TypeMultipartAlternative ContentType = "multipart/alternative"

	// TypeMultipartMixed represents the MIME type for a multipart message containing different parts.
	TypeMultipartMixed ContentType = "multipart/mixed"

	// TypeMultipartRelated represents the MIME type for a multipart message where each part is a related file
	// or resource.
	TypeMultipartRelated ContentType = "multipart/related"

	// TypePGPSignature represents the MIME type for PGP signed messages.
	TypePGPSignature ContentType = "application/pgp-signature"

	// TypePGPEncrypted represents the MIME type for PGP encrypted messages.
	TypePGPEncrypted ContentType = "application/pgp-encrypted"

	// TypeTextHTML represents the MIME type for HTML text content.
	TypeTextHTML ContentType = "text/html"

	// TypeTextPlain represents the MIME type for plain text content.
	TypeTextPlain ContentType = "text/plain"

	// TypeSMIMESigned represents the MIME type for S/MIME singed messages.
	TypeSMIMESigned ContentType = `application/pkcs7-signature; name="smime.p7s"`
)
