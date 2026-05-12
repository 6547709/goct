package deploy

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	gflags "github.com/6547709/goct/pkg/flags"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

var deployListColumns = []output.Column{
	{Header: "ID", Get: func(v any) string { return v.(adapter.Deploy).ID }},
	{Header: "Version", Get: func(v any) string { return v.(adapter.Deploy).VMID }},
}

func newLs() *cobra.Command {
	var sf gflags.SearchFlags
	c := &cobra.Command{
		Use:   "deploy.ls",
		Short: "List deploys",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			items, err := service.NewDeploy(cli).List(c.Context(), adapter.ListOpts{NameContains: sf.Name, Limit: sf.Limit})
			if err != nil {
				return err
			}
			out := make([]any, len(items))
			for i := range items {
				out[i] = items[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, deployListColumns)
		},
	}
	sf.Register(c)
	return c
}
