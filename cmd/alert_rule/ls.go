package alert_rule

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

var alertRuleListColumns = []output.Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.AlertRule).ID }},
	{Header: "Enabled", Get: func(v any) string {
		if v.(adapter.AlertRule).Enabled {
			return "true"
		}
		return "false"
	}},
	{Header: "Target Kind", Get: func(v any) string { return v.(adapter.AlertRule).TargetKind }},
	{Header: "Target ID", Get: func(v any) string { return v.(adapter.AlertRule).TargetID }},
}

func newLs() *cobra.Command {
	var clusterID string
	c := &cobra.Command{
		Use:   "alert-rule.ls",
		Short: "List alert rules",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			items, err := service.NewAlertRule(cli).List(c.Context(), adapter.ListOpts{ClusterID: clusterID})
			if err != nil {
				return err
			}
			out := make([]any, len(items))
			for i := range items {
				out[i] = items[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, alertRuleListColumns)
		},
	}
	c.Flags().StringVar(&clusterID, "cluster", "", "Filter by cluster ID")
	return c
}
