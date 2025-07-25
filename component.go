package mailgen

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"strings"
)

type Component interface {
	HTML(tmpl *htmltemplate.Template) (string, error)
	PlainText() string
}

var _ Component = &Action{}
var _ Component = &Line{}
var _ Component = &Table{}

// Action represents an action button in the email message.
type Action struct {
	Text  string
	URL   string
	Color string
}

type Line struct {
	Text string
}

// Table represents a simple table structure for the email message.
type Table struct {
	Headers []TableHeader
	Rows    [][]string
}

// TableHeader represents a header in the table.
type TableHeader struct {
	Text  string
	Width string
	Align string
}

func (a Action) HTML(tmpl *htmltemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "button", a)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (a Action) PlainText() string {
	return fmt.Sprintf("%s ( %s )", a.Text, a.URL)
}

func (l Line) HTML(tmpl *htmltemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "line", l)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (l Line) PlainText() string {
	return l.Text
}

func (t Table) HTML(tmpl *htmltemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "table", t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t Table) PlainText() string {
	var sb strings.Builder

	// Determine column widths
	colWidths := make([]int, len(t.Headers))
	for i, header := range t.Headers {
		colWidths[i] = len(header.Text)
	}
	for _, row := range t.Rows {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	// Write headers
	for i, header := range t.Headers {
		sb.WriteString(t.padCell(header.Text, colWidths[i], header.Align))
		if i < len(t.Headers)-1 {
			sb.WriteString(" | ")
		}
	}
	sb.WriteString("\n")

	// Write separator line
	for i, width := range colWidths {
		sb.WriteString(strings.Repeat("-", width))
		if i < len(colWidths)-1 {
			sb.WriteString("-+-")
		}
	}
	sb.WriteString("\n")

	// Write rows
	for _, row := range t.Rows {
		for i, cell := range row {
			align := t.Headers[i].Align
			sb.WriteString(t.padCell(cell, colWidths[i], align))
			if i < len(row)-1 {
				sb.WriteString(" | ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (t Table) padCell(text string, width int, align string) string {
	switch strings.ToLower(align) {
	case "right":
		return fmt.Sprintf("%*s", width, text)
	case "center":
		padding := width - len(text)
		left := padding / 2
		right := padding - left
		return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
	default: // left
		return fmt.Sprintf("%-*s", width, text)
	}
}
