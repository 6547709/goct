package metrics

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func newClusterMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:   "cluster.metrics <metric> [cluster-name]",
		Short: "Query Cluster metrics (zbs_cluster_*)",
		Long:  "Query Cluster metrics with optional cluster name filter. Example: cluster.metrics zbs_cluster_usage cluster001",
		Args:  func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return ListMetrics(cmd.OutOrStdout(), "cluster")
			}
			cli := client.From(cmd.Context())
			mc := metrics.NewMetricsClient(cli)

			clusterName := ""
			if len(args) > 1 {
				clusterName = args[1]
			}

			input := &models.GetClusterMetricInput{
				Metrics:  []string{args[0]},
				Range:    &rangeFlag,
				Clusters: &models.ClusterWhereInput{Name: &clusterName},
			}

			results, err := mc.GetClusterMetrics(cmd.Context(), input)
			if err != nil {
				return err
			}

			return renderMetricsResults(cmd.OutOrStdout(), results, clusterName, args[0], "cluster")
		},
	}
	return c
}
