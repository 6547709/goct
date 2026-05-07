package vm

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newGpuRm() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.gpu.rm [vm-name|vm-id] <gpu-device-id>",
		Short: "Remove a GPU device from VM",
		Long: `Remove a GPU device from a virtual machine.

Examples:
  goct vm.gpu.rm myvm gpu-device-001`,
		GroupID: groupID,
		Args:    cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			vmID := args[0]
			gpuDeviceID := args[1]
			ref, err := service.NewVM(cli).RemoveGpuDevice(c.Context(), vmID, gpuDeviceID)
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