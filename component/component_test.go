package component_test

import (
	htmltemplate "html/template"
	"testing"
	texttemplate "text/template"

	"github.com/ahmadfaizk/go-mailgen/component"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAction_HTML(t *testing.T) {
	tests := []struct {
		name     string
		action   component.Action
		template string
		expected string
		wantErr  bool
	}{
		{
			name: "valid action with all fields",
			action: component.Action{
				Text:  "Click Me",
				Link:  "https://example.com",
				Color: "blue",
			},
			template: `{{define "button"}}<a href="{{.Link}}" style="color:{{.Color}}">{{.Text}}</a>{{end}}`,
			expected: `<a href="https://example.com" style="color:blue">Click Me</a>`,
			wantErr:  false,
		},
		{
			name: "action with empty fields",
			action: component.Action{
				Text:  "",
				Link:  "",
				Color: "",
			},
			template: `{{define "button"}}<a href="{{.Link}}" style="color:{{.Color}}">{{.Text}}</a>{{end}}`,
			expected: `<a href="" style="color:"></a>`,
			wantErr:  false,
		},
		{
			name: "template execution error",
			action: component.Action{
				Text: "Test",
			},
			template: `{{define "button"}}{{.InvalidField}}{{end}}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := htmltemplate.New("test").Parse(tt.template)
			require.NoError(t, err)

			result, err := tt.action.HTML(tmpl)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestAction_PlainText(t *testing.T) {
	tests := []struct {
		name     string
		action   component.Action
		template string
		expected string
		wantErr  bool
	}{
		{
			name: "valid action with all fields",
			action: component.Action{
				Text:  "Click Me",
				Link:  "https://example.com",
				Color: "blue",
			},
			template: `{{define "button"}}{{.Text}} - {{.Link}}{{end}}`,
			expected: `Click Me - https://example.com`,
			wantErr:  false,
		},
		{
			name: "action with empty fields",
			action: component.Action{
				Text:  "",
				Link:  "",
				Color: "",
			},
			template: `{{define "button"}}{{.Text}} - {{.Link}}{{end}}`,
			expected: ` - `,
			wantErr:  false,
		},
		{
			name: "simple text template",
			action: component.Action{
				Text: "Download Report",
			},
			template: `{{define "button"}}{{.Text}}{{end}}`,
			expected: `Download Report`,
			wantErr:  false,
		},
		{
			name: "template execution error",
			action: component.Action{
				Text: "Test",
			},
			template: `{{define "button"}}{{.InvalidField}}{{end}}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := texttemplate.New("test").Parse(tt.template)
			require.NoError(t, err)

			result, err := tt.action.PlainText(tmpl)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestLine_HTML(t *testing.T) {
	tests := []struct {
		name     string
		line     component.Line
		template string
		expected string
		wantErr  bool
	}{
		{
			name: "valid line with text",
			line: component.Line{
				Text: "This is a line of text",
			},
			template: `{{define "line"}}<p>{{.Text}}</p>{{end}}`,
			expected: `<p>This is a line of text</p>`,
			wantErr:  false,
		},
		{
			name: "line with empty text",
			line: component.Line{
				Text: "",
			},
			template: `{{define "line"}}<p>{{.Text}}</p>{{end}}`,
			expected: `<p></p>`,
			wantErr:  false,
		},
		{
			name: "template execution error",
			line: component.Line{
				Text: "Test",
			},
			template: `{{define "line"}}{{.InvalidField}}{{end}}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := htmltemplate.New("test").Parse(tt.template)
			require.NoError(t, err)

			result, err := tt.line.HTML(tmpl)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestLine_PlainText(t *testing.T) {
	tests := []struct {
		name     string
		line     component.Line
		template string
		expected string
		wantErr  bool
	}{
		{
			name: "valid line with text",
			line: component.Line{
				Text: "This is a line of text",
			},
			template: `{{define "line"}}{{.Text}}{{end}}`,
			expected: `This is a line of text`,
			wantErr:  false,
		},
		{
			name: "line with empty text",
			line: component.Line{
				Text: "",
			},
			template: `{{define "line"}}{{.Text}}{{end}}`,
			expected: ``,
			wantErr:  false,
		},
		{
			name: "line with multiline text",
			line: component.Line{
				Text: "Line 1\nLine 2\nLine 3",
			},
			template: `{{define "line"}}{{.Text}}{{end}}`,
			expected: `Line 1
Line 2
Line 3`,
			wantErr:  false,
		},
		{
			name: "template execution error",
			line: component.Line{
				Text: "Test",
			},
			template: `{{define "line"}}{{.InvalidField}}{{end}}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := texttemplate.New("test").Parse(tt.template)
			require.NoError(t, err)

			result, err := tt.line.PlainText(tmpl)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestTable_HTML(t *testing.T) {
	tests := []struct {
		name     string
		table    component.Table
		template string
		expected string
		wantErr  bool
	}{
		{
			name: "valid table with data",
			table: component.Table{
				Data: [][]component.Entry{
					{
						{Key: "Name", Value: "John"},
						{Key: "Age", Value: "30"},
					},
					{
						{Key: "Name", Value: "Jane"},
						{Key: "Age", Value: "25"},
					},
				},
				Columns: component.Columns{
					CustomWidth: map[string]string{"Name": "50%"},
					CustomAlign: map[string]string{"Age": "center"},
				},
			},
			template: `{{define "table"}}<table>{{range .Data}}<tr>{{range .}}<td>{{.Value}}</td>{{end}}</tr>{{end}}</table>{{end}}`,
			expected: `<table><tr><td>John</td><td>30</td></tr><tr><td>Jane</td><td>25</td></tr></table>`,
			wantErr:  false,
		},
		{
			name: "empty table",
			table: component.Table{
				Data:    [][]component.Entry{},
				Columns: component.Columns{},
			},
			template: `{{define "table"}}<table>{{range .Data}}<tr>{{range .}}<td>{{.Value}}</td>{{end}}</tr>{{end}}</table>{{end}}`,
			expected: `<table></table>`,
			wantErr:  false,
		},
		{
			name: "table with custom width and align",
			table: component.Table{
				Data: [][]component.Entry{
					{
						{Key: "Header1", Value: "Value1"},
					},
				},
				Columns: component.Columns{
					CustomWidth: map[string]string{"Header1": "100px"},
					CustomAlign: map[string]string{"Header1": "left"},
				},
			},
			template: `{{define "table"}}<table>{{range .Data}}<tr>{{range .}}<td style="width:{{index $.Columns.CustomWidth .Key}};text-align:{{index $.Columns.CustomAlign .Key}}">{{.Value}}</td>{{end}}</tr>{{end}}</table>{{end}}`,
			expected: `<table><tr><td style="width:100px;text-align:left">Value1</td></tr></table>`,
			wantErr:  false,
		},
		{
			name: "table with nil columns",
			table: component.Table{
				Data: [][]component.Entry{
					{
						{Key: "Test", Value: "Data"},
					},
				},
			},
			template: `{{define "table"}}<table>{{range .Data}}<tr>{{range .}}<td>{{.Key}}: {{.Value}}</td>{{end}}</tr>{{end}}</table>{{end}}`,
			expected: `<table><tr><td>Test: Data</td></tr></table>`,
			wantErr:  false,
		},
		{
			name: "template execution error",
			table: component.Table{
				Data: [][]component.Entry{
					{{Key: "Test", Value: "Data"}},
				},
			},
			template: `{{define "table"}}{{.InvalidField}}{{end}}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := htmltemplate.New("test").Parse(tt.template)
			require.NoError(t, err)

			result, err := tt.table.HTML(tmpl)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestTable_PlainText(t *testing.T) {
	tests := []struct {
		name     string
		table    component.Table
		template string
		expected string
		wantErr  bool
	}{
		{
			name: "valid table with data",
			table: component.Table{
				Data: [][]component.Entry{
					{
						{Key: "Name", Value: "John"},
						{Key: "Age", Value: "30"},
					},
					{
						{Key: "Name", Value: "Jane"},
						{Key: "Age", Value: "25"},
					},
				},
				Columns: component.Columns{
					CustomWidth: map[string]string{"Name": "50%"},
					CustomAlign: map[string]string{"Age": "center"},
				},
			},
			template: `{{define "table"}}{{range .Data}}{{range .}}{{.Key}}: {{.Value}} {{end}}
{{end}}{{end}}`,
			expected: `Name: John Age: 30 
Name: Jane Age: 25 
`,
			wantErr: false,
		},
		{
			name: "empty table",
			table: component.Table{
				Data:    [][]component.Entry{},
				Columns: component.Columns{},
			},
			template: `{{define "table"}}{{range .Data}}{{range .}}{{.Value}} {{end}}{{end}}{{end}}`,
			expected: ``,
			wantErr:  false,
		},
		{
			name: "table with single row",
			table: component.Table{
				Data: [][]component.Entry{
					{
						{Key: "Header1", Value: "Value1"},
						{Key: "Header2", Value: "Value2"},
					},
				},
				Columns: component.Columns{
					CustomWidth: map[string]string{"Header1": "100px"},
					CustomAlign: map[string]string{"Header1": "left"},
				},
			},
			template: `{{define "table"}}{{range .Data}}{{range .}}{{.Value}} {{end}}{{end}}{{end}}`,
			expected: `Value1 Value2 `,
			wantErr:  false,
		},
		{
			name: "table with formatted output",
			table: component.Table{
				Data: [][]component.Entry{
					{
						{Key: "Product", Value: "Laptop"},
						{Key: "Price", Value: "$999"},
					},
					{
						{Key: "Product", Value: "Mouse"},
						{Key: "Price", Value: "$25"},
					},
				},
			},
			template: `{{define "table"}}{{range .Data}}{{range .}}{{.Key}}: {{.Value}}
{{end}}---
{{end}}{{end}}`,
			expected: `Product: Laptop
Price: $999
---
Product: Mouse
Price: $25
---
`,
			wantErr: false,
		},
		{
			name: "table with nil columns",
			table: component.Table{
				Data: [][]component.Entry{
					{
						{Key: "Test", Value: "Data"},
					},
				},
			},
			template: `{{define "table"}}{{range .Data}}{{range .}}{{.Key}}: {{.Value}}{{end}}{{end}}{{end}}`,
			expected: `Test: Data`,
			wantErr:  false,
		},
		{
			name: "template execution error",
			table: component.Table{
				Data: [][]component.Entry{
					{{Key: "Test", Value: "Data"}},
				},
			},
			template: `{{define "table"}}{{.InvalidField}}{{end}}`,
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := texttemplate.New("test").Parse(tt.template)
			require.NoError(t, err)

			result, err := tt.table.PlainText(tmpl)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}



