package datastore

import (
	"fmt"

	"github.com/6547709/goct/pkg/adapter"
	"github.com/6547709/goct/pkg/client"
	gflags "github.com/6547709/goct/pkg/flags"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newPoolLs() *cobra.Command {
	var sf gflags.SearchFlags
	var cluster string
	var idOnly bool
	c := &cobra.Command{
		Use:     "storage.pool.ls",
		Short:   "List hyperconverged storage pools (DiskPool)",
		GroupID: "datastore",
		Long: `List DiskPools - the hyperconverged storage pools on each host.
Each host has one DiskPool that aggregates all local disks into a distributed storage pool.

Examples:
  goct storage.pool.ls
  goct storage.pool.ls --cluster Cluster01
  goct storage.pool.ls --format json`,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			items, err := service.NewDatastore(cli).ListDiskPools(c.Context(),
				adapter.ListOpts{NameContains: sf.Name, ClusterID: cluster, Limit: sf.Limit})
			if err != nil {
				return err
			}
			if idOnly {
				for _, it := range items {
					_, _ = fmt.Fprintln(c.OutOrStdout(), it.ID)
				}
				return nil
			}
			out := make([]any, len(items))
			for i := range items {
				out[i] = items[i]
			}
			format, _ := c.Flags().GetString("format")
			return output.Render(c.OutOrStdout(), out, format, output.DiskPoolListColumns)
		},
	}
	sf.Register(c)
	c.Flags().StringVar(&cluster, "cluster", "", "Filter by cluster name or ID")
	c.Flags().BoolVar(&idOnly, "id-only", false, "Output only IDs, one per line (for scripting)")
	return c
}
