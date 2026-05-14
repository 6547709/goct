package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newDiskAdd() *cobra.Command {
	var size, bus, name, storagePolicy string
	c := &cobra.Command{
		Use:   "vm.disk.add [vm-name|vm-id]",
		Short: "Add a disk to VM",
		Long: `Add a new disk to a virtual machine.

Size must be an even number (in GiB). Storage policy must support the requested size.

Examples:
  goct vm.disk.add myvm --size 100G --name disk0
  goct vm.disk.add myvm --size 50G --bus SCSI
  goct vm.disk.add myvm --size 20G --storage-policy ELF_CP_REPLICA_2_THICK_PROVISION`,
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
			// 验证磁盘大小必须是偶数（单位 GiB）
			sizeGiB := sizeBytes / (1024 * 1024 * 1024)
			if sizeGiB%2 != 0 {
				return fmt.Errorf("disk size must be an even number (in GiB), got %dGiB", sizeGiB)
			}
			spec := adapter.DiskAddSpec{
				Name:          name,
				SizeBytes:     sizeBytes,
				Bus:           bus,
				StoragePolicy: storagePolicy,
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
	c.Flags().StringVar(&storagePolicy, "storage-policy", "", "Storage policy (e.g. ELF_CP_REPLICA_2_THICK_PROVISION, ELF_CP_REPLICA_3_THICK_PROVISION)")
	return c
}
