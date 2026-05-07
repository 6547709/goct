package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newInfo() *cobra.Command {
	var detail bool
	c := &cobra.Command{
		Use: "vm.info [name|id]", Short: "Show VM details", GroupID: groupID,
		Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			v, err := service.NewVM(cli).Resolve(c.Context(), id)
			if err != nil {
				return err
			}
			format, _ := c.Flags().GetString("format")
			if format == "json" {
				return output.Render(c.OutOrStdout(), []any{*v}, format, nil)
			}
			// key-value table
			var rows [][]string
			if detail {
				rows = output.VMDetailRows(*v)
			} else {
				rows = output.VMInfoRows(*v)
			}
			for _, row := range rows {
				fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", row[0]+":", row[1])
			}
			return nil
		},
	}
	c.Flags().BoolVar(&detail, "detail", false, "Show detailed VM information (includes BIOS UUID, GPU/USB devices, usage stats, etc.)")
	return c
}
