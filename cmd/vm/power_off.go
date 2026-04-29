package vm

import (
	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newPowerOff() *cobra.Command {
	var force bool
	c := &cobra.Command{
		Use: "vm.power.off <name|id>", Short: "Shut down / power off a VM", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewVM(cli).Power(c.Context(), args[0], adapter.PowerOff, force)
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
	c.Flags().BoolVar(&force, "force", false, "Force power off (skip graceful shutdown)")
	return c
}
