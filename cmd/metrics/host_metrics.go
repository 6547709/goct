package metrics

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func newHostMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:     "host.metrics <metric> [host-name]",
		Short:   "Query Host metrics",
		GroupID: "metrics",
		Long:  "Query Host metrics with optional host name filter. Example: host.metrics elf_host_cpu_usage host001",
		Args: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				names, err := GetMetricNames("host")
				if err != nil {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return names, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return ListMetrics(cmd.OutOrStdout(), "host")
			}
			cli := client.From(cmd.Context())
			mc := metrics.NewMetricsClient(cli)

			hostName := ""
			if len(args) > 1 {
				hostName = args[1]
			}

			input := &models.GetHostMetricInput{
				Metrics: []string{args[0]},
				Range:   &rangeFlag,
				Hosts:   &models.HostWhereInput{Name: &hostName},
			}

			results, err := mc.GetHostMetrics(cmd.Context(), input)
			if err != nil {
				return err
			}

			return renderMetricsResults(cmd.OutOrStdout(), results, hostName, args[0], "host")
		},
	}
	return c
}

func RegisterHostMetrics(root *cobra.Command) {
	c := newHostMetrics()
	c.Flags().BoolVar(&listFlag, "list", false, "List available metrics")
	c.Flags().StringVar(&rangeFlag, "range", "5m", "Time range: 5m, 1h, 1d, 7d")
	c.Flags().BoolVar(&latestFlag, "latest", false, "Show only latest value")
	c.Flags().StringVar(&formatFlag, "format", "table", "Output format: table, json, chart")
	root.AddCommand(c)
}