package mailgen_test

import (
	htmltemplate "html/template"
	"testing"
	texttemplate "text/template"

	"github.com/akfaiz/go-mailgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegisterTheme_CustomTheme(t *testing.T) {
	htmlTmpl := htmltemplate.Must(htmltemplate.New("index.html").Parse(`
		{{define "index.html"}}CUSTOM_HTML|{{.Greeting}}|{{range .ComponentsHTML}}{{.}}{{end}}{{end}}
		{{define "line"}}<span>CUSTOM_LINE:{{.Text}}</span>{{end}}
	`))
	textTmpl := texttemplate.Must(texttemplate.New("index.txt").Parse(`
		{{define "index.txt"}}CUSTOM_TEXT|{{.Greeting}}|{{range .ComponentsText}}{{.}}{{end}}{{end}}
	`))

	err := mailgen.RegisterTheme("my-custom-theme", mailgen.Theme{
		HTML:      htmlTmpl,
		PlainText: textTmpl,
	})
	require.NoError(t, err)

	msg, err := mailgen.New().
		Theme("my-custom-theme").
		Greeting("Hello").
		Name("John").
		UsePremailer(false).
		Line("Welcome").
		Build()
	require.NoError(t, err)

	assert.Contains(t, msg.HTML(), "CUSTOM_HTML|Hello John|<span>CUSTOM_LINE:Welcome</span>")
	assert.Contains(t, msg.PlainText(), "CUSTOM_TEXT|Hello John|Welcome")
}

func TestRegisterTheme_DefaultPlainTextFallback(t *testing.T) {
	htmlTmpl := htmltemplate.Must(htmltemplate.New("index.html").Parse(`
		{{define "index.html"}}CUSTOM_HTML_ONLY|{{range .ComponentsHTML}}{{.}}{{end}}{{end}}
		{{define "line"}}<p>{{.Text}}</p>{{end}}
	`))

	err := mailgen.RegisterTheme("html-only-theme", mailgen.Theme{HTML: htmlTmpl})
	require.NoError(t, err)

	msg, err := mailgen.New().
		Theme("html-only-theme").
		Line("Fallback plain text body").
		Build()
	require.NoError(t, err)

	assert.Contains(t, msg.HTML(), "CUSTOM_HTML_ONLY")
	assert.Contains(t, msg.PlainText(), "Fallback plain text body")
}

func TestRegisterTheme_Validation(t *testing.T) {
	err := mailgen.RegisterTheme("", mailgen.Theme{})
	require.ErrorIs(t, err, mailgen.ErrInvalidThemeName)

	err = mailgen.RegisterTheme("nil-html", mailgen.Theme{})
	require.ErrorIs(t, err, mailgen.ErrNilHTMLTemplate)
}
