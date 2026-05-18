package system

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/output"
	"github.com/spf13/cobra"
)

func newAbout() *cobra.Command {
	return &cobra.Command{
		Use:     "about",
		Short:   "Show CloudTower server version and connection info",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			info, err := cli.About(c.Context())
			if err != nil {
				return err
			}
			format, _ := c.Flags().GetString("format")
			if format == "json" {
				return output.Render(c.OutOrStdout(), []any{info}, format, nil)
			}
			fmt.Fprintf(c.OutOrStdout(), "Version: %s\nBuild:   %s\n", info.Version, info.Build)
			return nil
		},
	}
}
