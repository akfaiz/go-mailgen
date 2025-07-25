package templates

import (
	"embed"
	htmltemplate "html/template"
	texttemplate "text/template"
)

//go:embed default/*
var defaultFS embed.FS

//go:embed plain/*
var plainFS embed.FS

var DefaultHtmlTmpl = htmltemplate.Must(htmltemplate.New("index.html").
	ParseFS(defaultFS, "default/*.html"),
)
var DefaultPlainTextTmpl = texttemplate.Must(texttemplate.New("index.txt").
	ParseFS(defaultFS, "default/*.txt"),
)
var PlainHtmlTmpl = htmltemplate.Must(htmltemplate.New("index.html").
	ParseFS(plainFS, "plain/*.html"),
)
var PlainPlainTextTmpl = texttemplate.Must(texttemplate.New("index.txt").
	ParseFS(plainFS, "plain/*.txt"),
)
