package network

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	gflags "github.com/6547709/goct/pkg/flags"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newLs() *cobra.Command {
	var sf gflags.SearchFlags
	var idOnly bool
	c := &cobra.Command{
		Use: "network.ls", Short: "List virtual switches (VDS)", GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			items, err := service.NewNetwork(cli).List(c.Context(), adapter.ListOpts{NameContains: sf.Name, Limit: sf.Limit})
			if err != nil { return err }
			if idOnly {
				for _, it := range items {
					_, _ = fmt.Fprintln(c.OutOrStdout(), it.ID)
				}
				return nil
			}
			out := make([]any, len(items)); for i := range items { out[i] = items[i] }
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, output.NetworkListColumns)
		},
	}
	sf.Register(c)
	c.Flags().BoolVar(&idOnly, "id-only", false, "Output only IDs, one per line (for scripting)")
	return c
}
