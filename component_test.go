package mailgen_test

import (
	htmltemplate "html/template"
	"testing"

	"github.com/afkdevs/go-mailgen"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAction_HTML(t *testing.T) {
	tests := []struct {
		name     string
		action   mailgen.Action
		template string
		expected string
		wantErr  bool
	}{
		{
			name: "valid action with all fields",
			action: mailgen.Action{
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
			action: mailgen.Action{
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
			action: mailgen.Action{
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

func TestLine_HTML(t *testing.T) {
	tests := []struct {
		name     string
		line     mailgen.Line
		template string
		expected string
		wantErr  bool
	}{
		{
			name: "valid line with text",
			line: mailgen.Line{
				Text: "This is a line of text",
			},
			template: `{{define "line"}}<p>{{.Text}}</p>{{end}}`,
			expected: `<p>This is a line of text</p>`,
			wantErr:  false,
		},
		{
			name: "line with empty text",
			line: mailgen.Line{
				Text: "",
			},
			template: `{{define "line"}}<p>{{.Text}}</p>{{end}}`,
			expected: `<p></p>`,
			wantErr:  false,
		},
		{
			name: "template execution error",
			line: mailgen.Line{
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

func TestTable_HTML(t *testing.T) {
	tests := []struct {
		name     string
		table    mailgen.Table
		template string
		expected string
		wantErr  bool
	}{
		{
			name: "valid table with data",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "Name", Value: "John"},
						{Key: "Age", Value: "30"},
					},
					{
						{Key: "Name", Value: "Jane"},
						{Key: "Age", Value: "25"},
					},
				},
				Columns: mailgen.Columns{
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
			table: mailgen.Table{
				Data:    [][]mailgen.Entry{},
				Columns: mailgen.Columns{},
			},
			template: `{{define "table"}}<table>{{range .Data}}<tr>{{range .}}<td>{{.Value}}</td>{{end}}</tr>{{end}}</table>{{end}}`,
			expected: `<table></table>`,
			wantErr:  false,
		},
		{
			name: "table with custom width and align",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "Header1", Value: "Value1"},
					},
				},
				Columns: mailgen.Columns{
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
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
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
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
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
		table    mailgen.Table
		expected string
		wantErr  bool
	}{
		{
			name: "simple table with two columns",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "name", Value: "John"},
						{Key: "age", Value: "30"},
					},
					{
						{Key: "name", Value: "Jane"},
						{Key: "age", Value: "25"},
					},
				},
			},
			expected: "Name | Age\n-----+----\nJohn | 30 \nJane | 25 \n",
			wantErr:  false,
		},
		{
			name: "empty table data",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{},
			},
			expected: "",
			wantErr:  false,
		},
		{
			name: "table with empty first row",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{},
				},
			},
			expected: "",
			wantErr:  false,
		},
		{
			name: "table with custom width",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "id", Value: "1"},
						{Key: "name", Value: "John"},
					},
					{
						{Key: "id", Value: "2"},
						{Key: "name", Value: "Jane"},
					},
				},
				Columns: mailgen.Columns{
					CustomWidth: map[string]string{
						"id":   "5",
						"name": "10",
					},
				},
			},
			expected: "Id    | Name      \n------+-----------\n1     | John      \n2     | Jane      \n",
			wantErr:  false,
		},
		{
			name: "table with custom alignment",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "name", Value: "John"},
						{Key: "score", Value: "95"},
					},
					{
						{Key: "name", Value: "Jane"},
						{Key: "score", Value: "87"},
					},
				},
				Columns: mailgen.Columns{
					CustomAlign: map[string]string{
						"name":  "left",
						"score": "right",
					},
				},
			},
			expected: "Name | Score\n-----+------\nJohn |    95\nJane |    87\n",
			wantErr:  false,
		},
		{
			name: "table with center alignment",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "status", Value: "ok"},
					},
					{
						{Key: "status", Value: "fail"},
					},
				},
				Columns: mailgen.Columns{
					CustomWidth: map[string]string{
						"status": "8",
					},
					CustomAlign: map[string]string{
						"status": "center",
					},
				},
			},
			expected: " Status \n--------\n   ok   \n  fail  \n",
			wantErr:  false,
		},
		{
			name: "table with varying column widths",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "short", Value: "a"},
						{Key: "long", Value: "very long value"},
					},
					{
						{Key: "short", Value: "b"},
						{Key: "long", Value: "x"},
					},
				},
			},
			expected: "Short | Long           \n------+----------------\na     | very long value\nb     | x              \n",
			wantErr:  false,
		},
		{
			name: "table with invalid custom width",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "test", Value: "value"},
					},
				},
				Columns: mailgen.Columns{
					CustomWidth: map[string]string{
						"test": "invalid",
					},
				},
			},
			expected: "Test \n-----\nvalue\n",
			wantErr:  false,
		},
		{
			name: "single row table",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "column", Value: "data"},
					},
				},
			},
			expected: "Column\n------\ndata  \n",
			wantErr:  false,
		},
		{
			name: "table with missing values in subsequent rows",
			table: mailgen.Table{
				Data: [][]mailgen.Entry{
					{
						{Key: "a", Value: "1"},
						{Key: "b", Value: "2"},
					},
					{
						{Key: "a", Value: "3"},
					},
				},
			},
			expected: "A | B\n--+--\n1 | 2\n3 |  \n",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.table.PlainText()

			if tt.wantErr {
				assert.Error(t, err)
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
		action   mailgen.Action
		expected string
		wantErr  bool
	}{
		{
			name: "action with text and link",
			action: mailgen.Action{
				Text: "Click Here",
				Link: "https://example.com",
			},
			expected: "Click Here (https://example.com)",
			wantErr:  false,
		},
		{
			name: "action with empty text",
			action: mailgen.Action{
				Text: "",
				Link: "https://example.com",
			},
			expected: " (https://example.com)",
			wantErr:  false,
		},
		{
			name: "action with empty link",
			action: mailgen.Action{
				Text: "Click Here",
				Link: "",
			},
			expected: "Click Here ()",
			wantErr:  false,
		},
		{
			name: "action with both empty",
			action: mailgen.Action{
				Text: "",
				Link: "",
			},
			expected: " ()",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.action.PlainText()

			if tt.wantErr {
				assert.Error(t, err)
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
		line     mailgen.Line
		expected string
		wantErr  bool
	}{
		{
			name: "line with text",
			line: mailgen.Line{
				Text: "This is a line of text",
			},
			expected: "This is a line of text",
			wantErr:  false,
		},
		{
			name: "line with empty text",
			line: mailgen.Line{
				Text: "",
			},
			expected: "",
			wantErr:  false,
		},
		{
			name: "line with multiline text",
			line: mailgen.Line{
				Text: "Line 1\nLine 2",
			},
			expected: "Line 1\nLine 2",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.line.PlainText()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
