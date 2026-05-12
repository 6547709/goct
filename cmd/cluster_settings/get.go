package cluster_settings

import (
	"fmt"

	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/service"
	"github.com/spf13/cobra"
)

func newGet() *cobra.Command {
	var clusterID string
	c := &cobra.Command{
		Use:   "cluster-settings.get",
		Short: "Get cluster settings",
		GroupID: groupID,
		RunE: func(c *cobra.Command, _ []string) error {
			cli := client.From(c.Context())
			settings, err := service.NewClusterSettings(cli).GetSettings(c.Context(), clusterID)
			if err != nil {
				return err
			}
			_, _ = fmt.Fprintf(c.OutOrStdout(), "Cluster ID: %s\n", settings.ClusterID)
			return nil
		},
	}
	c.Flags().StringVar(&clusterID, "cluster", "", "Cluster ID")
	c.MarkFlagRequired("cluster")
	return c
}
