package vm

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newDiskRm() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.disk.rm [vm-name|vm-id] [disk-id]",
		Short: "Remove a disk from VM",
		Long: `Remove a disk from a virtual machine.

Examples:
  goct vm.disk.rm myvm disk-id`,
		GroupID: groupID,
		Args:    cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			vmID, diskID := args[0], args[1]
			ref, err := service.NewVM(cli).RemoveDisk(c.Context(), vmID, diskID)
			if err != nil {
				return err
			}
			if ref.IsSync() {
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	return c
}
