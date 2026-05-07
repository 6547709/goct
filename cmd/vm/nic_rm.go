package vm

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newNicRm() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.nic.rm [vm-name|vm-id] <nic-index>",
		Short: "Remove a NIC from VM by index",
		Long: `Remove a NIC from a virtual machine by its index.

Examples:
  goct vm.nic.rm myvm 0
  goct vm.nic.rm myvm 1`,
		GroupID: groupID,
		Args:    cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			vmID := args[0]
			var idx int32
			if _, err := parseIndex(args[1], &idx); err != nil {
				return err
			}
			ref, err := service.NewVM(cli).RemoveNic(c.Context(), vmID, idx)
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