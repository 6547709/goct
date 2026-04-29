package alert

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	gflags "github.com/6547709/goct/pkg/flags"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newLs() *cobra.Command {
	var sf gflags.SearchFlags
	c := &cobra.Command{
		Use: "alert.ls", Short: "List alerts", GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			items, err := service.NewAlert(cli).List(c.Context(), adapter.ListOpts{Limit: sf.Limit})
			if err != nil { return err }
			out := make([]any, len(items)); for i := range items { out[i] = items[i] }
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, output.AlertListColumns)
		},
	}
	sf.Register(c); return c
}
