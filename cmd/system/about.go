package system

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/spf13/cobra"
)

func newAbout() *cobra.Command {
	return &cobra.Command{
		Use:     "about",
		Short:   "Show CloudTower server version and connection info",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			info, err := client.From(c.Context()).About(c.Context())
			if err != nil {
				return err
			}
			fmt.Fprintf(c.OutOrStdout(), "Version: %s\nBuild:   %s\n", info.Version, info.Build)
			return nil
		},
	}
}
