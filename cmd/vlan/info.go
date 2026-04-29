package vlan

import (
	"fmt"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newInfo() *cobra.Command {
	return &cobra.Command{
		Use: "vlan.info <name|id>", Short: "Show VLAN details", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			v, err := service.NewVLAN(cli).Resolve(c.Context(), args[0])
			if err != nil { return err }
			format, _ := c.Flags().GetString("format")
			if format == "json" { return output.Render(c.OutOrStdout(), []any{*v}, format, nil) }
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "ID:", v.ID)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Name:", v.Name)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %d\n", "VLAN Tag:", v.VlanTag)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Type:", v.Type)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "VDS:", v.VdsID)
			return nil
		},
	}
}
