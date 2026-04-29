package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newSnapshotRevert() *cobra.Command {
	var vmFlag string
	c := &cobra.Command{
		Use: "vm.snapshot.revert <snapshot-id>", Short: "Revert VM to a snapshot", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewSnapshot(cli).Revert(c.Context(), vmFlag, args[0])
			if err != nil {
				return err
			}
			if ref.IsSync() {
				fmt.Fprintln(c.OutOrStdout(), "Snapshot reverted (sync)")
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().StringVar(&vmFlag, "vm", "", "VM name or ID (required)")
	_ = c.MarkFlagRequired("vm")
	return c
}
