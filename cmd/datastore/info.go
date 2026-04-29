package datastore

import (
	"fmt"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newInfo() *cobra.Command {
	return &cobra.Command{
		Use: "datastore.info <name|id>", Short: "Show datastore details", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			d, err := service.NewDatastore(cli).Resolve(c.Context(), args[0])
			if err != nil { return err }
			format, _ := c.Flags().GetString("format")
			if format == "json" {
				return output.Render(c.OutOrStdout(), []any{*d}, format, nil)
			}
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "ID:", d.ID)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Name:", d.Name)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Type:", d.Type)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Cluster:", d.ClusterID)
			return nil
		},
	}
}
