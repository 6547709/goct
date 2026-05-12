package content_library_image

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newDistribute() *cobra.Command {
	var clusterIDs []string
	c := &cobra.Command{
		Use:   "content-library-image.distribute <id>",
		Short: "Distribute content library image to clusters",
		GroupID: groupID,
		Args: cobra.ExactArgs(1),
		RunE: func(c *cobra.Command, args []string) error {
			cli := client.From(c.Context())
			task, err := service.NewContentLibraryImage(cli).Distribute(c.Context(), args[0], clusterIDs)
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
	c.Flags().StringSliceVar(&clusterIDs, "cluster", nil, "Target cluster IDs")
	c.MarkFlagRequired("cluster")
	return c
}
