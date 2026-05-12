package application

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newDeletePackage() *cobra.Command {
	c := &cobra.Command{
		Use:   "delete-cloudtower-application-package <id>",
		Short: "Delete CloudTower application package",
		GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			return service.NewCloudTowerApplication(cli).DeletePackage(c.Context(), args[0])
		},
	}
	return c
}
