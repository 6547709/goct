package find

import (
	"github.com/6547709/goct/pkg/output"
	"io"
)

// writeTable 渲染 find 的结果到表格。
func writeTable(w io.Writer, rows []row) error {
	items := make([]any, len(rows))
	for i := range rows {
		items[i] = rows[i]
	}
	cols := []output.Column{
		{Header: "TYPE", Get: func(v any) string { return v.(row).Type }},
		{Header: "ID", Get: func(v any) string { return v.(row).ID }},
		{Header: "NAME", Get: func(v any) string { return v.(row).Name }},
	}
	return output.Render(w, items, "table", cols)
}

// writeJSON 渲染为 JSON 数组。
func writeJSON(w io.Writer, rows []row) error {
	return output.Render(w, rows, "json", nil)
}
