package application

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newDeploy() *cobra.Command {
	var targetPackage string
	c := &cobra.Command{
		Use:   "deploy-cloudtower-application <name>",
		Short: "Deploy CloudTower application",
		GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			task, err := service.NewCloudTowerApplication(cli).Deploy(c.Context(), args[0], targetPackage, nil)
			if err != nil {
				return err
			}
			if task.IsSync() {
				_, _ = fmt.Fprintln(c.OutOrStdout(), "Done")
			} else {
				_, _ = fmt.Fprintf(c.OutOrStdout(), "Task ID: %s\n", task.ID)
			}
			return nil
		},
	}
	c.Flags().StringVar(&targetPackage, "package", "", "Target package ID")
	c.MarkFlagRequired("package")
	return c
}
