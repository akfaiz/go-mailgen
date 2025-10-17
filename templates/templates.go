package templates

import (
	"embed"
	htmltemplate "html/template"
	"strings"
	texttemplate "text/template"
	"unicode"
)

//go:embed default/*
var defaultFS embed.FS

var DefaultHTMLTmpl = htmltemplate.Must(htmltemplate.New("index.html").
	Funcs(htmlTemplateFuncs).
	ParseFS(defaultFS, "default/*.html"),
)

var DefaultPlainTextTmpl = texttemplate.Must(texttemplate.New("index.txt").
	Funcs(textTemplateFuncs).
	ParseFS(defaultFS, "default/*.txt"),
)

//go:embed plain/*
var plainFS embed.FS

var PlainHTMLTmpl = htmltemplate.Must(htmltemplate.New("index.html").
	Funcs(htmlTemplateFuncs).
	ParseFS(plainFS, "plain/*.html"),
)

var htmlTemplateFuncs = htmltemplate.FuncMap{
	"capitalize": capitalize,
}
var textTemplateFuncs = texttemplate.FuncMap{
	"boxString": boxString,
	"concat":    concat,
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func concat(a, b string) string {
	return a + b
}

func boxString(s string) string {
	// Find the max line length (in case of multi-line input)
	lines := strings.Split(s, "\n")
	maxLen := 0
	for _, line := range lines {
		if len(line) > maxLen {
			maxLen = len(line)
		}
	}

	// The border should match the longest line length
	border := strings.Repeat("*", maxLen)

	// Combine
	var b strings.Builder
	b.WriteString(border + "\n")
	b.WriteString(s + "\n")
	b.WriteString(border)

	return b.String()
}
