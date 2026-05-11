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
		Args: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				names, err := GetMetricNames("cluster")
				if err != nil {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return names, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
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

func RegisterClusterMetrics(root *cobra.Command) {
	c := newClusterMetrics()
	c.Flags().BoolVar(&listFlag, "list", false, "List available metrics")
	c.Flags().StringVar(&rangeFlag, "range", "5m", "Time range: 5m, 1h, 1d, 7d")
	c.Flags().BoolVar(&latestFlag, "latest", false, "Show only latest value")
	c.Flags().StringVar(&formatFlag, "format", "table", "Output format: table, json, chart")
	root.AddCommand(c)
}