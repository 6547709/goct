package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newMigrate() *cobra.Command {
	var host string
	c := &cobra.Command{
		Use: "vm.migrate [name|id]", Short: "Migrate a VM to another host (omit --host to let CloudTower choose)", GroupID: groupID,
		Long: `Migrate a VM to another host within the same cluster.

If --host is omitted, CloudTower will pick a target host (DRS-like behavior).
If --host is given (name or ID), goct validates it is not the current host
and rejects the migration in that case.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			id, err := resolveVMArg(args)
			if err != nil {
				return err
			}
			ref, err := service.NewVM(cli).Migrate(c.Context(), id, host)
			if err != nil {
				return err
			}
			if ref.IsSync() {
				fmt.Fprintln(c.OutOrStdout(), "VM migrated (sync)")
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().StringVar(&host, "host", "", "Target host name or ID (omit = let CloudTower choose)")
	return c
}
