package host

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
		Use:   "host.ls",
		Short: "List hosts",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			hosts, err := service.NewHost(cli).List(c.Context(),
				adapter.ListOpts{NameContains: sf.Name, Limit: sf.Limit, Skip: sf.Skip})
			if err != nil {
				return err
			}
			if idOnly {
				for _, h := range hosts {
					_, _ = fmt.Fprintln(c.OutOrStdout(), h.ID)
				}
				return nil
			}
			items := make([]any, len(hosts))
			for i := range hosts {
				items[i] = hosts[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), items, format, output.HostListColumns)
		},
	}
	sf.Register(c)
	c.Flags().BoolVar(&idOnly, "id-only", false, "Output only IDs, one per line (for scripting)")
	return c
}
