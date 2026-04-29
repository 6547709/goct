package host

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newShutdown() *cobra.Command {
	var force bool
	c := &cobra.Command{
		Use: "host.shutdown [name|id]", Short: "Shut down a host", GroupID: groupID,
		Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveHostArg(args)
			if err != nil {
				return err
			}
			ref, err := service.NewHost(cli).Shutdown(c.Context(), id, force)
			if err != nil {
				return err
			}
			if ref.IsSync() {
				return nil
			}
			return task.New(cli, task.Options{Out: c.OutOrStderr()}).Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().BoolVar(&force, "force", false, "Force shutdown")
	return c
}
