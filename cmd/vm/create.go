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
		name        string
		clusterID   string
		vcpu        int32
		memoryMiB   int64
		firmware    string
		description string
	)
	c := &cobra.Command{
		Use: "vm.create", Short: "Create a new VM", GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
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
	c.Flags().StringVar(&name, "name", "", "VM name (required)")
	c.Flags().StringVar(&clusterID, "cluster", "", "Target cluster ID (required)")
	c.Flags().Int32Var(&vcpu, "vcpu", 1, "Number of vCPUs")
	c.Flags().Int64Var(&memoryMiB, "memory", 1024, "Memory in MiB")
	c.Flags().StringVar(&firmware, "firmware", "BIOS", "Firmware: BIOS or UEFI")
	c.Flags().StringVar(&description, "description", "", "Description")
	_ = c.MarkFlagRequired("name")
	_ = c.MarkFlagRequired("cluster")
	return c
}
