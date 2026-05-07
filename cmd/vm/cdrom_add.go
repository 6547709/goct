package vm

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newCdRomAdd() *cobra.Command {
	var isoPath string
	c := &cobra.Command{
		Use:   "vm.cdrom.add [vm-name|vm-id]",
		Short: "Add a CD-ROM drive to VM",
		Long: `Add a CD-ROM drive to a VM, optionally with an ISO mounted.

Examples:
  goct vm.cdrom.add myvm
  goct vm.cdrom.add myvm --iso /path/to/boot.iso`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			ref, err := service.NewVM(cli).AddCdRom(c.Context(), id, isoPath)
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
	c.Flags().StringVar(&isoPath, "iso", "", "ISO path or image ID to mount")
	return c
}

func newCdRomEject() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.cdrom.eject [cdrom-id]",
		Short: "Eject ISO from CD-ROM drive",
		Long: `Eject the currently mounted ISO from a CD-ROM drive.`,
		GroupID: groupID,
		Args:    cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewVM(cli).EjectCdRom(c.Context(), args[0])
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

func newCdRomRm() *cobra.Command {
	c := &cobra.Command{
		Use:   "vm.cdrom.rm [vm-name|vm-id] [cdrom-id]",
		Short: "Remove a CD-ROM drive from VM",
		Long: `Remove a CD-ROM drive from a virtual machine.`,
		GroupID: groupID,
		Args:    cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			vmID, cdromID := args[0], args[1]
			ref, err := service.NewVM(cli).RemoveCdRom(c.Context(), vmID, cdromID)
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
