package mailgen_test

import (
	htmltemplate "html/template"
	"testing"

	"github.com/ahmadfaizk/go-mailgen"
	"github.com/stretchr/testify/assert"
)

func TestAction_HTML(t *testing.T) {
	tests := []struct {
		name     string
		action   mailgen.Action
		template string
		want     string
		wantErr  bool
	}{
		{
			name: "valid action with template",
			action: mailgen.Action{
				Text:  "Click Me",
				URL:   "https://example.com",
				Color: "blue",
			},
			template: `{{define "button"}}<a href="{{.URL}}" style="color:{{.Color}}">{{.Text}}</a>{{end}}`,
			want:     `<a href="https://example.com" style="color:blue">Click Me</a>`,
			wantErr:  false,
		},
		{
			name: "action with empty values",
			action: mailgen.Action{
				Text:  "",
				URL:   "",
				Color: "",
			},
			template: `{{define "button"}}<a href="{{.URL}}" style="color:{{.Color}}">{{.Text}}</a>{{end}}`,
			want:     `<a href="" style="color:"></a>`,
			wantErr:  false,
		},
		{
			name: "template execution error - missing button template",
			action: mailgen.Action{
				Text: "Click Me",
				URL:  "https://example.com",
			},
			template: `{{define "notbutton"}}invalid{{end}}`,
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := htmltemplate.New("test").Parse(tt.template)
			assert.NoError(t, err, "Failed to parse template")

			got, err := tt.action.HTML(tmpl)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestLine_HTML(t *testing.T) {
	tests := []struct {
		name     string
		line     mailgen.Line
		template string
		want     string
		wantErr  bool
	}{
		{
			name: "valid line with template",
			line: mailgen.Line{
				Text: "This is a line of text",
			},
			template: `{{define "line"}}<p>{{.Text}}</p>{{end}}`,
			want:     `<p>This is a line of text</p>`,
			wantErr:  false,
		},
		{
			name: "line with empty text",
			line: mailgen.Line{
				Text: "",
			},
			template: `{{define "line"}}<p>{{.Text}}</p>{{end}}`,
			want:     `<p></p>`,
			wantErr:  false,
		},
		{
			name: "line with special characters",
			line: mailgen.Line{
				Text: "Hello & welcome! <em>Enjoy</em>",
			},
			template: `{{define "line"}}<div>{{.Text}}</div>{{end}}`,
			want:     `<div>Hello &amp; welcome! &lt;em&gt;Enjoy&lt;/em&gt;</div>`,
			wantErr:  false,
		},
		{
			name: "template execution error - missing line template",
			line: mailgen.Line{
				Text: "Some text",
			},
			template: `{{define "notline"}}invalid{{end}}`,
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := htmltemplate.New("test").Parse(tt.template)
			assert.NoError(t, err, "Failed to parse template")

			got, err := tt.line.HTML(tmpl)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestTable_HTML(t *testing.T) {
	tests := []struct {
		name     string
		table    mailgen.Table
		template string
		want     string
		wantErr  bool
	}{
		{
			name: "valid table with headers and rows",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Name", Width: "50%", Align: "left"},
					{Text: "Age", Width: "25%", Align: "right"},
					{Text: "City", Width: "25%", Align: "center"},
				},
				Rows: [][]string{
					{"John", "25", "New York"},
					{"Jane", "30", "Los Angeles"},
				},
			},
			template: `{{define "table"}}<table>{{range .Headers}}<th>{{.Text}}</th>{{end}}{{range .Rows}}{{range .}}<td>{{.}}</td>{{end}}{{end}}</table>{{end}}`,
			want:     `<table><th>Name</th><th>Age</th><th>City</th><td>John</td><td>25</td><td>New York</td><td>Jane</td><td>30</td><td>Los Angeles</td></table>`,
			wantErr:  false,
		},
		{
			name: "empty table",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{},
				Rows:    [][]string{},
			},
			template: `{{define "table"}}<table>{{range .Headers}}<th>{{.Text}}</th>{{end}}{{range .Rows}}{{range .}}<td>{{.}}</td>{{end}}{{end}}</table>{{end}}`,
			want:     `<table></table>`,
			wantErr:  false,
		},
		{
			name: "table with only headers",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Header1", Width: "100%", Align: "left"},
				},
				Rows: [][]string{},
			},
			template: `{{define "table"}}<table>{{range .Headers}}<th>{{.Text}}</th>{{end}}{{range .Rows}}{{range .}}<td>{{.}}</td>{{end}}{{end}}</table>{{end}}`,
			want:     `<table><th>Header1</th></table>`,
			wantErr:  false,
		},
		{
			name: "table with empty cells",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Col1", Width: "", Align: ""},
					{Text: "Col2", Width: "", Align: ""},
				},
				Rows: [][]string{
					{"", "value"},
					{"data", ""},
				},
			},
			template: `{{define "table"}}<table>{{range .Headers}}<th>{{.Text}}</th>{{end}}{{range .Rows}}{{range .}}<td>{{.}}</td>{{end}}{{end}}</table>{{end}}`,
			want:     `<table><th>Col1</th><th>Col2</th><td></td><td>value</td><td>data</td><td></td></table>`,
			wantErr:  false,
		},
		{
			name: "template execution error - missing table template",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Header", Width: "100%", Align: "left"},
				},
				Rows: [][]string{
					{"Data"},
				},
			},
			template: `{{define "nottable"}}invalid{{end}}`,
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := htmltemplate.New("test").Parse(tt.template)
			assert.NoError(t, err, "Failed to parse template")

			got, err := tt.table.HTML(tmpl)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestAction_PlainText(t *testing.T) {
	tests := []struct {
		name   string
		action mailgen.Action
		want   string
	}{
		{
			name: "valid action",
			action: mailgen.Action{
				Text: "Click Me",
				URL:  "https://example.com",
			},
			want: "Click Me ( https://example.com )",
		},
		{
			name: "action with empty text",
			action: mailgen.Action{
				Text: "",
				URL:  "https://example.com",
			},
			want: " ( https://example.com )",
		},
		{
			name: "action with empty URL",
			action: mailgen.Action{
				Text: "Click Me",
				URL:  "",
			},
			want: "Click Me (  )",
		},
		{
			name: "action with both empty",
			action: mailgen.Action{
				Text: "",
				URL:  "",
			},
			want: " (  )",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.action.PlainText()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLine_PlainText(t *testing.T) {
	tests := []struct {
		name string
		line mailgen.Line
		want string
	}{
		{
			name: "valid line",
			line: mailgen.Line{
				Text: "This is a line of text",
			},
			want: "This is a line of text",
		},
		{
			name: "empty line",
			line: mailgen.Line{
				Text: "",
			},
			want: "",
		},
		{
			name: "line with special characters",
			line: mailgen.Line{
				Text: "Hello & welcome! <em>Enjoy</em>",
			},
			want: "Hello & welcome! <em>Enjoy</em>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.line.PlainText()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTable_PlainText(t *testing.T) {
	tests := []struct {
		name  string
		table mailgen.Table
		want  string
	}{
		{
			name: "simple table with left alignment",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Name", Align: "left"},
					{Text: "Age", Align: "left"},
				},
				Rows: [][]string{
					{"John", "25"},
					{"Jane", "30"},
				},
			},
			want: "Name | Age\n-----+----\nJohn | 25 \nJane | 30 \n",
		},
		{
			name: "table with different alignments",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Name", Align: "left"},
					{Text: "Age", Align: "right"},
					{Text: "City", Align: "center"},
				},
				Rows: [][]string{
					{"John", "25", "NYC"},
					{"Jane", "30", "LA"},
				},
			},
			want: "Name | Age | City\n-----+-----+-----\nJohn |  25 | NYC \nJane |  30 |  LA \n",
		},
		{
			name: "table with varying column widths",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Name", Align: "left"},
					{Text: "Description", Align: "left"},
				},
				Rows: [][]string{
					{"John", "Software Engineer"},
					{"Jane", "Designer"},
				},
			},
			want: "Name | Description      \n-----+------------------\nJohn | Software Engineer\nJane | Designer         \n",
		},
		{
			name: "table with empty cells",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Col1", Align: "left"},
					{Text: "Col2", Align: "left"},
				},
				Rows: [][]string{
					{"", "value"},
					{"data", ""},
				},
			},
			want: "Col1 | Col2 \n-----+------\n     | value\ndata |      \n",
		},
		{
			name: "table with center alignment",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Short", Align: "center"},
					{Text: "LongerHeader", Align: "center"},
				},
				Rows: [][]string{
					{"A", "B"},
					{"Test", "LongValue"},
				},
			},
			want: "Short | LongerHeader\n------+-------------\n  A   |      B      \nTest  |  LongValue  \n",
		},
		{
			name: "empty table",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{},
				Rows:    [][]string{},
			},
			want: "\n\n",
		},
		{
			name: "table with only headers",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Header1", Align: "left"},
					{Text: "Header2", Align: "left"},
				},
				Rows: [][]string{},
			},
			want: "Header1 | Header2\n--------+--------\n",
		},
		{
			name: "single column table",
			table: mailgen.Table{
				Headers: []mailgen.TableHeader{
					{Text: "Column", Align: "left"},
				},
				Rows: [][]string{
					{"Value1"},
					{"Value2"},
				},
			},
			want: "Column\n------\nValue1\nValue2\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.table.PlainText()
			assert.Equal(t, tt.want, got)
		})
	}
}
