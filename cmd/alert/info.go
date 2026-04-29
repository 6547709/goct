package alert

import (
	"fmt"
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newInfo() *cobra.Command {
	return &cobra.Command{
		Use: "alert.info <id>", Short: "Show alert details", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			a, err := service.NewAlert(cli).Get(c.Context(), args[0])
			if err != nil { return err }
			format, _ := c.Flags().GetString("format")
			if format == "json" { return output.Render(c.OutOrStdout(), []any{*a}, format, nil) }
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "ID:", a.ID)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Severity:", a.Severity)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Message:", a.Message)
			fmt.Fprintf(c.OutOrStdout(), "%-12s %s\n", "Cause:", a.Cause)
			return nil
		},
	}
}
