package host

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newMaintenanceEnter() *cobra.Command {
	return &cobra.Command{
		Use: "host.maintenance.enter <name|id>", Short: "Enter maintenance mode", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			ref, err := service.NewHost(cli).EnterMaintenance(c.Context(), args[0])
			if err != nil {
				return err
			}
			if ref.IsSync() {
				return nil
			}
			return task.New(cli, task.Options{Out: c.OutOrStderr()}).Watch(c.Context(), ref.ID)
		},
	}
}
