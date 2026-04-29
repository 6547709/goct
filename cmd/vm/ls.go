package vm

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
		Use:   "vm.ls",
		Short: "List virtual machines",
		Long: `List virtual machines with optional filtering.

Script-friendly: use --id-only to output only IDs (one per line),
then pipe to other commands:

  goct vm.ls --name web --id-only | xargs -I{} goct vm.info {}
  goct vm.ls --id-only | head -1 | xargs goct vm.power.on`,
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			vms, err := service.NewVM(cli).List(c.Context(),
				adapter.ListOpts{NameContains: sf.Name, Limit: sf.Limit, Skip: sf.Skip})
			if err != nil {
				return err
			}
			if idOnly {
				for _, v := range vms {
					fmt.Fprintln(c.OutOrStdout(), v.ID)
				}
				return nil
			}
			items := make([]any, len(vms))
			for i := range vms {
				items[i] = vms[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), items, format, output.VMListColumns)
		},
	}
	sf.Register(c)
	c.Flags().BoolVar(&idOnly, "id-only", false, "Output only IDs, one per line (for scripting)")
	return c
}
