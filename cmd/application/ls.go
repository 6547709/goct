package application

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	gflags "github.com/6547709/goct/pkg/flags"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

var cloudTowerApplicationListColumns = []output.Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.CloudTowerApplication).ID }},
	{Header: "Name", Get: func(v any) string { return v.(adapter.CloudTowerApplication).Name }},
	{Header: "State", Get: func(v any) string { return v.(adapter.CloudTowerApplication).State }},
	{Header: "Target Package", Get: func(v any) string { return v.(adapter.CloudTowerApplication).TargetPackage }},
}

func newGetApplications() *cobra.Command {
	var sf gflags.SearchFlags
	c := &cobra.Command{
		Use:   "get-cloudtower-applications",
		Short: "List CloudTower applications",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			items, err := service.NewCloudTowerApplication(cli).List(c.Context(), adapter.ListOpts{NameContains: sf.Name, Limit: sf.Limit})
			if err != nil {
				return err
			}
			out := make([]any, len(items))
			for i := range items {
				out[i] = items[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, cloudTowerApplicationListColumns)
		},
	}
	sf.Register(c)
	return c
}
