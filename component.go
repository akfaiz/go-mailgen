package mailgen

import (
	"bytes"
	"fmt"
	htmltemplate "html/template"
	"strconv"
	"strings"
	"unicode"
)

// Component represents a part of the email message, such as a button, line, or table.
type Component interface {
	// HTML generates the HTML representation of the component using the provided template.
	HTML(tmpl *htmltemplate.Template) (string, error)
	// PlainText generates the plain text representation of the component.
	PlainText() (string, error)
}

var _ Component = &Table{}
var _ Component = &Action{}
var _ Component = &Line{}

// Action represents a button or link in the email.
type Action struct {
	// Text is the text displayed on the button.
	Text string
	// Link is the URL the button points to.
	Link string
	// Color is hex color code for the button, e.g. "#3869D4".
	Color string
	// NoFallback if true, the action will not have a fallback text.
	NoFallback   bool
	FallbackText string
}

// Line represents a simple text line in the email.
type Line struct {
	Text string
}

// Table represents a structured table in the email.
// It contains data entries and column definitions.
//
// Example usage:
//
//	table := mailgen.Table{
//	    Data: [][]mailgen.Entry{
//	        {
//	            {"Key": "name", "Value": "John Doe"},
//	            {"Key": "email", "Value": "john@example.com"},
//	        },
//	    },
//	    Columns: mailgen.Columns{
//	        CustomWidth: map[string]string{
//	            "name":  "200px",
//	            "email": "300px",
//	        },
//	        CustomAlign: map[string]string{
//	            "name":  "left",
//	            "email": "right",
//	        },
//	    },
//	}
type Table struct {
	// Data contains the rows of the table, each row is a slice of Entry.
	// Each Entry has a Key and Value, where Key is the column name.
	Data [][]Entry
	// Columns defines column properties like width and alignment.
	Columns Columns
}

// Entry represents a single entry in the table with a key and value.
type Entry struct {
	Key   string
	Value string
}

// Columns defines the structure of the table columns.
type Columns struct {
	// CustomWidth allows setting specific widths for columns.
	CustomWidth map[string]string
	// CustomAlign allows setting specific alignments for columns.
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
