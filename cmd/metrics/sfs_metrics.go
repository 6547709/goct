package metrics

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newSFSMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:   "sfs.metrics <metric> [sfs-name]",
		Short: "Query SFS metrics (TODO: not yet implemented)",
		Long:  "Query SFS metrics. This feature is not yet available in the CloudTower SDK.",
		Args: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				names, err := GetMetricNames("sfs")
				if err != nil {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return names, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return ListMetrics(cmd.OutOrStdout(), "sfs")
			}
			return fmt.Errorf("SFS metrics API is not yet implemented in the CloudTower SDK")
		},
	}
	return c
}

func RegisterSFSMetrics(root *cobra.Command) {
	c := newSFSMetrics()
	c.Flags().BoolVar(&listFlag, "list", false, "List available metrics")
	c.Flags().StringVar(&rangeFlag, "range", "5m", "Time range: 5m, 1h, 1d, 7d")
	c.Flags().BoolVar(&latestFlag, "latest", false, "Show only latest value")
	c.Flags().StringVar(&formatFlag, "format", "table", "Output format: table, json, chart")
	root.AddCommand(c)
}