package vm

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newNicAdd() *cobra.Command {
	var nicType, model string
	c := &cobra.Command{
		Use:   "vm.nic.add [vm-name|vm-id]",
		Short: "Add a NIC to VM",
		Long: `Add a new NIC to a virtual machine.

Examples:
  goct vm.nic.add myvm --type VLAN --model VIRTIO
  goct vm.nic.add myvm --type VPC --model E1000`,
		GroupID: groupID,
		Args:    cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			spec := adapter.NicAddSpec{
				Type:  nicType,
				Model: model,
			}
			ref, err := service.NewVM(cli).AddNic(c.Context(), id, spec)
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
	c.Flags().StringVar(&nicType, "type", "VLAN", "NIC type (VLAN, VPC)")
	c.Flags().StringVar(&model, "model", "VIRTIO", "NIC model (VIRTIO, E1000, SRIOV)")
	return c
}