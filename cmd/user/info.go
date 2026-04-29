package user

import (
	"fmt"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newInfo() *cobra.Command {
	return &cobra.Command{
		Use: "user.info <name|id>", Short: "Show user details", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			u, err := service.NewUser(cli).Resolve(c.Context(), args[0])
			if err != nil { return err }
			format, _ := c.Flags().GetString("format")
			if format == "json" { return output.Render(c.OutOrStdout(), []any{*u}, format, nil) }
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "ID:", u.ID)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Name:", u.Name)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Username:", u.Username)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Role:", u.Role)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Source:", u.Source)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Email:", u.Email)
			return nil
		},
	}
}
