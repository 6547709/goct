package vm

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newDiskAdd() *cobra.Command {
	var size, bus, name string
	c := &cobra.Command{
		Use:   "vm.disk.add [vm-name|vm-id]",
		Short: "Add a disk to VM",
		Long: `Add a new disk to a virtual machine.

Examples:
  goct vm.disk.add myvm --size 100G --name disk0
  goct vm.disk.add myvm --size 50G --bus SCSI`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			sizeBytes, err := parseSize(size)
			if err != nil {
				return err
			}
			spec := adapter.DiskAddSpec{
				Name:      name,
				SizeBytes: sizeBytes,
				Bus:       bus,
			}
			ref, err := service.NewVM(cli).AddDisk(c.Context(), id, spec)
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
	c.Flags().StringVar(&size, "size", "", "Disk size (e.g. 100G, 50G)")
	c.Flags().StringVar(&bus, "bus", "SCSI", "Disk bus type (SCSI, IDE, VIRTIO)")
	c.Flags().StringVar(&name, "name", "", "Disk name")
	return c
}
