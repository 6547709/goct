package metrics

import (
	"github.com/6547709/goct/pkg/client"
	"github.com/6547709/goct/pkg/metrics"
	"github.com/spf13/cobra"
	"github.com/smartxworks/cloudtower-go-sdk/v2/models"
)

func newVolumeMetrics() *cobra.Command {
	c := &cobra.Command{
		Use:   "volume.metrics <metric> [volume-name]",
		Short: "Query Volume metrics",
		Long:  "Query independent volume metrics with optional volume name filter. Example: volume.metrics zbs_volume_read_iops volume001",
		Args: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return nil
			}
			return cobra.RangeArgs(1, 2)(cmd, args)
		},
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			if len(args) == 0 {
				names, err := GetMetricNames("volume")
				if err != nil {
					return nil, cobra.ShellCompDirectiveNoFileComp
				}
				return names, cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if listFlag {
				return ListMetrics(cmd.OutOrStdout(), "volume")
			}
			cli := client.From(cmd.Context())
			mc := metrics.NewMetricsClient(cli)

			volumeName := ""
			if len(args) > 1 {
				volumeName = args[1]
			}

			input := &models.GetVMVolumeMetricInput{
				Metrics:   []string{args[0]},
				Range:     &rangeFlag,
				VMVolumes: &models.VMVolumeWhereInput{Name: &volumeName},
			}

			results, err := mc.GetVmVolumeMetrics(cmd.Context(), input)
			if err != nil {
				return err
			}

			return renderMetricsResults(cmd.OutOrStdout(), results, volumeName, args[0], "volume")
		},
	}
	return c
}

func RegisterVolumeMetrics(root *cobra.Command) {
	c := newVolumeMetrics()
	c.Flags().BoolVar(&listFlag, "list", false, "List available metrics")
	c.Flags().StringVar(&rangeFlag, "range", "5m", "Time range: 5m, 1h, 1d, 7d")
	c.Flags().BoolVar(&latestFlag, "latest", false, "Show only latest value")
	c.Flags().StringVar(&formatFlag, "format", "table", "Output format: table, json, chart")
	root.AddCommand(c)
}