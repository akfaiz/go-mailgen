package mailgen

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"strconv"
	"strings"
	"unicode"
)

type Component interface {
	HTML(tmpl *htmltemplate.Template) (string, error)
	PlainText() (string, error)
}

var _ Component = &Table{}
var _ Component = &Action{}
var _ Component = &Line{}

type Action struct {
	Text         string
	Link         string
	Color        string
	NoFallback   bool // If true fallbacks are not used
	FallbackText string
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

func (a Action) PlainText() (string, error) {
	return a.Text + " (" + a.Link + ")", nil
}

func (l Line) HTML(tmpl *htmltemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "line", l)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (l Line) PlainText() (string, error) {
	return l.Text, nil
}

func (t Table) HTML(tmpl *htmltemplate.Template) (string, error) {
	var buf bytes.Buffer
	err := tmpl.ExecuteTemplate(&buf, "table", t)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (t Table) PlainText() (string, error) {
	if len(t.Data) == 0 || len(t.Data[0]) == 0 {
		return "", nil
	}

	// Extract column order from first row
	columnNames := make([]string, 0, len(t.Data[0]))
	for _, entry := range t.Data[0] {
		columnNames = append(columnNames, entry.Key)
	}

	// Calculate column widths
	colWidths := make(map[string]int)
	for _, col := range columnNames {
		colWidths[col] = len(col)
		if wStr, ok := t.Columns.CustomWidth[col]; ok {
			if w, err := strconv.Atoi(wStr); err == nil {
				colWidths[col] = w
			}
		}
	}

	// If no custom width, compute max width from data
	for _, row := range t.Data {
		for _, entry := range row {
			width := len(entry.Value)
			if width > colWidths[entry.Key] {
				colWidths[entry.Key] = width
			}
		}
	}

	var sb strings.Builder

	// Header row
	for i, col := range columnNames {
		sb.WriteString(t.padString(t.capitalize(col), colWidths[col], t.Columns.CustomAlign[col]))
		if i < len(columnNames)-1 {
			sb.WriteString(" | ")
		}
	}
	sb.WriteString("\n")

	// Separator row
	for i, col := range columnNames {
		sb.WriteString(strings.Repeat("-", colWidths[col]))
		if i < len(columnNames)-1 {
			sb.WriteString("-+-")
		}
	}
	sb.WriteString("\n")

	// Data rows
	for _, row := range t.Data {
		entryMap := make(map[string]string)
		for _, e := range row {
			entryMap[e.Key] = e.Value
		}
		for i, col := range columnNames {
			val := entryMap[col]
			sb.WriteString(t.padString(val, colWidths[col], t.Columns.CustomAlign[col]))
			if i < len(columnNames)-1 {
				sb.WriteString(" | ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String(), nil
}

func (t Table) padString(s string, width int, align string) string {
	switch align {
	case "right":
		return fmt.Sprintf("%*s", width, s)
	case "center":
		pad := width - len(s)
		left := pad / 2
		right := pad - left
		return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
	default: // left
		return fmt.Sprintf("%-*s", width, s)
	}
}

func (t Table) capitalize(s string) string {
	if s == "" {
		return ""
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
