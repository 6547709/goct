package vm

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/6547709/goct/pkg/task"
	"github.com/spf13/cobra"
)

func newMigrateAcrossCluster() *cobra.Command {
	var host string
	c := &cobra.Command{
		Use:   "vm.migrate.across [name|id] <cluster-name|id>",
		Short: "Migrate a VM to another cluster",
		Long: `Migrate a VM to another cluster, optionally specifying a target host.
Cluster and host can be names or IDs. If host is not specified,
CloudTower will automatically select an available host.

Examples:
  goct vm.migrate.across myvm Cluster02
  goct vm.migrate.across myvm Cluster02 --host SMTXOS03`,
		GroupID: groupID,
		Args:    cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			vmID := args[0]
			clusterIDOrName := args[1]
			ref, err := service.NewVM(cli).MigrateAcrossCluster(c.Context(), vmID, clusterIDOrName, host)
			if err != nil {
				return err
			}
			if ref.IsSync() {
				fmt.Fprintln(c.OutOrStdout(), "VM migrated across cluster (sync)")
				return nil
			}
			w := task.New(cli, task.Options{Out: c.OutOrStderr()})
			return w.Watch(c.Context(), ref.ID)
		},
	}
	c.Flags().StringVar(&host, "host", "", "Target host name or ID (optional, auto-select if not specified)")
	return c
}
