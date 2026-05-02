package mailgen

import (
	htmltemplate "html/template"
	"strings"
	"sync"
	texttemplate "text/template"

	"github.com/akfaiz/go-mailgen/templates"
)

// Theme defines HTML and plain text templates for a named mail theme.
//
// HTML must include all component templates (for example: index.html, line, button, table)
// expected by the builder.
//
// PlainText is optional. If omitted, the default plain-text template will be used.
type Theme struct {
	HTML      *htmltemplate.Template
	PlainText *texttemplate.Template
}

var (
	themeRegistryMu sync.RWMutex
	themeRegistry   = map[string]Theme{
		"default": {
			HTML:      templates.DefaultHTMLTmpl,
			PlainText: templates.DefaultPlainTextTmpl,
		},
		"plain": {
			HTML:      templates.PlainHTMLTmpl,
			PlainText: templates.DefaultPlainTextTmpl,
		},
	}
)

// RegisterTheme registers a custom theme that can be selected via Builder.Theme(name).
// Name is case-insensitive and stored in lowercase.
func RegisterTheme(name string, theme Theme) error {
	name = normalizeThemeName(name)
	if name == "" {
		return ErrInvalidThemeName
	}
	if theme.HTML == nil {
		return ErrNilHTMLTemplate
	}
	if theme.PlainText == nil {
		theme.PlainText = templates.DefaultPlainTextTmpl
	}

	themeRegistryMu.Lock()
	themeRegistry[name] = theme
	themeRegistryMu.Unlock()

	return nil
}

// MustRegisterTheme registers a custom theme and panics if registration fails.
func MustRegisterTheme(name string, theme Theme) {
	if err := RegisterTheme(name, theme); err != nil {
		panic(err)
	}
}

func resolveTheme(name string) Theme {
	name = normalizeThemeName(name)
	themeRegistryMu.RLock()
	theme, ok := themeRegistry[name]
	if !ok {
		theme = themeRegistry["default"]
	}
	themeRegistryMu.RUnlock()
	return theme
}

func normalizeThemeName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
