package content_library_image

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newImport() *cobra.Command {
	var clusterID string
	c := &cobra.Command{
		Use:   "content-library-image.import <path> <name>",
		Short: "Import content library image",
		GroupID: groupID,
		Args: cobra.ExactArgs(2),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			task, err := service.NewContentLibraryImage(cli).Import(c.Context(), args[0], args[1], clusterID)
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
	c.Flags().StringVar(&clusterID, "cluster", "", "Target cluster ID")
	c.MarkFlagRequired("cluster")
	return c
}
