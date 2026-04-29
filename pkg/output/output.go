// Package output 提供统一的 table / JSON 渲染入口。
//
// 命令层只调 Render(w, data, format, columns)；
// table 模式要求 data 为 []any，命令层负责把领域对象切片转为 []any；
// json 模式直接 encoding/json 编码 data，便于自动化下游消费。
package output

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/olekukonko/tablewriter"
)

// Column 描述一列：表头名 + 取值函数（接收单行的 any 值）。
type Column struct {
	Header string
	Get    func(any) string
}

// Render 是统一渲染入口。
// format 为 "" 等价于 "table"。未知 format 返回错误。
func Render(w io.Writer, data any, format string, columns []Column) error {
	switch format {
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		if err := enc.Encode(data); err != nil {
			return fmt.Errorf("encode json: %w", err)
		}
		return nil
	case "table", "":
		return renderTable(w, data, columns)
	default:
		return fmt.Errorf("unsupported format %q (want table|json)", format)
	}
}

// renderTable 渲染表格。data 必须是 []any（命令层负责转换）。
// 即使 rows 为空，也会输出表头，便于自动化判定"无结果"。
func renderTable(w io.Writer, data any, columns []Column) error {
	rows, ok := data.([]any)
	if !ok {
		return fmt.Errorf("table render expects []any, got %T", data)
	}
	tw := tablewriter.NewWriter(w)
	headers := make([]any, len(columns))
	for i, c := range columns {
		headers[i] = c.Header
	}
	tw.Header(headers...)
	for _, item := range rows {
		row := make([]any, len(columns))
		for i, c := range columns {
			row[i] = c.Get(item)
		}
		if err := tw.Append(row...); err != nil {
			return fmt.Errorf("append row: %w", err)
		}
	}
	if err := tw.Render(); err != nil {
		return fmt.Errorf("render table: %w", err)
	}
	return nil
}
