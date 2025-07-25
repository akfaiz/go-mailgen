package component

import (
	"bytes"
	htmltemplate "html/template"
	texttemplate "text/template"
)

type Component interface {
	HTML(tmpl *htmltemplate.Template) (string, error)
	PlainText(tmpl *texttemplate.Template) (string, error)
}

var _ Component = &Table{}
var _ Component = &Action{}
var _ Component = &Line{}

type Action struct {
	Text  string
	Link  string
	Color string
}

type Line struct {
	Text string
}

type Table struct {
	Data    [][]Entry
	Columns Columns
}

type Entry struct {
	Key   string
	Value string
}

type Columns struct {
	CustomWidth map[string]string
	CustomAlign map[string]string
}

func (a Action) HTML(tmpl *htmltemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "button", a)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a Action) PlainText(tmpl *texttemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "button", a)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (l Line) HTML(tmpl *htmltemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "line", l)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (l Line) PlainText(tmpl *texttemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "line", l)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t Table) HTML(tmpl *htmltemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "table", t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t Table) PlainText(tmpl *texttemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "table", t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
