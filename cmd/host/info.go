package host

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newInfo() *cobra.Command {
	return &cobra.Command{
		Use: "host.info <name|id>", Short: "Show host details", GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			h, err := service.NewHost(cli).Resolve(c.Context(), args[0])
			if err != nil {
				return err
			}
			format, _ := c.Flags().GetString("format")
			if format == "json" {
				return output.Render(c.OutOrStdout(), []any{*h}, format, nil)
			}
			for _, row := range output.HostInfoRows(*h) {
				fmt.Fprintf(c.OutOrStdout(), "%-16s %s\n", row[0]+":", row[1])
			}
			return nil
		},
	}
}
