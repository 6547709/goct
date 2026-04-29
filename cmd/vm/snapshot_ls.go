package vm

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newSnapshotLs() *cobra.Command {
	return &cobra.Command{
		Use: "vm.snapshot.ls <vm-name|id>", Short: "List snapshots of a VM", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			snaps, err := service.NewSnapshot(cli).List(c.Context(), args[0])
			if err != nil {
				return err
			}
			items := make([]any, len(snaps))
			for i := range snaps {
				items[i] = snaps[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), items, format, output.SnapshotListColumns)
		},
	}
}
