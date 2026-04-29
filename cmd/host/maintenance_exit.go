package host

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newMaintenanceExit() *cobra.Command {
	return &cobra.Command{
		Use: "host.maintenance.exit [name|id]", Short: "Exit maintenance mode", GroupID: groupID,
		Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveHostArg(args)
			if err != nil {
				return err
			}
			ref, err := service.NewHost(cli).ExitMaintenance(c.Context(), id)
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
