package application

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newUploadPackage() *cobra.Command {
	c := &cobra.Command{
		Use:   "upload-cloudtower-application-package <path> <name>",
		Short: "Upload CloudTower application package",
		GroupID: groupID,
		Args: cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			task, err := service.NewCloudTowerApplication(cli).UploadPackage(c.Context(), args[0], args[1])
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
	return c
}
