package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newCreate() *cobra.Command {
	var (
		name          string
		clusterID     string
		vcpu          int32
		memoryMiB     int64
		firmware      string
		description   string
		fromTemplate  string
		isFullCopy    bool
		nicType       string
		nicModel      string
		nicVlan       string
	)
	c := &cobra.Command{
		Use: "vm.create", Short: "Create a new VM", GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())

			if fromTemplate != "" {
				ref, err := service.NewVM(cli).CreateFromTemplate(c.Context(), adapter.VMCreateFromTemplateSpec{
					TemplateID:  fromTemplate,
					Name:        name,
					ClusterID:   clusterID,
					VCPU:        vcpu,
					MemoryBytes: memoryMiB * 1024 * 1024,
					Firmware:    firmware,
					Description: description,
					IsFullCopy:  isFullCopy,
					NIC: adapter.NicConfig{
						Type:   nicType,
						Model:  nicModel,
						VlanID: nicVlan,
					},
				})
				if err != nil {
					return err
				}
				if ref.IsSync() {
					fmt.Fprintln(c.OutOrStdout(), "VM created from template (sync)")
					return nil
				}
				w := task.New(cli, task.Options{Out: c.OutOrStderr()})
				return w.Watch(c.Context(), ref.ID)
			}

			ref, err := service.NewVM(cli).Create(c.Context(), adapter.VMCreateSpec{
				Name:        name,
				ClusterID:   clusterID,
				VCPU:        vcpu,
				MemoryBytes: memoryMiB * 1024 * 1024, // MiB → bytes
				Firmware:    firmware,
				Description: description,
			})
			if err != nil {
				return err
			}
			if ref.IsSync() {
				fmt.Fprintln(c.OutOrStdout(), "VM created (sync)")
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().StringVar(&name, "name", "", "VM name")
	c.Flags().StringVar(&clusterID, "cluster", "", "Target cluster ID")
	c.Flags().Int32Var(&vcpu, "vcpu", 0, "Number of vCPUs (0 = use template default)")
	c.Flags().Int64Var(&memoryMiB, "memory", 0, "Memory in MiB (0 = use template default)")
	c.Flags().StringVar(&firmware, "firmware", "", "Firmware: BIOS or UEFI (empty = use template default)")
	c.Flags().StringVar(&description, "description", "", "Description")
	c.Flags().StringVar(&fromTemplate, "from-template", "", "Template ID to create VM from")
	c.Flags().BoolVar(&isFullCopy, "full-copy", false, "Full copy when creating from template")
	c.Flags().StringVar(&nicType, "nic-type", "", "NIC type: VLAN or VPC (only for --from-template)")
	c.Flags().StringVar(&nicModel, "nic-model", "", "NIC model: E1000, SRIOV, VIRTIO (only for --from-template)")
	c.Flags().StringVar(&nicVlan, "nic-vlan", "", "VLAN ID (only for --from-template)")
	return c
}
