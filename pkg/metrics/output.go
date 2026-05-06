package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// OutputFormat 支持的三种输出格式。
const (
	FormatTable = "table"
	FormatJSON  = "json"
	FormatChart = "chart"
)

// RenderResult 将 MetricResult 列表渲染为指定格式。
func RenderResult(w io.Writer, results []MetricResult, format string, latest bool) error {
	switch format {
	case FormatJSON:
		return renderJSON(w, results)
	case FormatChart:
		return renderChart(w, results)
	default:
		return renderTable(w, results, latest)
	}
}

func renderJSON(w io.Writer, results []MetricResult) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func renderTable(w io.Writer, results []MetricResult, latest bool) error {
	tw := tablewriter.NewWriter(w)
	tw.Header([]any{"TIME", "METRIC", "VALUE", "UNIT"}...)

	for _, r := range results {
		if latest && r.Latest != nil {
			tw.Append([]any{
				r.Latest.Timestamp,
				r.Metric,
				fmt.Sprintf("%.2f", r.Latest.Value),
				r.Latest.Unit,
			})
		} else {
			for _, s := range r.Samples {
				tw.Append([]any{
					s.Timestamp,
					r.Metric,
					fmt.Sprintf("%.2f", s.Value),
					s.Unit,
				})
			}
		}
	}
	tw.Render()
	return nil
}

func renderChart(w io.Writer, results []MetricResult) error {
	for _, r := range results {
		if len(r.Samples) == 0 {
			continue
		}
		fmt.Fprintf(w, "\n%s (%s)\n", r.Metric, r.Target)

		// 找最大最小值
		var min, max float64
		for i, s := range r.Samples {
			if i == 0 || s.Value < min {
				min = s.Value
			}
			if s.Value > max {
				max = s.Value
			}
		}

		// 渲染 ASCII 图表
		barWidth := 50
		range_ := max - min
		for _, s := range r.Samples {
			var barLen int
			if range_ > 0 {
				barLen = int((s.Value - min) / range_ * float64(barWidth))
			}
			bar := strings.Repeat("█", barLen)
			fmt.Fprintf(w, "%5.1f|%s%s %.2f%%\n", s.Value, bar,
				strings.Repeat(" ", barWidth-barLen), s.Value)
		}
		fmt.Fprintln(w)
	}
	return nil
}