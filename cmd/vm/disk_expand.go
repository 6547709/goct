package vm

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newDiskExpand() *cobra.Command {
	var size string
	c := &cobra.Command{
		Use:   "vm.disk.expand [vm-name|vm-id] [disk-id]",
		Short: "Expand a VM disk",
		Long: `Expand (resize) a disk attached to a VM.

Examples:
  goct vm.disk.expand myvm disk-id --size 200G`,
		GroupID: groupID,
		Args:    cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			vmID, diskID := args[0], args[1]
			sizeBytes, err := parseSize(size)
			if err != nil {
				return err
			}
			ref, err := service.NewVM(cli).ExpandDisk(c.Context(), vmID, diskID, sizeBytes)
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
	c.Flags().StringVar(&size, "size", "", "New disk size (e.g. 200G)")
	return c
}
