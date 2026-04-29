package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newClone() *cobra.Command {
	var name, clusterID string
	c := &cobra.Command{
		Use: "vm.clone <source-name|id>", Short: "Clone a VM", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewVM(cli).Clone(c.Context(), args[0], adapter.VMCloneSpec{
				Name:            name,
				TargetClusterID: clusterID,
			})
			if err != nil {
				return err
			}
			if ref.IsSync() {
				fmt.Fprintln(c.OutOrStdout(), "VM cloned (sync)")
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().StringVar(&name, "name", "", "Name for cloned VM (required)")
	c.Flags().StringVar(&clusterID, "cluster", "", "Target cluster ID (optional, default same)")
	_ = c.MarkFlagRequired("name")
	return c
}
