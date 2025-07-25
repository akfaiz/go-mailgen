package templates

import (
	"embed"
	"fmt"
	htmltemplate "html/template"
	"strconv"
	"strings"
	texttemplate "text/template"
	"unicode"

	"github.com/afkdevs/go-mailgen/component"
)

//go:embed default/*
var defaultFS embed.FS

var DefaultHtmlTmpl = htmltemplate.Must(htmltemplate.New("index.html").
	Funcs(htmlTemplateFuncs).
	ParseFS(defaultFS, "default/*.html"),
)

//go:embed plain/*
var plainFS embed.FS

var PlainHtmlTmpl = htmltemplate.Must(htmltemplate.New("index.html").
	Funcs(htmlTemplateFuncs).
	ParseFS(plainFS, "plain/*.html"),
)

//go:embed plaintext/*
var plaintextFS embed.FS

var PlainTextTmpl = texttemplate.Must(texttemplate.New("index.txt").
	Funcs(textTemplateFuncs).
	ParseFS(plaintextFS, "plaintext/*.txt"),
)

var htmlTemplateFuncs = htmltemplate.FuncMap{
	"capitalize": capitalize,
}
var textTemplateFuncs = texttemplate.FuncMap{
	"boxString":     boxString,
	"concat":        concat,
	"generateTable": generateTable,
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

func generateTable(t component.Table) string {
	if len(t.Data) == 0 || len(t.Data[0]) == 0 {
		return ""
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
		sb.WriteString(padString(capitalize(col), colWidths[col], t.Columns.CustomAlign[col]))
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
			sb.WriteString(padString(val, colWidths[col], t.Columns.CustomAlign[col]))
			if i < len(columnNames)-1 {
				sb.WriteString(" | ")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func padString(s string, width int, align string) string {
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
