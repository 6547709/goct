package cluster

import (
	"fmt"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newInfo() *cobra.Command {
	return &cobra.Command{
		Use: "cluster.info <name|id>", Short: "Show cluster details", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			cl, err := service.NewCluster(cli).Resolve(c.Context(), args[0])
			if err != nil { return err }
			format, _ := c.Flags().GetString("format")
			if format == "json" {
				return output.Render(c.OutOrStdout(), []any{*cl}, format, nil)
			}
			fmt.Fprintf(c.OutOrStdout(), "%-16s %s\n", "ID:", cl.ID)
			fmt.Fprintf(c.OutOrStdout(), "%-16s %s\n", "Name:", cl.Name)
			fmt.Fprintf(c.OutOrStdout(), "%-16s %d\n", "CPU Cores:", cl.TotalCPUCores)
			fmt.Fprintf(c.OutOrStdout(), "%-16s %s\n", "Memory:", output.HumanBytes(cl.TotalMemoryBytes))
			fmt.Fprintf(c.OutOrStdout(), "%-16s %s\n", "Storage:", output.HumanBytes(cl.TotalDataCapacity))
			fmt.Fprintf(c.OutOrStdout(), "%-16s %d\n", "Running VMs:", cl.RunningVMs)
			return nil
		},
	}
}
