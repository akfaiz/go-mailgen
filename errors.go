package mailgen

import "errors"

var (
	// ErrInvalidThemeName indicates an empty theme name was provided.
	ErrInvalidThemeName = errors.New("mailgen: theme name cannot be empty")
	// ErrNilHTMLTemplate indicates a theme has no HTML template.
	ErrNilHTMLTemplate = errors.New("mailgen: theme HTML template cannot be nil")
)
