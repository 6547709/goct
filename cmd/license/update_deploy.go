package license

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newUpdateDeploy() *cobra.Command {
	c := &cobra.Command{
		Use:   "license.update-deploy <license-key>",
		Short: "Update deploy license",
		GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			_, err := service.NewLicense(cli).UpdateDeploy(c.Context(), args[0])
			return err
		},
	}
	return c
}
