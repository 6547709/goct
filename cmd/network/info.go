package network

import (
	"fmt"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newInfo() *cobra.Command {
	return &cobra.Command{
		Use: "network.info <name|id>", Short: "Show virtual switch details", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			n, err := service.NewNetwork(cli).Resolve(c.Context(), args[0])
			if err != nil { return err }
			format, _ := c.Flags().GetString("format")
			if format == "json" { return output.Render(c.OutOrStdout(), []any{*n}, format, nil) }
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "ID:", n.ID)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Name:", n.Name)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Type:", n.Type)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Cluster:", n.ClusterID)
			return nil
		},
	}
}
