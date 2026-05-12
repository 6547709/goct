// Package events 实现 govc 风格的 `goct events` 命令。
//
// CloudTower 没有"事件流"概念，最接近的实现是 user_audit_log（用户审计日志），
// 由 adapter.EventOps 暴露。本命令把它包装成 govc events 风格的输出。
package events

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/spf13/cobra"
)

const groupID = "task"

// Register 注册命令到 root。
func Register(root *cobra.Command) {
	root.AddCommand(newEvents())
}

// eventColumns 是表格输出列，与 govc events 的 "Date / User / Type / Target / Message" 对齐。
var eventColumns = []output.Column{
	{Header: "DATE", Get: func(v any) string { return shortTime(v.(adapter.Event).CreatedAt) }},
	{Header: "USER", Get: func(v any) string { return v.(adapter.Event).Username }},
	{Header: "ACTION", Get: func(v any) string { return v.(adapter.Event).Action }},
	{Header: "STATUS", Get: func(v any) string { return v.(adapter.Event).Status }},
	{Header: "TARGET", Get: func(v any) string {
		ev := v.(adapter.Event)
		if ev.ResourceType == "" && ev.ResourceID == "" {
			return ""
		}
		return ev.ResourceType + ":" + ev.ResourceID
	}},
	{Header: "MESSAGE", Get: func(v any) string { return truncate(v.(adapter.Event).Message, 80) }},
}

func newEvents() *cobra.Command {
	var (
		resourceID, resourceType, user, action string
		limit                                  int32
		follow                                 bool
		formatJSON                             bool
	)
	c := &cobra.Command{
		Use:     "events [resource-id]",
		Short:   "Show recent events / audit log entries",
		GroupID: groupID,
		Long: `Show CloudTower events (backed by user_audit_log).

Examples:
  goct events                                    # last 50 events, newest first
  goct events --limit 200
  goct events --user admin
  goct events --action create_vm
  goct events --type VM
  goct events vm-uuid                            # events for one resource
  goct events --follow                           # tail (refresh every 2s)`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			if len(args) > 0 && resourceID == "" {
				resourceID = args[0]
			}
			opts := adapter.EventListOpts{
				ResourceID:   resourceID,
				ResourceType: strings.ToUpper(strings.TrimSpace(resourceType)),
				Username:     user,
				ActionLike:   action,
				Limit:        limit,
			}

			if !follow {
				items, err := cli.ListEvents(c.Context(), opts)
				if err != nil {
					return err
				}
				return renderEvents(c.OutOrStdout(), items, formatJSON)
			}

			// follow 模式：每 2s 拉一次，去重已展示过的事件 ID（按 ID 集合维护）。
			seen := map[string]bool{}
			ticker := time.NewTicker(2 * time.Second)
			defer ticker.Stop()
			for {
				items, err := cli.ListEvents(c.Context(), opts)
				if err != nil {
					return err
				}
				// 反向遍历，按时间正序追加新条目。
				for i := len(items) - 1; i >= 0; i-- {
					ev := items[i]
					if seen[ev.ID] {
						continue
					}
					seen[ev.ID] = true
					if err := writeOneLine(c.OutOrStdout(), ev, formatJSON); err != nil {
						return err
					}
				}
				select {
				case <-c.Context().Done():
					return nil
				case <-ticker.C:
				}
			}
		},
	}
	c.Flags().StringVar(&resourceID, "resource", "", "Filter by resource ID (positional arg also accepted)")
	c.Flags().StringVar(&resourceType, "type", "", "Filter by resource type (VM/HOST/CLUSTER/...)")
	c.Flags().StringVar(&user, "user", "", "Filter by triggering username")
	c.Flags().StringVar(&action, "action", "", "Filter by action substring (e.g. create_vm)")
	c.Flags().Int32VarP(&limit, "limit", "n", 50, "Max events to show (ignored in --follow)")
	c.Flags().BoolVarP(&follow, "follow", "f", false, "Follow new events (poll every 2s)")
	c.Flags().BoolVar(&formatJSON, "json", false, "JSON output")
	return c
}

func renderEvents(w io.Writer, items []adapter.Event, asJSON bool) error {
	if asJSON {
		return output.Render(w, items, "json", nil)
	}
	rows := make([]any, len(items))
	for i := range items {
		rows[i] = items[i]
	}
	return output.Render(w, rows, "table", eventColumns)
}

// writeOneLine 用于 --follow 模式：每个新事件输出一行。
// JSON 模式输出 NDJSON（一行一对象），更适合 tail -f 配合 jq 流式处理。
func writeOneLine(w io.Writer, ev adapter.Event, asJSON bool) error {
	if asJSON {
		return output.Render(w, ev, "json", nil)
	}
	_, err := fmt.Fprintf(w, "%s  %-12s %-22s %-10s %-12s %s\n",
		shortTime(ev.CreatedAt),
		truncate(ev.Username, 12),
		truncate(ev.Action, 22),
		ev.Status,
		truncate(ev.ResourceType, 12),
		truncate(ev.Message, 80),
	)
	return err
}

// shortTime 把 ISO 时间戳缩短到 "MM-DD HH:MM:SS"，便于表格显示。
// 解析失败时原样返回。
func shortTime(s string) string {
	if s == "" {
		return ""
	}
	// CloudTower 时间戳通常是 RFC3339 / ISO8601。
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano, "2006-01-02T15:04:05Z"} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.Local().Format("01-02 15:04:05")
		}
	}
	return s
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}
